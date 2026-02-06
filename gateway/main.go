package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Gateway struct {
	finnhubURL string
	secret     string
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	upstream   chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
}

func NewGateway(finnhubURL, secret string) *Gateway {
	return &Gateway{
		finnhubURL: finnhubURL,
		secret:     secret,
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		upstream:   make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (g *Gateway) Run() {
	go g.handleConnections()
	go g.connectToSource()
}

func (g *Gateway) readFromClient(conn *websocket.Conn) {
	defer func() {
		g.unregister <- conn
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		g.upstream <- message
	}
}

func (g *Gateway) ServeWS(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Gateway-Token")
	if token != g.secret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	g.register <- conn

	go g.readFromClient(conn)
}

func (g *Gateway) handleConnections() {
	for {
		select {
		case conn := <-g.register:
			g.mutex.Lock()
			g.clients[conn] = true
			g.mutex.Unlock()

		case conn := <-g.unregister:
			g.mutex.Lock()
			if _, ok := g.clients[conn]; ok {
				delete(g.clients, conn)
				conn.Close()
			}
			g.mutex.Unlock()

		case message := <-g.broadcast:
			g.mutex.Lock()
			for conn := range g.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					conn.Close()
					delete(g.clients, conn)
				}
			}
			g.mutex.Unlock()
		}
	}
}

func (g *Gateway) connectToSource() {
	backoff := 5 * time.Second

	for {
		log.Printf("Connecting to Finnhub...")

		conn, _, err := websocket.DefaultDialer.Dial(g.finnhubURL, nil)
		if err != nil {
			log.Printf("Dial error: %v. Retrying in %v...", err, backoff)
			time.Sleep(backoff)
			continue
		}

		log.Println("Connected to Finnhub successfully!")

		g.handleSession(conn)

		log.Println("Disconnected. Reconnecting in 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}

func (g *Gateway) handleSession(conn *websocket.Conn) {
	defer conn.Close()

	readErr := make(chan error, 1)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				readErr <- err
				return
			}
			g.broadcast <- message
		}
	}()

	for {
		select {
		case msg := <-g.upstream:
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("Finnhub Write Error: %v", err)
				return
			}

		case err := <-readErr:
			log.Printf("Finnhub Read Error: %v", err)
			return
		}
	}
}

func initializeEnvironment() (string, string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	apiKey := os.Getenv("FINNHUB_API_KEY")
	secret := os.Getenv("GATEWAY_SECRET")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if apiKey == "" || secret == "" {
		log.Fatal("Error: FINNHUB_API_KEY and GATEWAY_SECRET are required.")
	}

	return apiKey, secret, port
}

func main() {
	apiKey, secret, port := initializeEnvironment()

	finnhubURL := "wss://ws.finnhub.io?token=" + apiKey
	gateway := NewGateway(finnhubURL, secret)
	gateway.Run()

	http.HandleFunc("/ws", gateway.ServeWS)

	log.Printf("Gateway started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
