# ระบบการทำงานของ New Carmen Backend

## 📋 สารบัญ
1. [ภาพรวมระบบ](#ภาพรวมระบบ)
2. [Architecture Layers](#architecture-layers)
3. [Request Flow](#request-flow)
4. [Authentication Flow](#authentication-flow)
5. [Document Management Flow](#document-management-flow)
6. [AI-Powered Search Flow](#ai-powered-search-flow)
7. [Indexing Flow](#indexing-flow)

---

## ภาพรวมระบบ

New Carmen เป็นระบบ Document Management ที่ใช้ AI ช่วยในการค้นหา โดยมีส่วนประกอบหลัก:

- **Backend API** (Go Fiber) - RESTful API
- **PostgreSQL** - Relational database สำหรับเก็บ metadata
- **ChromaDB** - Vector database สำหรับ semantic search
- **Ollama** - LLM สำหรับ AI processing
- **GitHub** - เก็บ markdown files จาก wiki.js

---

## Architecture Layers

ระบบใช้ **Clean Architecture** แบ่งเป็น 4 layers:

```
┌─────────────────────────────────────┐
│   Presentation Layer (API)         │  ← HTTP Handlers
├─────────────────────────────────────┤
│   Business Logic Layer (Services)   │  ← Business Rules
├─────────────────────────────────────┤
│   Data Access Layer (Repositories)  │  ← Database Operations
├─────────────────────────────────────┤
│   Infrastructure                    │  ← DB, Config, Utils
└─────────────────────────────────────┘
```

### 1. Presentation Layer (`internal/api/`)
**หน้าที่**: จัดการ HTTP requests/responses

- รับ request จาก client
- Validate input data
- เรียกใช้ Service layer
- ส่ง response กลับ

**ไฟล์**:
- `auth_handler.go` - Authentication endpoints
- `document_handler.go` - Document CRUD
- `search_handler.go` - Search endpoints
- `indexing_handler.go` - Indexing endpoints

### 2. Business Logic Layer (`internal/services/`)
**หน้าที่**: ตรรกะทางธุรกิจ

- **AuthService**: จัดการ authentication, authorization
- **DocumentService**: จัดการ documents, permissions, versions
- **SearchService**: Hybrid search (keyword + semantic), AI clarification
- **IndexingService**: Index content จาก GitHub

### 3. Data Access Layer (`internal/repositories/`)
**หน้าที่**: Database operations

- **UserRepository**: CRUD operations สำหรับ users
- **DocumentRepository**: CRUD operations สำหรับ documents
- ใช้ GORM เป็น ORM

### 4. Infrastructure
- **Database**: PostgreSQL connection & migrations
- **Config**: Environment variables
- **Middleware**: Auth, CORS, Logger
- **Utils**: JWT, Password hashing

---

## Request Flow

### ขั้นตอนการทำงานของ Request

```
1. Client Request
   ↓
2. Fiber Router (routes.go)
   ↓
3. Middleware (Logger, CORS, Auth)
   ↓
4. API Handler (internal/api/)
   ↓
5. Service Layer (internal/services/)
   ↓
6. Repository Layer (internal/repositories/)
   ↓
7. Database (PostgreSQL)
   ↓
8. Response กลับไปยัง Client
```

### ตัวอย่าง: การสร้าง Document

```
1. POST /api/documents
   ↓
2. AuthMiddleware ตรวจสอบ JWT token
   ↓
3. DocumentHandler.CreateDocument()
   - Parse request body
   - เรียก DocumentService.CreateDocument()
   ↓
4. DocumentService.CreateDocument()
   - Validate business rules
   - เรียก DocumentRepository.Create()
   ↓
5. DocumentRepository.Create()
   - GORM insert ลง PostgreSQL
   ↓
6. Response กลับไปยัง Client
```

---

## Authentication Flow

### 1. Registration (สมัครสมาชิก)

```
User → POST /api/auth/register
  ↓
AuthHandler.Register()
  ↓
AuthService.Register()
  ├─ ตรวจสอบ email ซ้ำ
  ├─ Hash password (bcrypt)
  ├─ สร้าง User ใน database
  └─ Generate JWT token
  ↓
Response: { token, user }
```

**Code Flow**:
```go
// 1. Handler รับ request
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req services.RegisterRequest
    c.BodyParser(&req)
    
    // 2. เรียก Service
    response, err := h.authService.Register(req)
    
    // 3. ส่ง response
    return c.JSON(response)
}

// 3. Service ทำ business logic
func (s *AuthService) Register(req RegisterRequest) {
    // Check duplicate
    existing, _ := s.userRepo.GetByEmail(req.Email)
    
    // Hash password
    hashedPassword := utils.HashPassword(req.Password)
    
    // Create user
    user := &domain.User{...}
    s.userRepo.Create(user)
    
    // Generate JWT
    token := utils.GenerateToken(user.ID, user.Email, roles)
    
    return &AuthResponse{Token: token, User: user}
}
```

### 2. Login

```
User → POST /api/auth/login
  ↓
AuthHandler.Login()
  ↓
AuthService.Login()
  ├─ หา User จาก email
  ├─ ตรวจสอบ password
  ├─ ตรวจสอบ status (active)
  └─ Generate JWT token
  ↓
Response: { token, user }
```

### 3. Protected Routes

```
User → GET /api/documents/:id
  ↓
AuthMiddleware()
  ├─ อ่าน Authorization header
  ├─ Validate JWT token
  ├─ Extract claims (user_id, email, roles)
  └─ เก็บใน c.Locals()
  ↓
DocumentHandler.GetDocument()
  ├─ อ่าน user_id จาก c.Locals()
  ├─ ตรวจสอบ permission
  └─ เรียก DocumentService
  ↓
Response
```

**Middleware Code**:
```go
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 1. อ่าน token
        tokenString := c.Get("Authorization")
        
        // 2. Validate token
        claims, err := utils.ValidateToken(tokenString)
        
        // 3. เก็บใน context
        c.Locals("user_id", claims.UserID)
        c.Locals("email", claims.Email)
        c.Locals("roles", claims.Roles)
        
        return c.Next()
    }
}
```

---

## Document Management Flow

### 1. Create Document

```
User → POST /api/documents
Body: { title, description, is_public }
  ↓
AuthMiddleware (ตรวจสอบ token)
  ↓
DocumentHandler.CreateDocument()
  ├─ อ่าน user_id จาก c.Locals()
  ├─ Parse request body
  └─ เรียก DocumentService.CreateDocument()
  ↓
DocumentService.CreateDocument()
  ├─ สร้าง Document object
  ├─ Set owner_id = user_id
  └─ เรียก DocumentRepository.Create()
  ↓
DocumentRepository.Create()
  └─ GORM INSERT → PostgreSQL
  ↓
Response: Document object
```

### 2. Get Document (with Permission Check)

```
User → GET /api/documents/:id
  ↓
AuthMiddleware
  ↓
DocumentHandler.GetDocument()
  ├─ Parse document ID
  ├─ อ่าน user_id จาก c.Locals()
  └─ เรียก DocumentService.CheckPermission()
  ↓
DocumentService.CheckPermission()
  ├─ ตรวจสอบว่าเป็น owner หรือไม่
  ├─ ตรวจสอบ is_public
  └─ ตรวจสอบ document_permissions
  ↓
ถ้ามี permission → DocumentService.GetDocument()
  └─ DocumentRepository.GetByID()
  ↓
Response: Document with versions
```

### 3. Add Version

```
User → POST /api/documents/:id/versions
Body: { content, content_html, file_path }
  ↓
DocumentHandler.AddVersion()
  ├─ ตรวจสอบ permission (write)
  └─ เรียก DocumentService.AddVersion()
  ↓
DocumentService.AddVersion()
  ├─ หา latest version
  ├─ เพิ่ม version number
  └─ สร้าง DocumentVersion
  ↓
DocumentRepository.AddVersion()
  └─ INSERT → document_versions table
  ↓
Response: DocumentVersion
```

---

## AI-Powered Search Flow

### 1. Question Clarity Analysis

```
User → POST /api/search/analyze
Body: { question: "..." }
  ↓
SearchHandler.AnalyzeClarity()
  ↓
SearchService.AnalyzeQuestionClarity()
  ├─ เรียก Ollama.AnalyzeQuestionClarity()
  │   └─ LLM วิเคราะห์ว่าคำถามคลุมเครือหรือไม่
  │
  ├─ ถ้าคลุมเครือ:
  │   ├─ Query ChromaDB หา related concepts
  │   ├─ Generate clarifying question (Ollama)
  │   └─ Return clarification options
  │
  └─ ถ้าไม่คลุมเครือ:
      └─ Return { needs_clarification: false }
  ↓
Response: ClarificationResponse
```

**Flow Diagram**:
```
Question → Ollama Analysis
    ↓
Is Ambiguous?
    ├─ Yes → Search ChromaDB
    │         ↓
    │      Generate Clarifying Question
    │         ↓
    │      Return Options
    │
    └─ No → Continue to Search
```

### 2. Hybrid Search (Keyword + Semantic)

```
User → GET /api/search?q=query
  ↓
SearchHandler.Search()
  ↓
SearchService.Search()
  ├─ 1. Keyword Search (PostgreSQL)
  │   └─ DocumentRepository.Search()
  │      └─ SQL: WHERE title ILIKE '%query%'
  │
  ├─ 2. Semantic Search (ChromaDB)
  │   └─ chromadb.Query(query)
  │      └─ Vector similarity search
  │
  └─ 3. Combine & Rank Results
      ├─ ChromaDB results (higher relevance)
      └─ Keyword results (lower relevance)
  ↓
Response: SearchResponse { results, total, query }
```

### 3. Search with AI-Generated Answer (RAG)

```
User → GET /api/search/context?q=query
  ↓
SearchHandler.SearchWithContext()
  ↓
SearchService.SearchWithContext()
  ├─ 1. Retrieve Context from ChromaDB
  │   └─ chromadb.Query(query, limit=5)
  │      └─ Get top-K similar documents
  │
  ├─ 2. Build Context String
  │   └─ Combine documents into context
  │
  └─ 3. Generate Answer (Ollama)
      └─ ollama.GenerateAnswer(context, query)
         └─ LLM สร้างคำตอบจาก context
  ↓
Response: { answer, query }
```

**RAG Flow**:
```
Query → ChromaDB (Vector Search)
    ↓
Retrieve Top-K Context
    ↓
Build Context String
    ↓
Send to Ollama: Context + Question
    ↓
LLM Generate Answer
    ↓
Return Answer
```

---

## Indexing Flow

### 1. Load Content from GitHub

```
Admin → POST /api/indexing/load-github
  ↓
RoleMiddleware (ตรวจสอบ admin role)
  ↓
IndexingHandler.LoadFromGitHub()
  ↓
IndexingService.LoadContentFromGitHub()
  ├─ เรียก GitHub API
  │   └─ github.ListMarkdownFiles()
  │
  ├─ Loop through files
  │   └─ github.GetFileContent(path)
  │
  └─ Return files list
  ↓
Response: { files: [...], count: N }
```

### 2. Index All Content

```
Admin → POST /api/indexing/index-all
  ↓
IndexingHandler.IndexAll()
  ↓
IndexingService.IncrementalIndexing()
  ├─ 1. Load files from GitHub
  │
  ├─ 2. For each file:
  │   ├─ Create embedding (Ollama)
  │   ├─ Store in ChromaDB
  │   └─ Store metadata in PostgreSQL
  │
  └─ 3. Track indexed files
  ↓
Response: { message: "Indexing completed" }
```

**Indexing Process**:
```
GitHub Repo
    ↓
Load Markdown Files
    ↓
For each file:
    ├─ Extract content
    ├─ Create embedding
    ├─ Store in ChromaDB (vector)
    └─ Store metadata in PostgreSQL
    ↓
Done
```

---

## Database Schema

### Tables

1. **users**
   - `id`, `email`, `password_hash`, `name`, `status`
   - `created_at`, `updated_at`

2. **roles**
   - `id`, `name` (superadmin, admin, guest)

3. **user_roles** (junction table)
   - `user_id`, `role_id`

4. **documents**
   - `id`, `title`, `description`, `owner_id`, `status`, `is_public`
   - `created_at`, `updated_at`

5. **document_versions**
   - `id`, `document_id`, `version`, `content`, `content_html`, `file_path`
   - `created_by`, `created_at`

6. **document_permissions** (junction table)
   - `document_id`, `user_id`, `permission` (read, write, admin)

---

## API Endpoints Summary

### Public Endpoints
- `POST /api/auth/register` - สมัครสมาชิก
- `POST /api/auth/login` - เข้าสู่ระบบ
- `GET /health` - Health check

### Protected Endpoints (ต้องมี JWT)
- `GET /api/documents/:id` - ดู document
- `POST /api/documents` - สร้าง document
- `PUT /api/documents/:id` - แก้ไข document
- `DELETE /api/documents/:id` - ลบ document
- `POST /api/documents/:id/versions` - เพิ่ม version
- `POST /api/search/analyze` - วิเคราะห์คำถาม
- `GET /api/search?q=...` - ค้นหา documents
- `GET /api/search/context?q=...` - ค้นหาพร้อม AI answer

### Admin Only
- `POST /api/indexing/load-github` - โหลดไฟล์จาก GitHub
- `POST /api/indexing/index-all` - Index เนื้อหาทั้งหมด

---

## Security

1. **JWT Authentication**
   - Token หมดอายุตาม config (default: 24h)
   - เก็บใน Authorization header: `Bearer <token>`

2. **Password Security**
   - Hash ด้วย bcrypt
   - ไม่เก็บ plain text password

3. **Permission System**
   - Owner มีสิทธิ์ทั้งหมด
   - Public documents อ่านได้ทุกคน
   - Specific permissions ผ่าน document_permissions table

4. **Role-Based Access Control (RBAC)**
   - Superadmin, Admin, Guest
   - Middleware ตรวจสอบ role ก่อนเข้าถึง admin endpoints

---

## Error Handling

ระบบจัดการ errors แบบ layer-by-layer:

```
Repository Error → Service Error → Handler Error → HTTP Response
```

**ตัวอย่าง**:
```go
// Repository
if err := r.db.Create(user).Error; err != nil {
    return err  // Database error
}

// Service
if err := s.userRepo.Create(user); err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}

// Handler
if err != nil {
    return c.Status(500).JSON(fiber.Map{
        "error": err.Error()
    })
}
```

---

## Configuration

ระบบใช้ environment variables (`.env`):

- **Server**: PORT, HOST, ENVIRONMENT
- **Database**: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
- **JWT**: JWT_SECRET, JWT_EXPIRY
- **Ollama**: OLLAMA_URL, OLLAMA_MODEL
- **ChromaDB**: CHROMADB_URL, CHROMADB_COLLECTION
- **GitHub**: GITHUB_TOKEN, GITHUB_REPO_OWNER, GITHUB_REPO_NAME

---

## Next Steps

ดู [ARCHITECTURE.md](ARCHITECTURE.md) สำหรับรายละเอียดเพิ่มเติม
