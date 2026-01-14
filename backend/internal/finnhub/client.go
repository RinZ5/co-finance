package finnhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rinz5/co-finance/backend/internal/models"
)

const DefaultBaseURL = "https://finnhub.io/api/v1"

type Client struct {
	ApiKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey:  apiKey,
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetQuote(symbol string) (*models.StockQuote, error) {
	url := fmt.Sprintf("%s/quote?symbol=%s&token=%s", c.BaseURL, symbol, c.ApiKey)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var quote models.StockQuote
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return nil, err
	}

	return &quote, nil
}

func (c *Client) GetBasicFinancials(symbol string) (*models.BasicFinancials, error) {
	url := fmt.Sprintf("%s/stock/metric?symbol=%s&metric=all&token=%s", c.BaseURL, symbol, c.ApiKey)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var financials models.BasicFinancials
	if err := json.NewDecoder(resp.Body).Decode(&financials); err != nil {
		return nil, err
	}

	return &financials, nil
}

func (c *Client) GetEarnings(symbol string) ([]models.EarningsSurprise, error) {
	url := fmt.Sprintf("%s/stock/earnings?symbol=%s&token=%s", c.BaseURL, symbol, c.ApiKey)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var earnings []models.EarningsSurprise
	if err := json.NewDecoder(resp.Body).Decode(&earnings); err != nil {
		return nil, err
	}

	return earnings, nil
}

func (c *Client) GetRecommendations(symbol string) ([]models.RecommendationTrend, error) {
	url := fmt.Sprintf("%s/stock/recommendation?symbol=%s&token=%s", c.BaseURL, symbol, c.ApiKey)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var recommendations []models.RecommendationTrend
	if err := json.NewDecoder(resp.Body).Decode(&recommendations); err != nil {
		return nil, err
	}

	return recommendations, nil
}

func (c *Client) GetInsiderTransactions(symbol string) ([]models.InsiderTransaction, error) {
	url := fmt.Sprintf("%s/stock/insider-transactions?symbol=%s&token=%s", c.BaseURL, symbol, c.ApiKey)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var wrapper struct {
		Data []models.InsiderTransaction `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	return wrapper.Data, nil
}

func (c *Client) GetCompanyNews(symbol, from, to string) ([]models.CompanyNews, error) {
	url := fmt.Sprintf("%s/company-news?symbol=%s&from=%s&to=%s&token=%s", c.BaseURL, symbol, from, to, c.ApiKey)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var news []models.CompanyNews
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return nil, err
	}

	return news, nil
}
