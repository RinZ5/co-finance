package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"github.com/rinz5/stock-tracker/backend/internal/finnhub"
	"github.com/rinz5/stock-tracker/backend/internal/models"
)

type DashboardResponse struct {
	Quote           *models.StockQuote           `json:"quote"`
	Financials      *models.BasicFinancials      `json:"financials"`
	Earnings        []models.EarningsSurprise    `json:"earnings"`
	Recommendations []models.RecommendationTrend `json:"recommendations"`
	Insiders        []models.InsiderTransaction  `json:"insiders"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: .evn file is not found or could not be loaded.")
	}

	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: FINNHUB_API_KEY is not set.")
	}

	client := finnhub.NewClient(apiKey)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	r.Use(cors.New(config))

	r.GET("/api/quote", func(ctx *gin.Context) {
		symbol := ctx.Query("symbol")
		if symbol == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
			return
		}

		quote, err := client.GetQuote(symbol)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, quote)
	})

	r.GET("/api/financials", func(ctx *gin.Context) {
		symbol := ctx.Query("symbol")
		if symbol == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
			return
		}

		financials, err := client.GetBasicFinancials(symbol)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, financials)
	})

	r.GET("/api/earnings", func(ctx *gin.Context) {
		symbol := ctx.Query("symbol")
		if symbol == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
			return
		}

		earnings, err := client.GetEarnings(symbol)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, earnings)
	})

	r.GET("/api/recommendations", func(ctx *gin.Context) {
		symbol := ctx.Query("symbol")
		if symbol == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
			return
		}

		recommendations, err := client.GetRecommendations(symbol)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, recommendations)
	})

	r.GET("/api/insider", func(ctx *gin.Context) {
		symbol := ctx.Query("symbol")
		if symbol == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
			return
		}

		transactions, err := client.GetInsiderTransactions(symbol)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, transactions)
	})

	r.GET("/api/dashboard", func(ctx *gin.Context) {
		symbol := ctx.Query("symbol")
		if symbol == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
			return
		}

		var res DashboardResponse

		g, _ := errgroup.WithContext(ctx.Request.Context())

		g.Go(func() error {
			var err error
			res.Quote, err = client.GetQuote(symbol)
			return err
		})

		g.Go(func() error {
			var err error
			res.Financials, err = client.GetBasicFinancials(symbol)
			return err
		})

		g.Go(func() error {
			var err error
			res.Earnings, err = client.GetEarnings(symbol)
			return err
		})

		g.Go(func() error {
			var err error
			res.Recommendations, err = client.GetRecommendations(symbol)
			return err
		})

		g.Go(func() error {
			var err error
			res.Insiders, err = client.GetInsiderTransactions(symbol)
			return err
		})

		if err := g.Wait(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard data: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, res)
	})

	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
