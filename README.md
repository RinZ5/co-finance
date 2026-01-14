# Co-Finance

A financial dashboard application with Vue.js frontend and Go backend, providing real-time stock market data through Finnhub API.

## Quick Start (Docker)

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Set your Finnhub API key in `.env`:
   ```
   FINNHUB_API_KEY=your_actual_api_key
   ```

3. Build and run:
   ```bash
   docker-compose up --build -d
   ```

4. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## Local Development

### Backend
```bash
cd backend
go mod download
go run ./cmd/server/main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## Environment Variables

- `FINNHUB_API_KEY`: Your Finnhub API key (required)
- `ALLOWED_ORIGINS`: CORS allowed origins (defaults to localhost ports)