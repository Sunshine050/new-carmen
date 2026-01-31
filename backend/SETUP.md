# Setup Guide

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Ollama (for AI processing)
- ChromaDB (for vector database)

## Installation Steps

### 1. Navigate to Backend Directory

```bash
cd backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 2. Configure Environment

Copy `.env.example` to `.env` and fill in your configuration:

```bash
cp .env.example .env
```

Edit `.env` with your settings:
- Database credentials
- JWT secret
- Ollama URL and model
- ChromaDB URL
- GitHub token and repo info

### 3. Setup Database

```bash
# Create database
createdb carmen_db

# Run migrations (using GORM auto-migrate for now)
# The app will auto-migrate on startup
```

### 4. Start Services

#### Start PostgreSQL
```bash
# Windows
pg_ctl start

# Linux/Mac
sudo systemctl start postgresql
```

#### Start Ollama
```bash
ollama serve
```

#### Start ChromaDB
```bash
# Using Docker
docker run -d -p 8000:8000 chromadb/chroma

# Or using pip
pip install chromadb
chroma run --host localhost --port 8000
```

### 5. Run the Application

```bash
# Development
go run cmd/server/main.go

# Or using Make
make run

# Build
make build
./bin/server
```

## API Testing

### Register User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Search (with token)
```bash
curl -X GET "http://localhost:8080/api/search?q=test" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Project Structure

```
backend/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── database/               # DB connection
│   ├── models/                 # Data models
│   ├── repositories/           # Data access
│   ├── services/               # Business logic
│   ├── handlers/               # HTTP handlers
│   ├── middleware/             # HTTP middleware
│   ├── routes/                 # Route definitions
│   └── utils/                  # Utilities
├── pkg/
│   ├── ollama/                 # Ollama client
│   ├── chromadb/               # ChromaDB client
│   └── github/                 # GitHub client
└── migrations/                 # DB migrations
```

## Development Workflow

1. **Make changes** to code
2. **Test locally** with `go run cmd/server/main.go`
3. **Run tests** with `make test`
4. **Check linter** errors

## Troubleshooting

### Database Connection Error
- Check PostgreSQL is running
- Verify credentials in `.env`
- Ensure database exists

### Ollama Connection Error
- Verify Ollama is running on configured port
- Check model name is correct
- Test with: `curl http://localhost:11434/api/tags`

### ChromaDB Connection Error
- Verify ChromaDB is running
- Check collection name in config
- Create collection if needed via API

### GitHub API Error
- Verify GitHub token is valid
- Check repo name and owner
- Ensure token has repo access permissions
