# Architecture Overview

## โครงสร้างโปรเจค

โปรเจคนี้ใช้ Clean Architecture pattern แบ่งเป็น layers ต่างๆ ดังนี้:

### 1. Presentation Layer (`internal/handlers`)
- จัดการ HTTP requests/responses
- Validate input
- เรียกใช้ service layer

### 2. Business Logic Layer (`internal/services`)
- **AuthService**: Authentication & Authorization
- **DocumentService**: Document management logic
- **SearchService**: Search functionality with AI assistance
- **AIService**: AI/LLM operations
- **IndexingService**: Content indexing from GitHub

### 3. Data Access Layer (`internal/repositories`)
- **UserRepository**: User data operations
- **DocumentRepository**: Document data operations
- ใช้ GORM สำหรับ database operations

### 4. Models (`internal/models`)
- **User, Role**: User management
- **Document, DocumentVersion, DocumentPermission**: Document management
- **SearchQuery, SearchResult**: Search models

### 5. External Services (`pkg/`)
- **ollama**: Ollama LLM client
- **chromadb**: ChromaDB vector database client
- **github**: GitHub API client for wiki.js integration

### 6. Infrastructure
- **config**: Configuration management
- **database**: Database connection & migrations
- **middleware**: HTTP middleware (auth, CORS, logger)
- **utils**: Utility functions (JWT, password hashing)

## Workflow

### Search Flow with Clarification

1. **User sends question** → Chat UI → Backend API
2. **Clarification Phase**:
   - Backend analyzes question clarity using Ollama
   - Searches ChromaDB for related concepts
   - Generates clarifying question if needed
   - Returns clarification request to user
3. **User clarifies** → Backend receives clarified question
4. **Answer Generation**:
   - Create query embedding
   - Retrieve top-K context from ChromaDB
   - Send context + question to Ollama
   - Return generated answer

### Indexing Flow

1. **Load content from GitHub** (wiki.js repo)
2. **Process markdown files**
3. **Create embeddings** (via Ollama or embedding model)
4. **Store in ChromaDB** with metadata
5. **Incremental updates** detect new/changed files

## Database Schema

ตาม ERD ที่ออกแบบไว้:
- `users` - User accounts
- `roles` - User roles (Superadmin, Admin, Guest)
- `user_roles` - Many-to-many relationship
- `documents` - Document metadata
- `document_versions` - Version history
- `document_permissions` - Access control

## API Endpoints

### Public
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login

### Protected (requires auth)
- `GET /api/documents/:id` - Get document
- `POST /api/documents` - Create document
- `PUT /api/documents/:id` - Update document
- `DELETE /api/documents/:id` - Delete document
- `POST /api/documents/:id/versions` - Add version
- `POST /api/search/analyze` - Analyze question clarity
- `GET /api/search?q=...` - Search documents
- `GET /api/search/context?q=...` - Search with AI-generated answer

### Admin Only
- `POST /api/indexing/load-github` - Load content from GitHub
- `POST /api/indexing/index-all` - Index all content

## Configuration

ใช้ environment variables (`.env` file):
- Database connection
- JWT secret
- Ollama URL & model
- ChromaDB URL & collection
- GitHub token & repo info

## Next Steps (Week 1-4)

### Week 1: Foundation ✅
- [x] Project structure
- [x] Database models
- [x] Basic API endpoints
- [ ] GitHub integration testing
- [ ] Ollama setup verification

### Week 2: Indexing & Semantic Search
- [ ] Implement embedding generation
- [ ] ChromaDB collection setup
- [ ] Incremental indexing logic
- [ ] Hybrid search (keyword + semantic)

### Week 3: AI-Assisted Search
- [ ] Query understanding & expansion
- [ ] Semantic matching improvements
- [ ] Result ranking with AI
- [ ] Clarification flow completion

### Week 4: UX & Performance
- [ ] Result snippets
- [ ] Keyword highlighting
- [ ] Performance optimization
- [ ] End-to-end testing
