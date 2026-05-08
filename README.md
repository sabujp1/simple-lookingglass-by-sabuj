# Looking Glass - ISP Network Diagnostics Platform

[![GitHub Stars](https://img.shields.io/github/stars/sabujp1/simple-lookingglass-by-sabuj)](https://github.com/sabujp1/simple-lookingglass-by-sabuj/stargazers)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

> A full-stack web application for performing network diagnostics (ping, traceroute, BGP lookup) across multiple routers from various vendors (MikroTik, Juniper, Cisco, Huawei).

## Links

- GitHub Repository: https://github.com/sabujp1/simple-lookingglass-by-sabuj.git
- Report a Bug: https://github.com/sabujp1/simple-lookingglass-by-sabuj/issues

## Features

- Multi-Vendor Support: SSH-based connection to MikroTik, Juniper, Cisco, and Huawei routers
- Network Diagnostics: Execute ping, traceroute, and BGP lookup queries through a web interface
- Real-time Streaming: Live query results via WebSocket connections
- Role-Based Access Control: Admin, Operator, and User roles with appropriate permissions
- Audit Logging: Complete audit trail of all user actions
- Responsive UI: Modern Next.js frontend with dark/light mode support
- Caching: Redis-based caching for improved performance
- Job Queue: Background processing for long-running queries

---

## Installation

### Prerequisites

Before you begin, ensure you have the following installed:

| Software | Version | Download |
|----------|---------|----------|
| Docker | 20.10+ | https://docs.docker.com/get-docker/ |
| Docker Compose | 2.0+ | https://docs.docker.com/compose/install/ |
| Git | Any | https://git-scm.com/downloads |

Windows users: Install Docker Desktop which includes both Docker Engine and Docker Compose.

### Step 1: Clone the Repository

```bash
git clone https://github.com/sabujp1/simple-lookingglass-by-sabuj.git
cd simple-lookingglass-by-sabuj
```

### Step 2: Configure Environment Variables

Create the environment configuration files:

```bash
# Backend environment
cp backend/.env.example backend/.env

# Frontend environment
cp frontend/.env.local.example frontend/.env.local
```

Edit backend/.env with your settings:

```env
APP_ENV=production
APP_PORT=8080
APP_SECRET=your-super-secret-key

DATABASE_URL=postgres://lookingglass:password@localhost:5432/lookingglass?sslmode=disable

REDIS_URL=redis://localhost:6379

JWT_SECRET=your-jwt-secret
JWT_EXPIRY=24h

CORS_ORIGINS=http://localhost:3000
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

Edit frontend/.env.local:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

### Step 3: Start the Application

Option A: Using Docker Compose (Recommended)

```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f

# Check service status
docker-compose ps
```

Option B: Using the Setup Script

```bash
# Make the script executable (Linux/Mac)
chmod +x setup.sh

# Run the setup script
./setup.sh
```

The setup script will:
- Check prerequisites (Docker, Git)
- Create environment files
- Build and start all containers
- Create a default admin user
- Display access URLs

### Step 4: Verify Installation

After starting, verify all services are running:

```bash
# Check container status
docker-compose ps

# Test frontend
curl http://localhost:3000

# Test backend API
curl http://localhost:8080/health
```

### Step 5: Access the Application

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| API Health | http://localhost:8080/health |

Default Admin Login:
- Username: admin
- Password: admin123

IMPORTANT: Change the default admin password immediately after first login!

---

## Adding Routers

### Via Web UI (Easiest)

1. Navigate to http://localhost:3000
2. Login with admin credentials
3. Go to Admin > Routers
4. Click Add Router
5. Fill in router details:

| Field | Description | Example |
|-------|-------------|---------|
| Name | Display name | Core Router MK-01 |
| Hostname/IP | Router management IP | 192.168.1.1 |
| Vendor | Router vendor | mikrotik, juniper, cisco, huawei |
| Port | SSH port | 22 |
| Username | SSH username | admin |
| Password | SSH password | your-password |
| Zone | Network zone | Core, Edge, etc. |

### Via JSON Configuration

Edit backend/seeds/example-routers.json with your router details.

---

## Architecture

```
Client (Browser)
     Next.js + React Query
           |
           v
   Nginx Reverse Proxy
   (SSL Termination, Routing)
           |
    +------+------+
    v             v
 Frontend     Backend
 (Next.js)    (Go API)
              + WebSocket
                  |
         +--------+--------+
         v                 v
    PostgreSQL           Redis
    (Database)          (Cache)
```

## Tech Stack

### Backend
- Go 1.21 - Server language
- Gin - HTTP web framework
- GORM - ORM for PostgreSQL
- gorilla/websocket - WebSocket support
- golang.org/x/crypto/ssh - SSH client

### Frontend
- Next.js 14 - React framework
- TypeScript - Type safety
- TailwindCSS - Styling
- Tanstack React Query - Data fetching
- Zustand - State management
- shadcn/ui - UI components

### Infrastructure
- PostgreSQL 15 - Primary database
- Redis 7 - Caching and job queue
- Docker - Containerization
- Nginx - Reverse proxy

---

## Useful Commands

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Restart a specific service
docker-compose restart backend

# Shell into PostgreSQL
docker exec -it lookingglass-postgres psql -U lookingglass

# Shell into Redis
docker exec -it lookingglass-redis redis-cli

# Rebuild after code changes
docker-compose up -d --build

# Clean up everything (WARNING: deletes data)
docker-compose down -v --rmi all
```

## Manual Development Setup

### Backend

```bash
cd backend
go mod download
go run cmd/server/main.go migrate
go run cmd/server/main.go serve
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

---

## Project Structure

```
looking-glass/
├── backend/
│   ├── cmd/server/         # Main application entry
│   ├── internal/
│   │   ├── api/            # HTTP handlers and routes
│   │   │   ├── handlers/   # Request handlers
│   │   │   ├── middleware/ # HTTP middleware
│   │   │   └── websocket/  # WebSocket handler
│   │   ├── services/       # Business logic
│   │   ├── repository/     # Data access layer
│   │   ├── vendor/         # Router vendor implementations
│   │   ├── ssh/            # SSH client management
│   │   ├── queue/          # Job queue
│   │   ├── database/       # Database connection
│   │   └── cache/          # Redis cache
│   ├── seeds/              # Database seed data
│   └── Dockerfile
├── frontend/
│   ├── app/                # Next.js app router pages
│   ├── components/         # React components
│   ├── lib/                # Utility functions
│   ├── store/             # Zustand stores
│   └── Dockerfile
├── nginx/                  # Nginx configuration
├── docker-compose.yml
├── setup.sh               # Auto-setup script
└── README.md
```

---

## Security

- All passwords stored with encryption
- JWT-based authentication
- Rate limiting on API endpoints
- SSH credentials stored securely
- Audit logging for all actions

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create a feature branch (git checkout -b feature/AmazingFeature)
3. Commit your changes (git commit -m 'Add AmazingFeature')
4. Push to the branch (git push origin feature/AmazingFeature)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Author

Sabuj - GitHub: @sabujp1
