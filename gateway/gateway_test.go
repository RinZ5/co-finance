package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func startMockFinnhub(t *testing.T, message string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			t.Logf("Mock Finnhub write error: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
	}))
}

func TestGatewayBroadcastWithAuth(t *testing.T) {
	expectedMsg := `{"data":"test_data"}`
	mockFinnhub := startMockFinnhub(t, expectedMsg)
	defer mockFinnhub.Close()

	finnhubWSUrl := "ws" + strings.TrimPrefix(mockFinnhub.URL, "http")
	testSecret := "very-secure-secret-token"

	gateway := NewGateway(finnhubWSUrl, testSecret)
	go gateway.Run()

	proxyServer := httptest.NewServer(http.HandlerFunc(gateway.ServeWS))
	defer proxyServer.Close()
	proxyUrl := "ws" + strings.TrimPrefix(proxyServer.URL, "http")

	headers := http.Header{}
	headers.Add("X-Gateway-Token", testSecret)

	client, _, err := websocket.DefaultDialer.Dial(proxyUrl, headers)
	if err != nil {
		t.Fatalf("Client failed to connect: %v", err)
	}
	defer client.Close()

	_, msg, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if string(msg) != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, string(msg))
	}
}

func TestGatewayRejectsUnauthorized(t *testing.T) {
	gateway := NewGateway("ws://dummy", "correct-secret")
	proxyServer := httptest.NewServer(http.HandlerFunc(gateway.ServeWS))
	defer proxyServer.Close()
	proxyUrl := "ws" + strings.TrimPrefix(proxyServer.URL, "http")

	_, resp, err := websocket.DefaultDialer.Dial(proxyUrl, nil)

	if err == nil {
		t.Fatal("Expected connection to fail, but it succeeded")
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}
