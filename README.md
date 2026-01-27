# GoFiver - Go + Fiber + React Blog Platform

A modern, enterprise-grade blog platform built with **Go** (Fiber framework) and **React** (TypeScript). 
Designed for Laravel developers transitioning to Go.

## 🚀 Tech Stack

### Backend
- **Go 1.24** - Fast, compiled language
- **Fiber v2** - Express-like web framework (fastest in Go)
- **GORM** - Go ORM (like Eloquent)
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
│       ├── middleware/        # Auth middleware
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
- **API**: http://localhost:8090
- **Frontend**: http://localhost:3000

### 3. Demo Accounts
```
admin@example.com / password123
john@example.com / password123
jane@example.com / password123
```

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
| POST | `/api/auth/register` | Register |
| POST | `/api/auth/login` | Login |
| GET | `/api/blogs` | List blogs |
| GET | `/api/blogs/:id` | Get blog |
| GET | `/api/users/:id/blogs` | User's blogs |

### Protected Endpoints (JWT Required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/me` | Current user |
| GET | `/api/users` | List users |
| PUT | `/api/users/:id` | Update user |
| DELETE | `/api/users/:id` | Delete user |
| POST | `/api/blogs` | Create blog |
| PUT | `/api/blogs/:id` | Update blog |
| DELETE | `/api/blogs/:id` | Delete blog |
| GET | `/api/blogs/my` | My blogs |

## 🔧 Development

### Hot Reload
Both backend and frontend support hot reload:
- Backend: Changes auto-reload via Go run
- Frontend: Vite HMR

### Adding New Features

1. **Add Model** → `backend/internal/models/`
2. **Add Repository** → `backend/internal/repositories/`
3. **Add Service** → `backend/internal/services/`
4. **Add Controller** → `backend/internal/controllers/`
5. **Add Routes** → `backend/internal/routes/routes.go`

## 🧪 Testing API

```bash
# Health check
curl http://localhost:8090/api/health

# Login
curl -X POST http://localhost:8090/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'

# Get blogs
curl http://localhost:8090/api/blogs
```

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

## 🎯 Why Fiber?

- **Express-like syntax** - Familiar for JS/Laravel devs
- **Fastest Go framework** - Built on fasthttp
- **Zero memory allocation** - Highly optimized
- **Middleware support** - Like Laravel middleware
- **Easy to learn** - Simple API

## 📝 License

MIT
