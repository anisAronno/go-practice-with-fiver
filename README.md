# GoFiver - Go + Fiber + React Blog Platform

A modern, enterprise-grade blog platform built with **Go** (Fiber framework) and **React** (TypeScript). 
Features user roles, soft delete/restore, image uploads, and admin dashboard.

## 🚀 Tech Stack

### Backend
- **Go 1.24** - Fast, compiled language
- **Fiber v2** - Express-like web framework (fastest in Go)
- **GORM** - Go ORM with soft delete support
- **JWT** - Authentication
- **MySQL** - Database
- **Redis** - Cache

### Frontend
- **React 18** with TypeScript
- **Vite** - Fast build tool
- **TailwindCSS** - Utility-first CSS
- **Zustand** - State management
- **React Router** - Routing
- **Axios** - HTTP client
- **Lucide React** - Icons

## 📁 Project Structure

```
gofiver/
├── backend/                    # Go API
│   ├── cmd/api/main.go        # Entry point
│   └── internal/
│       ├── bootstrap/         # App initialization
│       ├── config/            # Configuration
│       ├── controllers/       # Request handlers
│       ├── database/          # DB & Redis
│       │   └── seeders/       # Data seeders
│       ├── dto/               # Request/Response types
│       ├── middleware/        # Auth & Admin middleware
│       ├── models/            # GORM models
│       ├── repositories/      # Data access layer
│       ├── response/          # API responses
│       ├── routes/            # Route definitions
│       └── services/          # Business logic
├── frontend/                   # React App
│   └── src/
│       ├── components/        # UI components
│       ├── pages/             # Page components
│       ├── services/          # API services
│       ├── store/             # State management
│       └── types/             # TypeScript types
├── docker/                     # Docker files
├── docker-compose.yml
├── Makefile                    # Easy commands
└── .env
```

## ✨ Features

### User Roles
- **Admin**: Full access - manage all blogs, users, trash, restore
- **Author**: Create/edit/delete own blogs only

### Dashboard (Admin)
- View all blogs with author info
- Manage all users (toggle roles, delete)
- Trashed blogs tab (restore/force delete)
- Trashed users tab (restore/force delete)
- Confirmation modals for all destructive actions

### Dashboard (Author)
- View own blogs only
- Create, edit, delete own blogs

### Blog Images
- Upload feature image for each blog
- Supports JPG, PNG, GIF, WebP formats
- Images served via nginx with caching
- Persistent storage via Docker volume

### Soft Delete & Restore
- Delete moves items to trash
- Restore from trash
- Permanent delete (force delete)
- Works for both blogs and users

## 🏃 Quick Start

### Prerequisites
- Docker & Docker Compose
- Common Docker network (`common-net`) with MySQL and Redis

### 1. Start Services
```bash
# Start everything
make up

# Or individually
docker compose up -d
```

### 2. Access
- **App**: http://localhost:8088
- **API**: http://localhost:8088/api

### 3. Demo Accounts
| Email | Password | Role |
|-------|----------|------|
| admin@example.com | password123 | Admin |
| john@example.com | password123 | Author |
| jane@example.com | password123 | Author |

## 🛠 Commands

```bash
make up              # Start all services
make down            # Stop all services
make build           # Build containers
make logs            # View all logs
make api-logs        # View API logs
make shell-api       # Open shell in API container
make restart         # Restart services
make fresh           # Rebuild everything
```

## 📡 API Endpoints

### Public Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login user |
| GET | `/api/blogs` | List all blogs |
| GET | `/api/blogs/:id` | Get single blog |
| GET | `/api/users/:id/blogs` | Get user's blogs |

### Protected Endpoints (JWT Required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/me` | Current user profile |
| GET | `/api/users` | List all users |
| GET | `/api/users/:id` | Get user details |
| PUT | `/api/users/:id` | Update user |
| PATCH | `/api/users/:id/role` | Update user role (admin only) |
| DELETE | `/api/users/:id` | Soft delete user |
| POST | `/api/blogs` | Create blog |
| GET | `/api/blogs/my` | Get current user's blogs |
| PUT | `/api/blogs/:id` | Update blog |
| DELETE | `/api/blogs/:id` | Soft delete blog |
| POST | `/api/blogs/:id/image` | Upload blog image |

