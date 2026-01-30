# Co-Finance

A financial dashboard application with Vue.js frontend and Go backend, providing real-time stock market data through Finnhub API.

## Quick Start

### Option 1: Docker-compose
1. Copy environment file: `cp .env.example .env`
2. Set your Finnhub API key in `.env`
3. Run: `docker-compose up --build`
4. Access: http://localhost:3000

### Option 2: Separate Deployment
1. Backend: `cp backend/.env.example backend/.env`
2. Frontend: `cp frontend/.env.example frontend/.env`
3. Configure both files and deploy separately

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

### Docker-compose (Root .env)
- **FINNHUB_API_KEY**: Your Finnhub API key
- **ALLOWED_ORIGINS**: Comma-separated allowed domains

### Separate Deployment

**Backend (backend/.env)**:
- **FINNHUB_API_KEY**: Your Finnhub API key
- **ALLOWED_ORIGINS**: Comma-separated allowed domains

**Frontend (frontend/.env)**:
- **VITE_BACKEND_URL**: Backend API base URL