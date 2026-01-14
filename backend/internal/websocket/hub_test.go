package websocket

import (
	"testing"
)

type MockClient struct {
	send chan []byte
}

func TestHubBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	testMessage := []byte(`{"type":"trade", "data":[]}`)

	go func() {
		hub.Broadcast <- testMessage
	}()
}
