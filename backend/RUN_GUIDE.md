# 🚀 คู่มือการรันระบบ New Carmen Backend

คู่มือนี้จะบอกขั้นตอนการรันระบบตั้งแต่เริ่มต้น พร้อมบอก path ที่ต้อง cd ไป

---

## 📍 ตำแหน่งโปรเจค

```
D:\New-carmen\
└── backend\          ← ต้อง cd มาที่นี่
    ├── cmd\
    ├── internal\
    ├── pkg\
    └── go.mod
```

---

## ✅ ขั้นตอนที่ 1: ตรวจสอบ Prerequisites

### 1.1 ตรวจสอบ Go

เปิด Command Prompt หรือ Terminal แล้วรัน:

```bash
go version
```

**ต้องได้**: `go version go1.21.x` หรือสูงกว่า

**ถ้ายังไม่มี**:
- Download: https://golang.org/dl/
- ติดตั้งตามปกติ
- ตรวจสอบ PATH environment variable

### 1.2 ตรวจสอบ PostgreSQL

```bash
psql --version
```

**ต้องได้**: `psql (PostgreSQL) 12.x` หรือสูงกว่า

**ถ้ายังไม่มี**:
- Windows: Download จาก https://www.postgresql.org/download/windows/
- หรือใช้ Chocolatey: `choco install postgresql`

---

## 📂 ขั้นตอนที่ 2: เข้าไปในโฟลเดอร์ Backend

เปิด Command Prompt หรือ Terminal แล้วรัน:

```bash
# ถ้าคุณอยู่ที่ D:\New-carmen
cd backend

# หรือถ้าคุณอยู่ที่อื่น ให้ cd ไปที่ path เต็ม
cd D:\New-carmen\backend
```

**ตรวจสอบว่าอยู่ที่ถูกต้อง**:
```bash
# Windows
dir

# Linux/Mac
ls

# ต้องเห็นไฟล์ go.mod, Makefile, cmd/, internal/, pkg/
```

---

## 📦 ขั้นตอนที่ 3: ติดตั้ง Dependencies

เมื่ออยู่ที่ `D:\New-carmen\backend` แล้ว รัน:

```bash
go mod download
```

**หรือใช้ Make**:
```bash
make deps
```

**ผลลัพธ์ที่ควรเห็น**:
```
go: downloading github.com/gofiber/fiber/v2 v2.52.0
go: downloading gorm.io/gorm v1.25.5
...
```

---

## ⚙️ ขั้นตอนที่ 4: สร้างไฟล์ .env

### 4.1 สร้างไฟล์ .env

```bash
# Windows (PowerShell)
New-Item -Path .env -ItemType File

# Windows (CMD)
type nul > .env

# Linux/Mac
touch .env
```

### 4.2 เปิดไฟล์ .env ด้วย text editor

```bash
# Windows
notepad .env

# Linux/Mac
nano .env
# หรือ
code .env
```

### 4.3 ใส่ค่าตามนี้ (แก้ไขให้ตรงกับเครื่องคุณ):

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost
ENVIRONMENT=development

# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=carmen_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production-min-32-chars-please-change-this
JWT_EXPIRY=24h

# Ollama Configuration
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama2

# ChromaDB Configuration
CHROMADB_URL=http://localhost:8000
CHROMADB_COLLECTION=carmen_documents

# GitHub Configuration (ถ้ายังไม่มี ให้เว้นว่างไว้ก่อน)
GITHUB_TOKEN=
GITHUB_REPO_OWNER=
GITHUB_REPO_NAME=
GITHUB_BRANCH=main

# Git Configuration
GIT_REPO_PATH=./wiki-content
GIT_REPO_URL=
```

**สำคัญ**: แก้ไข `DB_PASSWORD` ให้ตรงกับ password ของ PostgreSQL ที่คุณตั้งไว้

---

## 🗄️ ขั้นตอนที่ 5: Setup PostgreSQL Database

### 5.1 ตรวจสอบว่า PostgreSQL รันอยู่

```bash
# Windows - ตรวจสอบใน Services
# กด Win+R พิมพ์ services.msc
# หา "postgresql-x64-xx" และตรวจสอบว่า Running

