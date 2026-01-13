package finnhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rinz5/stock-tracker/backend/internal/models"
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
