# Database Setup Guide

## Database Type
**PostgreSQL** - ตามที่กำหนดใน ERD

## Connection Details

การเชื่อมต่อฐานข้อมูลอยู่ที่: `internal/database/database.go`

ใช้ **GORM** เป็น ORM และ **PostgreSQL driver** (`gorm.io/driver/postgres`)

## Database Schema

ตาม ERD ที่ออกแบบไว้:

### 1. `users` table
- `id` (bigint, primary key)
- `email` (varchar, unique, not null)
- `password_hash` (varchar, not null)
- `name` (varchar)
- `status` (varchar)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `deleted_at` (timestamp, soft delete)

### 2. `roles` table
- `id` (int, primary key)
- `name` (varchar, unique, not null)

### 3. `user_roles` table (junction table)
- `user_id` (bigint, primary key, foreign key → users.id)
- `role_id` (int, primary key, foreign key → roles.id)

### 4. `documents` table
- `id` (bigint, primary key)
- `title` (varchar, not null)
- `description` (text)
- `owner_id` (bigint, foreign key → users.id)
- `status` (varchar)
- `is_public` (boolean)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `deleted_at` (timestamp, soft delete)

### 5. `document_versions` table
- `id` (bigint, primary key)
- `document_id` (bigint, foreign key → documents.id)
- `version` (int, not null)
- `content` (text)
- `content_html` (text)
- `file_path` (varchar)
- `created_by` (bigint, foreign key → users.id)
- `created_at` (timestamp)

### 6. `document_permissions` table (junction table)
- `document_id` (bigint, primary key, foreign key → documents.id)
- `user_id` (bigint, primary key, foreign key → users.id)
- `permission` (varchar, not null) - values: "read", "write", "admin"

## Migration

ใช้ **GORM AutoMigrate** - จะสร้าง/อัปเดตตารางอัตโนมัติเมื่อ start server

Migration code อยู่ที่: `cmd/server/main.go`

```go
database.Migrate(
    &models.User{},
    &models.Role{},
    &models.Document{},
    &models.DocumentVersion{},
    &models.DocumentPermission{},
)
```

## Setup Steps

### 1. Install PostgreSQL
```bash
# Windows (using Chocolatey)
choco install postgresql

# Linux (Ubuntu/Debian)
sudo apt-get install postgresql postgresql-contrib

# Mac (using Homebrew)
brew install postgresql
```

### 2. Create Database
```sql
-- Connect to PostgreSQL
psql -U postgres

-- Create database
CREATE DATABASE carmen_db;

-- Create user (optional)
CREATE USER carmen_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE carmen_db TO carmen_user;
```

### 3. Configure .env
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=carmen_db
DB_SSLMODE=disable
```

### 4. Run Application
```bash
go run cmd/server/main.go
```

Application จะ:
1. Load configuration จาก `.env`
2. Connect ไปยัง PostgreSQL
3. Auto-migrate tables (สร้างตารางอัตโนมัติ)
4. Start server

## Connection String Format

```
host=localhost user=postgres password=postgres dbname=carmen_db port=5432 sslmode=disable
```

## Roles Setup

หลังจากสร้าง database แล้ว ควร seed roles เริ่มต้น:

```sql
INSERT INTO roles (name) VALUES ('superadmin');
INSERT INTO roles (name) VALUES ('admin');
INSERT INTO roles (name) VALUES ('guest');
```

หรือสร้าง API endpoint สำหรับ seed data ในอนาคต

## Verification

ตรวจสอบว่าเชื่อมต่อได้:

```bash
# Test connection
psql -h localhost -U postgres -d carmen_db -c "SELECT version();"

# Check tables
psql -h localhost -U postgres -d carmen_db -c "\dt"
```

## Troubleshooting

### Connection Refused
- ตรวจสอบว่า PostgreSQL service กำลังรัน
- ตรวจสอบ port (default: 5432)
- ตรวจสอบ firewall settings

### Authentication Failed
- ตรวจสอบ username/password ใน `.env`
- ตรวจสอบ `pg_hba.conf` สำหรับ authentication method

### Database Does Not Exist
- สร้าง database ก่อน: `CREATE DATABASE carmen_db;`

### Migration Errors
- ตรวจสอบว่า user มีสิทธิ์สร้างตาราง
- ตรวจสอบ logs สำหรับ error messages
