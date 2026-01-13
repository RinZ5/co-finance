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
			"d": 2.29,
			"dp": 0.88,
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

	earning := earnings[0]
	if earning.Actual != 1.88 {
		t.Errorf("Expected actual 1.88, got %f", earning.Actual)
	}

	if earning.Quarter != 1 {
		t.Errorf("Expected quarter 1, got %d", earning.Quarter)
	}

	if earning.Year != 2023 {
		t.Errorf("Expected year 2023, got %d", earning.Year)
	}
}

func TestGetRecommendations(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stock/recommendation" {
			t.Errorf("Expected path /stock/recommendation, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"buy": 24,
				"hold": 7,
				"period": "2025-03-01",
				"sell": 0,
				"strongBuy": 13,
				"strongSell": 0,
				"symbol": "AAPL"
			},
			{
				"buy": 17,
				"hold": 13,
				"period": "2025-02-01",
				"sell": 5,
				"strongBuy": 13,
				"strongSell": 0,
				"symbol": "AAPL"
			}
		]`))
	}))
	defer mockServer.Close()

	client := NewClient("fake-key")
	client.BaseURL = mockServer.URL

	recommendations, err := client.GetRecommendations("AAPL")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(recommendations) != 2 {
		t.Errorf("Expected 2 recommendation periods, got %d", len(recommendations))
	}

	rec := recommendations[0]
	if rec.Period != "2025-03-01" {
		t.Errorf("Expected period 2025-03-01, got %s", rec.Period)
	}
	if rec.Buy != 24 {
		t.Errorf("Expected 24 buys, got %d", rec.Buy)
	}
	if rec.Hold != 7 {
		t.Errorf("Expected 7 holds, got %d", rec.Hold)
	}
	if rec.StrongBuy != 13 {
		t.Errorf("Expected 13 strong buys, got %d", rec.StrongBuy)
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

func TestGetCompanyNews(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/company-news" {
			t.Errorf("Expected path /company-news, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("symbol") != "AAPL" {
			t.Errorf("Expected symbol AAPL, got %s", query.Get("symbol"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"category": "company news",
				"datetime": 1569550360,
				"headline": "More sops needed to boost electronic manufacturing",
				"id": 25286,
				"image": "https://img.etimg.com/photo.jpg",
				"related": "AAPL",
				"source": "The Economic Times India",
				"summary": "India may have to offer electronic manufacturers...",
				"url": "https://economictimes.indiatimes.com/articleshow/71321308.cms"
			},
			{
				"category": "company news",
				"datetime": 1569528720,
				"headline": "How to disable comments on your YouTube videos",
				"id": 25287,
				"image": "https://amp.businessinsider.com/images/5d8d1618.jpg",
				"related": "AAPL",
				"source": "Business Insider",
				"summary": "You can disable comments on your own YouTube video...",
				"url": "https://www.businessinsider.com/how-to-disable-comments-on-youtube"
			},
			{
				"category": "company news",
				"datetime": 1569526180,
				"headline": "Apple iPhone 11 Pro Teardowns Look Encouraging",
				"id": 25341,
				"image": "http://s.thestreet.com/files/tsc/v2008/photos.png",
				"related": "AAPL",
				"source": "TheStreet",
				"summary": "STMicroelectronics and Sony each appear to be supplying...",
				"url": "https://realmoney.thestreet.com/investing/technology/iphone-11-pro"
			}
		]`))
	}))
	defer mockServer.Close()

	client := NewClient("fake-key")
	client.BaseURL = mockServer.URL

	news, err := client.GetCompanyNews("AAPL", "2023-01-01", "2023-01-07")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(news) != 3 {
		t.Errorf("Expected 3 news items, got %d", len(news))
	}

	item := news[0]
	if item.Id != 25286 {
		t.Errorf("Expected ID 25286, got %d", item.Id)
	}
	if item.Source != "The Economic Times India" {
		t.Errorf("Expected source 'The Economic Times India', got %s", item.Source)
	}
	if item.Datetime != 1569550360 {
		t.Errorf("Expected datetime 1569550360, got %d", item.Datetime)
	}
}
