package finnhub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type StreamClient struct {
	Token         string
	Symbols       []string
	BaseURL       string
	subscribeChan chan string
	conn          *websocket.Conn
	mu            sync.Mutex
}

func NewStreamClient(token string, symbols []string) *StreamClient {
	return &StreamClient{
		Token:         token,
		Symbols:       symbols,
		BaseURL:       "wss://ws.finnhub.io?token=",
		subscribeChan: make(chan string),
	}
}

func (s *StreamClient) Subscribe(symbol string) {
	s.subscribeChan <- symbol
}

func (s *StreamClient) Start(outputChan chan<- []byte) {
	url := s.BaseURL
	if s.BaseURL == "wss://ws.finnhub.io?token=" {
		url = s.BaseURL + s.Token
	}

	var err error
	s.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("Finnhub connection error: %v", err)
		return
	}
	defer s.conn.Close()

	for _, sym := range s.Symbols {
		s.sendSubscribe(sym)
	}

	go s.readLoop(outputChan)

	for symbol := range s.subscribeChan {
		s.sendSubscribe(symbol)
	}
}

func (s *StreamClient) sendSubscribe(symbol string) {
	msg := map[string]any{"type": "subscribe", "symbol": symbol}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.conn.WriteJSON(msg); err != nil {
		log.Printf("Subscribe error: %v", err)
	} else {
		log.Printf("Subscribed to %s", symbol)
	}
}

func (s *StreamClient) readLoop(outputChan chan<- []byte) {
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}
		outputChan <- message
	}
}
