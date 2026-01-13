package models

// https://finnhub.io/docs/api/quote
type StockQuote struct {
	CurrentPrice  float64 `json:"c"`
	Change        float64 `json:"d"`
	PercentChange float64 `json:"dp"`
	HighPrice     float64 `json:"h"`
	LowPrice      float64 `json:"l"`
	OpenPrice     float64 `json:"o"`
	PrevClose     float64 `json:"pc"`
	Timestamp     float64 `json:"t"`
}

// https://finnhub.io/docs/api/company-basic-financials
type BasicFinancials struct {
	Series struct {
		Annual struct {
			CurrentRatio []struct {
				Period string  `json:"period"`
				V      float64 `json:"v"`
			} `json:"currentRatio"`
			SalesPerShare []struct {
				Period string  `json:"period"`
				V      float64 `json:"v"`
			} `json:"salesPerShare"`
			NetMargin []struct {
				Period string  `json:"period"`
				V      float64 `json:"v"`
			} `json:"netMargin"`
		}
	} `json:"Series"`
	Metric struct {
		PeBasicExclExtraTTM  float64 `json:"peBasicExclExtraTTM"`
		MarketCapitalization float64 `json:"marketCapitalization"`
		High52Week           float64 `json:"52WeekHigh"`
		Low52Week            float64 `json:"52WeekLow"`
		DividendYield        float64 `json:"dividendYieldIndicatedAnnual"`
		Beta                 float64 `json:"beta"`
	} `json:"metric"`
	MetricType string `json:"metricType"`
	Symbol     string `json:"symbol"`
}

// https://finnhub.io/docs/api/company-earnings
type EarningsSurprise struct {
	Actual          float64 `json:"actual"`
	Estimate        float64 `json:"estimate"`
	Period          string  `json:"period"`
	Quarter         int     `json:"quarter"`
	Year            int     `json:"year"`
	Surprise        float64 `json:"surprise"`
	SurprisePercent float64 `json:"surprisePercent"`
	Symbol          string  `json:"symbol"`
}

// https://finnhub.io/docs/api/recommendation-trends
type RecommendationTrend struct {
	Buy        int    `json:"buy"`
	Hold       int    `json:"hold"`
	Period     string `json:"period"`
	Sell       int    `json:"sell"`
	StrongBuy  int    `json:"strongBuy"`
	StrongSell int    `json:"strongSell"`
	Symbol     string `json:"symbol"`
}

// https://finnhub.io/docs/api/insider-transactions
type InsiderTransactionResponse struct {
	Data   []InsiderTransaction `json:"data"`
	Symbol string               `json:"symbol"`
}

type InsiderTransaction struct {
	Name             string  `json:"name"`
	Share            float64 `json:"share"`
	Change           float64 `json:"change"`
	FilingDate       string  `json:"filingDate"`
	TransactionDate  string  `json:"transactionDate"`
	TransactionPrice float64 `json:"transactionPrice"`
	Symbol           string  `json:"symbol"`
}
