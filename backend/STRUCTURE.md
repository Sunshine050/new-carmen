# Project Structure

โครงสร้างโปรเจคที่ปรับให้ตรงกับ Go Standard Project Layout และ Fiber Framework best practices

## โครงสร้างโฟลเดอร์

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
│
├── internal/                    # Private application code
│   ├── api/                     # HTTP handlers (Fiber convention)
│   │   ├── auth_handler.go
│   │   ├── document_handler.go
│   │   ├── indexing_handler.go
│   │   └── search_handler.go
│   │
│   ├── config/                  # Configuration management
│   │   └── config.go
│   │
│   ├── database/                # Database connection & setup
│   │   └── database.go
│   │
│   ├── domain/                  # Domain models (Clean Architecture)
│   │   ├── document.go
│   │   ├── search.go
│   │   └── user.go
│   │
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   │
│   ├── repositories/            # Data access layer
│   │   ├── document_repository.go
│   │   └── user_repository.go
│   │
│   ├── router/                  # Route definitions (Fiber convention)
│   │   └── routes.go
│   │
│   ├── services/                # Business logic layer
│   │   ├── ai_service.go
│   │   ├── auth_service.go
│   │   ├── document_service.go
│   │   ├── indexing_service.go
│   │   └── search_service.go
│   │
│   └── utils/                   # Utility functions
│       ├── jwt.go
│       └── password.go
│
├── pkg/                         # Public library code
│   ├── chromadb/                # ChromaDB client
│   │   └── client.go
│   ├── github/                  # GitHub API client
│   │   └── client.go
│   └── ollama/                  # Ollama LLM client
│       └── client.go
│
├── migrations/                  # Database migrations
│
├── .env.example                 # Environment template
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## คำอธิบายโครงสร้าง

### `cmd/`
- **Purpose**: Main applications สำหรับโปรเจค
- **Convention**: ตาม Go standard project layout
- **Content**: Entry point ของ application (`main.go`)

### `internal/`
- **Purpose**: Private application code ที่ไม่ควร import จากภายนอก
- **Convention**: ตาม Go standard project layout

#### `internal/api/`
- **Purpose**: HTTP handlers สำหรับ Fiber framework
- **Naming**: ใช้ `api` แทน `handlers` ตาม Fiber convention
- **Responsibility**: จัดการ HTTP requests/responses, validation, เรียกใช้ services

#### `internal/router/`
- **Purpose**: Route definitions และ setup
- **Naming**: ใช้ `router` แทน `routes` ตาม Fiber convention
- **Responsibility**: กำหนด routes, middleware groups, route handlers

#### `internal/domain/`
- **Purpose**: Domain models และ entities
- **Naming**: ใช้ `domain` แทน `models` ตาม Clean Architecture
- **Responsibility**: กำหนดโครงสร้างข้อมูลหลักของระบบ

#### `internal/services/`
- **Purpose**: Business logic layer
- **Responsibility**: ตรรกะทางธุรกิจ, orchestration, validation

#### `internal/repositories/`
- **Purpose**: Data access layer
- **Responsibility**: Database operations, queries, data persistence

#### `internal/middleware/`
- **Purpose**: HTTP middleware functions
- **Responsibility**: Authentication, CORS, logging, request processing

#### `internal/config/`
- **Purpose**: Configuration management
- **Responsibility**: Load environment variables, application settings

#### `internal/database/`
- **Purpose**: Database connection and setup
- **Responsibility**: Database connection, migrations, connection pool

#### `internal/utils/`
- **Purpose**: Utility functions
- **Responsibility**: Helper functions (JWT, password hashing, etc.)

### `pkg/`
- **Purpose**: Public library code ที่สามารถ import จากภายนอกได้
- **Convention**: ตาม Go standard project layout
- **Content**: External service clients (Ollama, ChromaDB, GitHub)

## Naming Conventions

### เปลี่ยนชื่อจาก:
- `handlers/` → `api/` (Fiber convention)
- `routes/` → `router/` (Fiber convention)
- `models/` → `domain/` (Clean Architecture)

### เหตุผล:
1. **`api/`**: ตรงกับ Fiber documentation และ naming convention
2. **`router/`**: ตรงกับ Fiber router pattern
3. **`domain/`**: ตรงกับ Clean Architecture และ DDD principles

## References

- [Go Standard Project Layout](https://github.com/golang-standards/project-layout)
- [Fiber Documentation](https://docs.gofiber.io/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
