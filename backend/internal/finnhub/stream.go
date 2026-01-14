package finnhub

import (
	"log"

	"github.com/gorilla/websocket"
)

type StreamClient struct {
	Token   string
	Symbols []string
	BaseURL string
}

func NewStreamClient(token string, symbols []string) *StreamClient {
	return &StreamClient{
		Token:   token,
		Symbols: symbols,
		BaseURL: "wss://ws.finnhub.io?token=",
	}
}

func (s *StreamClient) Start(outputChan chan<- []byte) {
	url := s.BaseURL
	if s.BaseURL == "wss://ws.finnhub.io?token=" {
		url = s.BaseURL + s.Token
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("Finnhub connection error: %v", err)
		return
	}
	defer conn.Close()

	for _, sym := range s.Symbols {
		msg := map[string]interface{}{"type": "subscribe", "symbol": sym}
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Subscribe error: %v", err)
			return
		}
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}
		outputChan <- message
	}
}