### Admin Endpoints (Admin Role Required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/blogs/trashed` | List trashed blogs |
| POST | `/api/admin/blogs/:id/restore` | Restore trashed blog |
| DELETE | `/api/admin/blogs/:id/force` | Permanently delete blog |
| GET | `/api/admin/users/trashed` | List trashed users |
| POST | `/api/admin/users/:id/restore` | Restore trashed user |
| DELETE | `/api/admin/users/:id/force` | Permanently delete user |

## 🧪 Testing API

```bash
# Health check
curl http://localhost:8088/api/health

# Login as admin
curl -X POST http://localhost:8088/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'

# Create blog
curl -X POST http://localhost:8088/api/blogs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title":"My Blog","content":"Content here..."}'

# Upload blog image
curl -X POST http://localhost:8088/api/blogs/1/image \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@/path/to/image.jpg"

# Get trashed blogs (admin only)
curl http://localhost:8088/api/admin/blogs/trashed \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Restore a blog (admin only)
curl -X POST http://localhost:8088/api/admin/blogs/1/restore \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

## 🗄️ Database Schema

### Users Table
| Column | Type | Description |
|--------|------|-------------|
| id | uint | Primary key |
| name | string | User name |
| email | string | Unique email |
| password | string | Hashed password |
| role | string | 'admin' or 'author' |
| created_at | timestamp | Created timestamp |
| updated_at | timestamp | Updated timestamp |
| deleted_at | timestamp | Soft delete timestamp |

### Blogs Table
| Column | Type | Description |
|--------|------|-------------|
| id | uint | Primary key |
| title | string | Blog title |
| content | text | Blog content |
| image | string | Feature image path |
| user_id | uint | Author foreign key |
| created_at | timestamp | Created timestamp |
| updated_at | timestamp | Updated timestamp |
| deleted_at | timestamp | Soft delete timestamp |

## 📦 Architecture (Laravel Comparison)

| Laravel | GoFiver |
|---------|---------|
| Controllers | controllers/ |
| Models | models/ |
| Migrations | GORM AutoMigrate |
| Seeders | database/seeders/ |
| FormRequest | dto/ |
| Services | services/ |
| Repositories | repositories/ |
| Middleware | middleware/ |
| routes/api.php | routes/routes.go |
| config/ | config/ |
| SoftDeletes Trait | GORM DeletedAt |

## 🔒 Middleware

### AuthMiddleware
- Validates JWT token
- Attaches user to request context
- Returns 401 for invalid/missing token

### AdminMiddleware
- Checks user role is 'admin'
- Returns 403 for non-admin users
- Used for trash management routes

## 🎨 Frontend Components

### ConfirmModal
Reusable confirmation modal with types:
- `delete` - Red theme for delete actions
- `restore` - Green theme for restore actions
- `force-delete` - Red theme for permanent delete
- `warning` - Yellow theme for warnings

### LogoutModal
Beautiful logout confirmation modal with escape key support.

## 🚀 Deployment

### Environment Variables
```env
BUILD_TARGET=development
NGINX_PORT=8088
DB_HOST=mysql
DB_PORT=3306
DB_DATABASE=gofiver
DB_USERNAME=root
DB_PASSWORD=your_password
JWT_SECRET=your_secret_key
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password
```

### Production Build
```bash
# Set production target
export BUILD_TARGET=production

# Build and deploy
docker compose build
docker compose up -d
```

## 🎯 Why Fiber?

- **Express-like syntax** - Familiar for JS/Laravel devs
- **Fastest Go framework** - Built on fasthttp
- **Zero memory allocation** - Highly optimized
- **Middleware support** - Like Laravel middleware
- **Easy to learn** - Simple API

## 📝 License

MIT