# หรือใช้ Command Line
pg_ctl status
```

**ถ้ายังไม่รัน**:
```bash
# Windows
# ไปที่ Services และ Start PostgreSQL service
# หรือ
pg_ctl start -D "C:\Program Files\PostgreSQL\15\data"
```

### 5.2 สร้าง Database

```bash
# เข้า PostgreSQL
psql -U postgres

# ถ้าเจอ password prompt ให้ใส่ password ที่ตั้งไว้
```

**ใน psql prompt**:
```sql
-- สร้าง database
CREATE DATABASE carmen_db;

-- ตรวจสอบว่าสร้างสำเร็จ
\l

-- ออกจาก psql
\q
```

---

## 🐳 ขั้นตอนที่ 6: Setup Ollama (Optional - สำหรับ AI Search)

### 6.1 ติดตั้ง Ollama

- Download: https://ollama.ai/download
- ติดตั้งตามปกติ

### 6.2 Start Ollama

เปิด Terminal ใหม่ (ไม่ต้อง cd ไปไหน):

```bash
ollama serve
```

**ผลลัพธ์**: `Ollama is running on http://localhost:11434`

### 6.3 Download Model (ใน Terminal ใหม่)

```bash
ollama pull llama2
```

**หมายเหตุ**: ถ้ายังไม่ต้องการใช้ AI Search สามารถข้ามขั้นตอนนี้ได้

---

## 🗄️ ขั้นตอนที่ 7: Setup ChromaDB (Optional - สำหรับ Vector Search)

### 7.1 ใช้ Docker (แนะนำ)

```bash
docker run -d -p 8000:8000 chromadb/chroma
```

**ตรวจสอบ**:
```bash
curl http://localhost:8000/api/v1/heartbeat
```

**หมายเหตุ**: ถ้ายังไม่ต้องการใช้ Vector Search สามารถข้ามขั้นตอนนี้ได้

---

## 🚀 ขั้นตอนที่ 8: รันระบบ

### 8.1 กลับไปที่โฟลเดอร์ backend

```bash
cd D:\New-carmen\backend
```

### 8.2 รันระบบ

**วิธีที่ 1: ใช้ go run (แนะนำสำหรับ development)**
```bash
go run cmd/server/main.go
```

**วิธีที่ 2: ใช้ Make**
```bash
make run
```

**ผลลัพธ์ที่ควรเห็น**:
```
Database connected successfully
Database migration completed
Server starting on port 8080
```

**ถ้าเห็นข้อความนี้ = ระบบรันสำเร็จ! ✅**

---

## 🧪 ขั้นตอนที่ 9: ทดสอบระบบ

### 9.1 Health Check

เปิด Terminal ใหม่ (หรือใช้ browser):

```bash
# Windows (PowerShell)
Invoke-WebRequest -Uri http://localhost:8080/health

# Linux/Mac
curl http://localhost:8080/health
```

**หรือเปิด browser**: http://localhost:8080/health

**ผลลัพธ์ที่ควรเห็น**:
```json
{"status":"ok"}
```

### 9.2 ทดสอบ Register

```bash
# Windows (PowerShell)
Invoke-WebRequest -Method POST -Uri http://localhost:8080/api/auth/register `
  -ContentType "application/json" `
  -Body '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Linux/Mac
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```

**ผลลัพธ์ที่ควรเห็น**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "test@example.com",
    "name": "Test User"
  }
}
```

### 9.3 ทดสอบ Login

```bash
# Windows (PowerShell)
Invoke-WebRequest -Method POST -Uri http://localhost:8080/api/auth/login `
  -ContentType "application/json" `
  -Body '{"email":"test@example.com","password":"password123"}'

# Linux/Mac
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

---

## 📝 สรุปขั้นตอนทั้งหมด

```bash
# 1. ตรวจสอบ Go และ PostgreSQL
go version
psql --version

