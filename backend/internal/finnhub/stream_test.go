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
		authToken := r.Header.Get("X-Gateway-Token")
		if authToken != "fake-token" {
			t.Errorf("Expected X-Gateway-Token 'fake-token', got '%s'", authToken)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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
	client := NewStreamClient(wsURL, "fake-token", []string{"AAPL"})

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
		authToken := r.Header.Get("X-Gateway-Token")
		if authToken != "fake-token" {
			t.Errorf("Expected X-Gateway-Token 'fake-token', got '%s'", authToken)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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
	client := NewStreamClient(wsURL, "fake-token", []string{})

	go client.Start(make(chan []byte))

	time.Sleep(50 * time.Millisecond)
	client.Subscribe("TSLA")

	time.Sleep(50 * time.Millisecond)
}

func TestStreamClientAuthenticationFailure(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("X-Gateway-Token")
		if authToken != "valid-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
	}))
	defer mockServer.Close()

	wsURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
	client := NewStreamClient(wsURL, "invalid-token", []string{})
	dataChan := make(chan []byte)

	go client.Start(dataChan)

	time.Sleep(100 * time.Millisecond)
}

func TestStreamClientConnectionRetry(t *testing.T) {
	connectionAttempts := 0
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connectionAttempts++

		authToken := r.Header.Get("X-Gateway-Token")
		if authToken != "test-token" {
			t.Errorf("Expected X-Gateway-Token 'test-token', got '%s'", authToken)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if connectionAttempts == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		_, msg, _ := c.ReadMessage()
		if !strings.Contains(string(msg), "subscribe") {
			t.Errorf("Expected subscribe message, got %s", msg)
		}

		fakeTrade := `{"type":"trade","data":[{"p":150.0,"s":"GOOGL"}]}`
		c.WriteMessage(websocket.TextMessage, []byte(fakeTrade))
	}))
	defer mockServer.Close()

	wsURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
	client := NewStreamClient(wsURL, "test-token", []string{"GOOGL"})
	dataChan := make(chan []byte)

	go client.Start(dataChan)

	select {
	case msg := <-dataChan:
		if !strings.Contains(string(msg), "150.0") {
			t.Errorf("Expected trade data 150.0, got %s", msg)
		}
		if connectionAttempts != 2 {
			t.Errorf("Expected 2 connection attempts, got %d", connectionAttempts)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Did not receive message after retry")
	}
}

func TestStreamClientEmptySecret(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("X-Gateway-Token")
		if authToken != "" {
			t.Errorf("Expected no X-Gateway-Token header, got '%s'", authToken)
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		_, msg, _ := c.ReadMessage()
		if !strings.Contains(string(msg), "subscribe") {
			t.Errorf("Expected subscribe message, got %s", msg)
		}

		fakeTrade := `{"type":"trade","data":[{"p":200.0,"s":"MSFT"}]}`
		c.WriteMessage(websocket.TextMessage, []byte(fakeTrade))
	}))
	defer mockServer.Close()

	wsURL := "ws" + strings.TrimPrefix(mockServer.URL, "http")
	client := NewStreamClient(wsURL, "", []string{"MSFT"})
	dataChan := make(chan []byte)

	go client.Start(dataChan)

	select {
	case msg := <-dataChan:
		if !strings.Contains(string(msg), "200.0") {
			t.Errorf("Expected trade data 200.0, got %s", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Did not receive message from mock server in time")
	}
}
