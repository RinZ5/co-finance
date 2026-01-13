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

func TestGetBasicFinancials(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stock/metric" {
			t.Errorf("Expected path /stock/metric, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"series": {
				"annual": {
					"currentRatio": [
						{"period": "2019-09-28", "v": 1.5401},
						{"period": "2018-09-29", "v": 1.1329}
					]
				}
			},
			"metric": {
				"10DayAverageTradingVolume": 32.50,
				"52WeekHigh": 310.43,
				"peBasicExclExtraTTM": 25.5
			},
			"metricType": "all",
			"symbol": "AAPL"
		}`))
	}))
	defer mockServer.Close()

	client := NewClient("fake-key")
	client.BaseURL = mockServer.URL

	financials, err := client.GetBasicFinancials("AAPL")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if financials.Metric.High52Week != 310.43 {
		t.Errorf("Expected 52WeekHigh 310.43, got %f", financials.Metric.High52Week)
	}

	if len(financials.Series.Annual.CurrentRatio) > 0 {
		firstRatio := financials.Series.Annual.CurrentRatio[0].V
		if firstRatio != 1.5401 {
			t.Errorf("Expected first Current Ratio 1.5401, got %f", firstRatio)
		}
	}

	if financials.MetricType != "all" {
		t.Errorf("Expected metric type 'all', got %s", financials.MetricType)
	}
}
