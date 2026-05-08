# Looking Glass - ISP Network Diagnostics Platform

A full-stack web application for performing network diagnostics (ping, traceroute, BGP lookup) across multiple routers from various vendors (MikroTik, Juniper, Cisco, Huawei).

## Features

- **Multi-Vendor Support**: SSH-based connection to MikroTik, Juniper, Cisco, and Huawei routers
- **Network Diagnostics**: Execute ping, traceroute, and BGP lookup queries through a web interface
- **Real-time Streaming**: Live query results via WebSocket connections
- **Role-Based Access Control**: Admin, Operator, and User roles with appropriate permissions
- **Audit Logging**: Complete audit trail of all user actions
- **Responsive UI**: Modern Next.js frontend with dark/light mode support
- **Caching**: Redis-based caching for improved performance
- **Job Queue**: Background processing for long-running queries

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Client (Browser)                         │
│                   Next.js + React Query                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Nginx Reverse Proxy                       │
│                   (SSL Termination, Routing)                 │
└─────────────────────────────────────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    ▼                   ▼
        ┌───────────────────┐  ┌───────────────────┐
        │   Frontend (3000) │  │  Backend (8080)   │
        │   Next.js Static  │  │   Go REST API     │
        └───────────────────┘  │   + WebSocket     │
                               └───────────────────┘
                                          │
                              ┌───────────┴───────────┐
                              ▼                       ▼
                    ┌─────────────────┐      ┌─────────────────┐
                    │   PostgreSQL    │      │     Redis       │
                    │   (Metadata)    │      │    (Cache)      │
                    └─────────────────┘      └─────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │     Routers     │
                    │  (MikroTik, etc)│
                    └─────────────────┘
```

## Tech Stack

### Backend
- **Go 1.21** - Server language
- **Gin** - HTTP web framework
- **GORM** - ORM for PostgreSQL
- **gorilla/websocket** - WebSocket support
- **ssh** - SSH client for router connections

### Frontend
- **Next.js 14** - React framework
- **TypeScript** - Type safety
- **TailwindCSS** - Styling
- **Tanstack React Query** - Data fetching
- **Zustand** - State management
- **Radix UI** - UI primitives

### Infrastructure
- **PostgreSQL 15** - Primary database
- **Redis 7** - Caching and job queue
- **Docker** - Containerization
- **Nginx** - Reverse proxy

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Git

### 1. Clone the repository
```bash
git clone https://github.com/your-org/looking-glass.git
cd looking-glass
```

### 2. Configure environment
```bash
cp backend/.env.example backend/.env
cp frontend/.env.local.example frontend/.env.local
```

Edit the `.env` files with your configuration.

### 3. Start with Docker
```bash
docker-compose up -d
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

### 4. Default Login
- **Username**: admin
- **Password**: admin123

> ⚠️ Change the default password immediately in production!

## Manual Setup

### Backend

```bash
cd backend

# Install dependencies
go mod download

# Run database migrations
go run cmd/server/main.go migrate

# Start the server
go run cmd/server/main.go serve
```

### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.local.example .env.local

# Start development server
npm run dev
```

## Configuration

### Environment Variables

#### Backend (`backend/.env`)
| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Server port | `8080` |
| `APP_SECRET` | Application secret | `changeme` |
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `REDIS_URL` | Redis connection string | Required |
| `JWT_SECRET` | JWT signing secret | Required |
| `JWT_EXPIRY` | Token expiration time | `24h` |

#### Frontend (`frontend/.env.local`)
| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:8080` |
| `NEXT_PUBLIC_WS_URL` | WebSocket URL | `ws://localhost:8080` |

## API Documentation

### Authentication

```
POST /api/auth/login
POST /api/auth/register
POST /api/auth/refresh
GET  /api/auth/me
```

### Routers

```
GET    /api/routers         - List all routers
POST   /api/routers         - Create a router
GET    /api/routers/:id     - Get router details
PUT    /api/routers/:id     - Update router
DELETE /api/routers/:id     - Delete router
POST   /api/routers/:id/test - Test router connection
GET    /api/routers/:id/status - Get router status
```

### Zones

```
GET    /api/zones           - List all zones
POST   /api/zones           - Create a zone
GET    /api/zones/:id       - Get zone details
PUT    /api/zones/:id       - Update zone
DELETE /api/zones/:id       - Delete zone
```

### Queries

```
POST /api/queries/ping       - Execute ping
POST /api/queries/traceroute - Execute traceroute
POST /api/queries/bgp        - Execute BGP lookup
GET  /api/queries           - List recent queries
GET  /api/queries/:id       - Get query results
```

### WebSocket

```
WS /ws?q=ping&router_id=xxx&target=8.8.8.8
```

## Development

### Running Tests

```bash
# Backend
cd backend
go test -v ./...

# Frontend
cd frontend
npm run test
```

### Code Style

```bash
# Backend (Go)
cd backend
go fmt ./...
go vet ./...

# Frontend (ESLint)
cd frontend
npm run lint
```

## Project Structure

```
looking-glass/
├── backend/
│   ├── cmd/
│   │   └── server/         # Main application entry
│   ├── internal/
│   │   ├── api/            # HTTP handlers and routes
│   │   │   ├── handlers/   # Request handlers
│   │   │   ├── middleware/ # HTTP middleware
│   │   │   └── websocket/  # WebSocket handler
│   │   ├── services/      # Business logic
│   │   ├── repository/     # Data access layer
│   │   ├── vendor/         # Router vendor implementations
│   │   ├── ssh/           # SSH client management
│   │   ├── queue/         # Job queue
│   │   ├── database/      # Database connection
│   │   └── cache/         # Redis cache
│   └── go.mod
├── frontend/
│   ├── app/               # Next.js app router pages
│   ├── components/        # React components
│   ├── lib/               # Utility functions
│   ├── store/            # Zustand stores
│   └── package.json
├── nginx/                 # Nginx configuration
├── docker-compose.yml
└── README.md
```

## License

MIT License - See [LICENSE](LICENSE) for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## Support

For issues and questions, please open a GitHub issue.
