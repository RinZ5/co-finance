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

func TestGetEarnings(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stock/earnings" {
			t.Errorf("Expected path /stock/earnings, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"actual": 1.88,
				"estimate": 1.9744,
				"period": "2023-03-31",
				"quarter": 1,
				"surprise": -0.0944,
				"surprisePercent": -4.7812,
				"symbol": "AAPL",
				"year": 2023
			},
			{
				"actual": 1.29,
				"estimate": 1.2957,
				"period": "2022-12-31",
				"quarter": 4,
				"surprise": -0.0057,
				"surprisePercent": -0.4399,
				"symbol": "AAPL",
				"year": 2022
			},
			{
				"actual": 1.2,
				"estimate": 1.1855,
				"period": "2022-09-30",
				"quarter": 3,
				"surprise": 0.0145,
				"surprisePercent": 1.2231,
				"symbol": "AAPL",
				"year": 2022
			}
		]`))
	}))
	defer mockServer.Close()

	client := NewClient("fake-key")
	client.BaseURL = mockServer.URL

	earnings, err := client.GetEarnings("AAPL")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(earnings) != 3 {
		t.Errorf("Expected 3 earnings records, got %d", len(earnings))
	}

	if earnings[0].Actual != 1.88 {
		t.Errorf("Expected actual 1.88, got %f", earnings[0].Actual)
	}

	if earnings[0].Quarter != 1 {
		t.Errorf("Expected quarter 1, got %d", earnings[0].Quarter)
	}

	if earnings[0].Year != 2023 {
		t.Errorf("Expected year 2023, got %d", earnings[0].Year)
	}
}

func TestGetInsiderTransactions(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stock/insider-transactions" {
			t.Errorf("Expected path /stock/insider-transactions, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": [
				{
					"name": "Kirkhorn Zachary",
					"share": 57234,
					"change": -1250,
					"filingDate": "2021-03-19",
					"transactionDate": "2021-03-17",
					"transactionCode": "S",
					"transactionPrice": 655.81
				},
				{
					"name": "Baglino Andrew D",
					"share": 20614,
					"change": 1000,
					"filingDate": "2021-03-31",
					"transactionDate": "2021-03-29",
					"transactionCode": "M",
					"transactionPrice": 41.57
				},
				{
					"name": "Baglino Andrew D",
					"share": 19114,
					"change": -1500,
					"filingDate": "2021-03-31",
					"transactionDate": "2021-03-29",
					"transactionCode": "S",
					"transactionPrice": 615.75
				}
			],
			"symbol": "TSLA"
		}`))
	}))
	defer mockServer.Close()

	client := NewClient("fake-key")
	client.BaseURL = mockServer.URL

	transactions, err := client.GetInsiderTransactions("AAPL")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(transactions) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(transactions))
	}

	if transactions[0].Name != "Kirkhorn Zachary" {
		t.Errorf("Expected name Kirkhorn Zachary, got %s", transactions[0].Name)
	}
	if transactions[0].Change != -1250 {
		t.Errorf("Expected change -1250, got %f", transactions[0].Change)
	}
}
