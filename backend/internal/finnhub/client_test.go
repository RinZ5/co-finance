package finnhub

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetQuote(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/quote" {
			t.Errorf("Expected path /quote, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
		  "c": 261.74,
		  "h": 263.31,
		  "l": 260.68,
		  "o": 261.07,
		  "pc": 259.45,
		  "t": 1582641000
		}`))
	}))
	defer mockServer.Close()

	client := NewClient("fake-key")
	client.BaseURL = mockServer.URL

	quote, err := client.GetQuote("AAPL")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if quote.CurrentPrice != 261.74 {
		t.Errorf("Expected price 261.74, got %f", quote.CurrentPrice)
	}
}
