package finnhub

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type StreamClient struct {
	URL           string
	Secret        string
	Symbols       []string
	subscribeChan chan string
	conn          *websocket.Conn
	mu            sync.Mutex
}

func NewStreamClient(url, secret string, symbols []string) *StreamClient {
	return &StreamClient{
		URL:           url,
		Secret:        secret,
		Symbols:       symbols,
		subscribeChan: make(chan string),
	}
}

func (s *StreamClient) Subscribe(symbol string) {
	go func() {
		s.subscribeChan <- symbol
	}()
}

func (s *StreamClient) Start(outputChan chan<- []byte) {
	for {
		log.Printf("Connecting to Gateway: %s", s.URL)

		headers := http.Header{}
		if s.Secret != "" {
			headers.Add("X-Gateway-Token", s.Secret)
		}

		var err error
		conn, _, err := websocket.DefaultDialer.Dial(s.URL, headers)
		if err != nil {
			log.Printf("Gateway connection failed: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		s.mu.Lock()
		s.conn = conn
		s.mu.Unlock()

		log.Println("Connected to Gateway!")
		log.Printf("Gateway connection established from local address: %s", conn.LocalAddr())

		for _, sym := range s.Symbols {
			s.sendSubscribe(sym)
		}

		done := make(chan struct{})
		go s.readLoop(outputChan, done)

	SessionLoop:
		for {
			select {
			case symbol := <-s.subscribeChan:
				s.addSymbol(symbol)
				s.sendSubscribe(symbol)

			case <-done:
				break SessionLoop
			}
		}

		s.mu.Lock()
		if s.conn != nil {
			s.conn.Close()
			s.conn = nil
		}
		s.mu.Unlock()

		log.Println("Disconnected. Reconnecting in 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}

func (s *StreamClient) addSymbol(symbol string) {
	for _, exist := range s.Symbols {
		if exist == symbol {
			return
		}
	}
	s.Symbols = append(s.Symbols, symbol)
}

func (s *StreamClient) sendSubscribe(symbol string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		return
	}

	msg := map[string]any{"type": "subscribe", "symbol": symbol}
	if err := s.conn.WriteJSON(msg); err != nil {
		log.Printf("Subscribe error: %v", err)
	} else {
		log.Printf("Subscribed to %s", symbol)
	}
}

func (s *StreamClient) readLoop(outputChan chan<- []byte, done chan struct{}) {
	defer close(done)

	for {
		if s.conn == nil {
			return
		}

		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}
		outputChan <- message
	}
}
