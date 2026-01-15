package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"github.com/rinz5/co-finance/backend/internal/finnhub"
	"github.com/rinz5/co-finance/backend/internal/models"
	"github.com/rinz5/co-finance/backend/internal/websocket"

	ws "github.com/gorilla/websocket"
)

const ErrSymbolRequired = "Symbol is required"

type DashboardResponse struct {
	Quote           *models.StockQuote           `json:"quote"`
	Financials      *models.BasicFinancials      `json:"financials"`
	Earnings        []models.EarningsSurprise    `json:"earnings"`
	Recommendations []models.RecommendationTrend `json:"recommendations"`
	Insiders        []models.InsiderTransaction  `json:"insiders"`
}

type Server struct {
	hub      *websocket.Hub
	streamer *finnhub.StreamClient
	client   *finnhub.Client
}

func initializeEnvironment() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: FINNHUB_API_KEY is not set.")
	}

	return apiKey
}

func setupServer(apiKey string) *Server {
	hub := websocket.NewHub()
	go hub.Run()

	streamer := finnhub.NewStreamClient(apiKey, []string{"AAPL"})
	go streamer.Start(hub.Broadcast)

	client := finnhub.NewClient(apiKey)

	return &Server{
		hub:      hub,
		streamer: streamer,
		client:   client,
	}
}

func setupOrigins() (map[string]bool, []string) {
	rawOrigins := os.Getenv("ALLOWED_ORIGINS")
	if rawOrigins == "" {
		rawOrigins = "http://localhost:5173,http://127.0.0.1:5173"
	}

	wsAllowedOrigins := make(map[string]bool)
	var corsAllowedOrigins []string

	for origin := range strings.SplitSeq(rawOrigins, ",") {
		cleanOrigin := strings.TrimSpace(origin)
		if cleanOrigin == "" {
			continue
		}

		wsAllowedOrigins[cleanOrigin] = true
		corsAllowedOrigins = append(corsAllowedOrigins, cleanOrigin)
	}

	return wsAllowedOrigins, corsAllowedOrigins
}

func (s *Server) setupWebSocketHandler(wsAllowedOrigins map[string]bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var upgrader = ws.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")

				if wsAllowedOrigins[origin] {
					return true
				}

				log.Printf("SECURITY: Blocked connection attempt from unauthorized Origin: %s", origin)
				return false
			},
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("Upgrade failed:", err)
			return
		}

		s.hub.Register <- conn

		defer func() {
			s.hub.Unregister <- conn
			conn.Close()
		}()

		for {
			var msg struct {
				Type   string `json:"type"`
				Symbol string `json:"symbol"`
			}

			if err := conn.ReadJSON(&msg); err != nil {
				break
			}

			if msg.Type == "subscribe" && msg.Symbol != "" {
				log.Printf("Frontend requested subscription to: %s", msg.Symbol)
				s.streamer.Subscribe(msg.Symbol)
			}
		}
	}
}

func (s *Server) validateSymbol(ctx *gin.Context) (string, bool) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrSymbolRequired})
		return "", false
	}
	return symbol, true
}

func (s *Server) handleQuote(ctx *gin.Context) {
	symbol, ok := s.validateSymbol(ctx)
	if !ok {
		return
	}

	quote, err := s.client.GetQuote(symbol)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, quote)
}

func (s *Server) handleFinancials(ctx *gin.Context) {
	symbol, ok := s.validateSymbol(ctx)
	if !ok {
		return
	}

	financials, err := s.client.GetBasicFinancials(symbol)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, financials)
}

func (s *Server) handleEarnings(ctx *gin.Context) {
	symbol, ok := s.validateSymbol(ctx)
	if !ok {
		return
	}

	earnings, err := s.client.GetEarnings(symbol)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, earnings)
}

func (s *Server) handleRecommendations(ctx *gin.Context) {
	symbol, ok := s.validateSymbol(ctx)
	if !ok {
		return
	}

	recommendations, err := s.client.GetRecommendations(symbol)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recommendations)
}

func (s *Server) handleInsider(ctx *gin.Context) {
	symbol, ok := s.validateSymbol(ctx)
	if !ok {
		return
	}

	transactions, err := s.client.GetInsiderTransactions(symbol)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (s *Server) handleDashboard(ctx *gin.Context) {
	symbol, ok := s.validateSymbol(ctx)
	if !ok {
		return
	}

	var res DashboardResponse

	g, _ := errgroup.WithContext(ctx.Request.Context())

	g.Go(func() error {
		var err error
		res.Quote, err = s.client.GetQuote(symbol)
		return err
	})

	g.Go(func() error {
		var err error
		res.Financials, err = s.client.GetBasicFinancials(symbol)
		return err
	})

	g.Go(func() error {
		var err error
		res.Earnings, err = s.client.GetEarnings(symbol)
		return err
	})

	g.Go(func() error {
		var err error
		res.Recommendations, err = s.client.GetRecommendations(symbol)
		return err
	})

	g.Go(func() error {
		var err error
		res.Insiders, err = s.client.GetInsiderTransactions(symbol)
		return err
	})

	if err := g.Wait(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard data: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *Server) handleCompanyNews(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	from := ctx.Query("from")
	to := ctx.Query("to")

	if symbol == "" || from == "" || to == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Symbol, from, and to parameters are required"})
		return
	}

	news, err := s.client.GetCompanyNews(symbol, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, news)
}

func (s *Server) handleMarketStatus(ctx *gin.Context) {
	exchange := ctx.Query("exchange")

	if exchange == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Exchange parametes is required"})
		return
	}

	market, err := s.client.GetMarketStatus(exchange)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, market)
}

func (s *Server) setupRoutes(r *gin.Engine) {
	wsAllowedOrigins, _ := setupOrigins()

	r.GET("/ws", s.setupWebSocketHandler(wsAllowedOrigins))
	r.GET("/api/quote", s.handleQuote)
	r.GET("/api/financials", s.handleFinancials)
	r.GET("/api/earnings", s.handleEarnings)
	r.GET("/api/recommendations", s.handleRecommendations)
	r.GET("/api/insider", s.handleInsider)
	r.GET("/api/dashboard", s.handleDashboard)
	r.GET("/api/company-news", s.handleCompanyNews)
	r.GET("/api/market-status", s.handleMarketStatus)
}

func main() {
	apiKey := initializeEnvironment()
	server := setupServer(apiKey)

	r := gin.Default()

	_, corsAllowedOrigins := setupOrigins()

	config := cors.DefaultConfig()
	config.AllowOrigins = corsAllowedOrigins
	r.Use(cors.New(config))

	server.setupRoutes(r)

	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
