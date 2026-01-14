package finnhub

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func TestStreamClient(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		_, msg, _ := c.ReadMessage()
		if !strings.Contains(string(msg), "subscribe") {
			t.Errorf("Expected subscribe message, got %s", msg)
		}

		fakeTrade := `{"type":"trade","data":[{"p":100.5,"s":"AAPL"}]}`
		c.WriteMessage(websocket.TextMessage, []byte(fakeTrade))
	}))
	defer mockServer.Close()

	wsURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")

	dataChan := make(chan []byte)
	client := NewStreamClient("fake-token", []string{"AAPL"})

	client.BaseURL = wsURL

	go client.Start(dataChan)

	select {
	case msg := <-dataChan:
		if !strings.Contains(string(msg), "100.5") {
			t.Errorf("Expected trade data 100.5, got %s", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Did not receive message from mock server in time")
	}
}

func TestStreamClientSubscription(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		_, msg, _ := c.ReadMessage()
		if !strings.Contains(string(msg), "subscribe") || !strings.Contains(string(msg), "TSLA") {
			t.Errorf("Expected subscription for TSLA, got %s", msg)
		}
	}))
	defer mockServer.Close()

	wsURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
	client := NewStreamClient("fake-token", []string{})
	client.BaseURL = wsURL

	go client.Start(make(chan []byte))

	time.Sleep(50 * time.Millisecond)
	client.Subscribe("TSLA")

	time.Sleep(50 * time.Millisecond)
}