# 2. เข้าไปในโฟลเดอร์ backend
cd D:\New-carmen\backend

# 3. ติดตั้ง dependencies
go mod download

# 4. สร้างไฟล์ .env และแก้ไขค่า
notepad .env

# 5. สร้าง database
psql -U postgres
CREATE DATABASE carmen_db;
\q

# 6. (Optional) Start Ollama
ollama serve

# 7. (Optional) Start ChromaDB
docker run -d -p 8000:8000 chromadb/chroma

# 8. รันระบบ
go run cmd/server/main.go

# 9. ทดสอบ
curl http://localhost:8080/health
```

---

## ⚠️ Troubleshooting

### Error: "Failed to connect to database"

**สาเหตุ**: PostgreSQL ไม่ได้รัน หรือ credentials ผิด

**แก้ไข**:
1. ตรวจสอบว่า PostgreSQL service รันอยู่
2. ตรวจสอบ `DB_USER`, `DB_PASSWORD`, `DB_NAME` ใน `.env`
3. ทดสอบ connection:
   ```bash
   psql -U postgres -d carmen_db
   ```

### Error: "No .env file found"

**สาเหตุ**: ไฟล์ `.env` ไม่มี หรือไม่ได้อยู่ในโฟลเดอร์ `backend`

**แก้ไข**:
```bash
# ตรวจสอบว่าอยู่ที่ backend directory
cd D:\New-carmen\backend

# ตรวจสอบว่ามีไฟล์ .env
dir .env

# ถ้าไม่มี ให้สร้างใหม่
notepad .env
```

### Error: "module not found"

**สาเหตุ**: Dependencies ยังไม่ได้ download

**แก้ไข**:
```bash
cd D:\New-carmen\backend
go mod download
go mod tidy
```

### Error: "port 8080 already in use"

**สาเหตุ**: Port 8080 ถูกใช้งานอยู่

**แก้ไข**:
1. หา process ที่ใช้ port 8080:
   ```bash
   # Windows
   netstat -ano | findstr :8080
   
   # Linux/Mac
   lsof -i :8080
   ```
2. Kill process หรือเปลี่ยน port ใน `.env`:
   ```env
   SERVER_PORT=8081
   ```

---

## 🎯 Quick Start (สรุปสั้นๆ)

```bash
# 1. cd ไปที่ backend
cd D:\New-carmen\backend

# 2. ติดตั้ง dependencies
go mod download

# 3. สร้าง .env และแก้ไขค่า
notepad .env

# 4. สร้าง database
psql -U postgres -c "CREATE DATABASE carmen_db;"

# 5. รันระบบ
go run cmd/server/main.go
```

---

## 📚 เอกสารเพิ่มเติม

- [SETUP.md](SETUP.md) - Setup guide แบบละเอียด
- [DATABASE.md](DATABASE.md) - Database setup
- [WORKFLOW.md](WORKFLOW.md) - การทำงานของระบบ
- [ARCHITECTURE.md](ARCHITECTURE.md) - Architecture overview

---

## ✅ Checklist

ก่อนรันระบบ ตรวจสอบว่า:

- [ ] Go 1.21+ ติดตั้งแล้ว
- [ ] PostgreSQL ติดตั้งและรันอยู่
- [ ] Database `carmen_db` สร้างแล้ว
- [ ] ไฟล์ `.env` สร้างและแก้ไขค่าแล้ว
- [ ] Dependencies download แล้ว (`go mod download`)
- [ ] อยู่ที่โฟลเดอร์ `D:\New-carmen\backend`
- [ ] (Optional) Ollama รันอยู่ (ถ้าต้องการ AI Search)
- [ ] (Optional) ChromaDB รันอยู่ (ถ้าต้องการ Vector Search)

**พร้อมแล้ว! รัน `go run cmd/server/main.go` ได้เลย** 🚀
