# New Carmen Project

โปรเจค New Carmen Web - ระบบ Document Management พร้อม AI-Powered Search

## โครงสร้างโปรเจค

```
New-carmen/
├── backend/          # Backend API (Go Fiber)
│   ├── cmd/         # Application entry point
│   ├── internal/    # Internal packages
│   ├── pkg/         # External service clients
│   └── ...
└── frontend/        # Frontend (Next.js) - ยังไม่สร้าง
```

## Quick Start

### Backend

```bash
# 1. เข้าไปในโฟลเดอร์ backend
cd backend

# 2. ติดตั้ง dependencies
go mod download

# 3. สร้างไฟล์ .env และแก้ไขค่า
# Windows: notepad .env
# Linux/Mac: nano .env

# 4. สร้าง database
psql -U postgres -c "CREATE DATABASE carmen_db;"

# 5. รันระบบ
go run cmd/server/main.go
```

**📖 ดูคู่มือการรันแบบละเอียด**: [backend/RUN_GUIDE.md](backend/RUN_GUIDE.md)

## Development Plan

- **Week 1**: Foundation & Content Management Flow
- **Week 2**: Indexing & Semantic Search
- **Week 3**: AI-Assisted Search Intelligence
- **Week 4**: UX, Performance & Delivery

## Documentation

- [Backend Architecture](backend/ARCHITECTURE.md)
- [Setup Guide](backend/SETUP.md)
- [Database Setup](backend/DATABASE.md)
- [Next Steps](backend/NEXT_STEPS.md)
