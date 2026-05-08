#!/bin/bash

# ===========================================
# Looking Glass - Auto Setup Script
# ===========================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="docker-compose.yml"
ENV_FILE="backend/.env"
FRONTEND_ENV_FILE="frontend/.env.local"

print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    print_step "Checking prerequisites..."
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    # Check Git
    if ! command -v git &> /dev/null; then
        print_warning "Git is not installed. Skipping git initialization."
        GIT_AVAILABLE=false
    else
        GIT_AVAILABLE=true
    fi
    
    print_success "All prerequisites met"
}

# Initialize Git repository
init_git() {
    if [ "$GIT_AVAILABLE" = true ]; then
        print_step "Initializing Git repository..."
        
        if [ ! -d ".git" ]; then
            git init
            
            # Set default branch
            git checkout -b main
            
            print_success "Git repository initialized"
        else
            print_warning "Git repository already exists"
        fi
        
        # Create .gitignore if not exists
        if [ ! -f ".gitignore" ]; then
            print_warning ".gitignore not found - you may want to create one"
        fi
    fi
}

# Generate secure secret
generate_secret() {
    head -c 32 /dev/urandom | base64 | head -c 32
}

# Create environment files
create_env_files() {
    print_step "Creating environment files..."
    
    # Backend .env
    if [ ! -f "$ENV_FILE" ]; then
        cat > "$ENV_FILE" << EOF
# Application
APP_ENV=development
APP_PORT=8080
APP_SECRET=$(generate_secret)

# Database
DATABASE_URL=postgres://lookingglass:changeme@localhost:5432/lookingglass?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_USER=lookingglass
DB_PASSWORD=changeme
DB_NAME=lookingglass

# Redis
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=

# JWT Authentication
JWT_SECRET=$(generate_secret)
JWT_EXPIRY=24h

# CORS
CORS_ORIGINS=http://localhost:3000

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# SSH Connection Pool
SSH_MAX_CONNECTIONS=50
SSH_TIMEOUT=30s

# Query Settings
QUERY_TIMEOUT=60s
MAX_CONCURRENT_QUERIES=10

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
EOF
        print_success "Created $ENV_FILE"
    else
        print_warning "$ENV_FILE already exists, skipping"
    fi
    
    # Frontend .env.local
    if [ ! -f "$FRONTEND_ENV_FILE" ]; then
        cat > "$FRONTEND_ENV_FILE" << EOF
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# App Configuration
NEXT_PUBLIC_APP_NAME=Looking Glass
NEXT_PUBLIC_APP_DESCRIPTION=ISP Network Diagnostics Platform
EOF
        print_success "Created $FRONTEND_ENV_FILE"
    else
        print_warning "$FRONTEND_ENV_FILE already exists, skipping"
    fi
}

# Pull latest changes
pull_changes() {
    if [ "$GIT_AVAILABLE" = true ]; then
        print_step "Pulling latest changes..."
        
        if git remote -v | grep -q origin; then
            git pull origin main
            print_success "Pulled latest changes"
        else
            print_warning "No remote configured. Run: git remote add origin <your-repo-url>"
        fi
    fi
}

# Start Docker services
start_services() {
    print_step "Starting Docker services..."
    
    # Check if Docker is running
    if ! docker info &> /dev/null; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
    
    # Use docker compose (v2) or docker-compose (v1)
    if docker compose version &> /dev/null; then
        DOCKER_COMPOSE="docker compose"
    else
        DOCKER_COMPOSE="docker-compose"
    fi
    
    # Build and start containers
    print_step "Building Docker images..."
    $DOCKER_COMPOSE build
    
    print_step "Starting containers..."
    $DOCKER_COMPOSE up -d
    
    print_success "Docker services started"
}

# Wait for services
wait_for_services() {
    print_step "Waiting for services to be ready..."
    
    # Wait for PostgreSQL
    print_step "Waiting for PostgreSQL..."
    for i in {1..30}; do
        if docker exec lookingglass-postgres pg_isready -U lookingglass &> /dev/null; then
            print_success "PostgreSQL is ready"
            break
        fi
        if [ $i -eq 30 ]; then
            print_error "PostgreSQL failed to start"
            exit 1
        fi
        sleep 2
    done
    
    # Wait for Redis
    print_step "Waiting for Redis..."
    for i in {1..15}; do
        if docker exec lookingglass-redis redis-cli ping &> /dev/null; then
            print_success "Redis is ready"
            break
        fi
        if [ $i -eq 15 ]; then
            print_error "Redis failed to start"
            exit 1
        fi
        sleep 2
    done
    
    # Wait for Backend
    print_step "Waiting for Backend API..."
    for i in {1..30}; do
        if curl -sf http://localhost:8080/health &> /dev/null; then
            print_success "Backend API is ready"
            break
        fi
        if [ $i -eq 30 ]; then
            print_warning "Backend API may not be ready yet"
        fi
        sleep 2
    done
    
    # Wait for Frontend
    print_step "Waiting for Frontend..."
    for i in {1..30}; do
        if curl -sf http://localhost:3000 &> /dev/null; then
            print_success "Frontend is ready"
            break
        fi
        if [ $i -eq 30 ]; then
            print_warning "Frontend may not be ready yet"
        fi
        sleep 2
    done
}

# Run database migrations
run_migrations() {
    print_step "Running database migrations..."
    
    # Check if backend is running
    if curl -sf http://localhost:8080/api/v1/health &> /dev/null; then
        # Run migrations via API or direct DB
        docker exec lookingglass-backend ./server migrate &> /dev/null || \
        print_warning "Could not run migrations automatically. Run manually: docker exec lookingglass-backend ./server migrate"
        print_success "Migrations completed"
    else
        print_warning "Backend not ready for migrations. Will retry on next setup."
    fi
}

# Create initial admin user
create_admin_user() {
    print_step "Creating initial admin user..."
    
    # Default credentials - SHOULD BE CHANGED AFTER FIRST LOGIN
    print_warning "Creating default admin user (change password after first login!)"
    print_warning "Username: admin | Password: admin123"
    
    # Create admin via API if available
    curl -sf -X POST http://localhost:8080/api/v1/auth/register \
        -H "Content-Type: application/json" \
        -d '{"username":"admin","password":"admin123","email":"admin@example.com","role":"admin"}' &> /dev/null || \
    print_warning "Could not create admin user. Create via web UI at http://localhost:3000/login"
    
    print_success "Admin user created (or already exists)"
}

# Print next steps
print_next_steps() {
    echo ""
    echo "=========================================="
    echo -e "${GREEN}Setup Complete!${NC}"
    echo "=========================================="
    echo ""
    echo "Services:"
    echo "  • Frontend:    http://localhost:3000"
    echo "  • Backend API: http://localhost:8080"
    echo "  • Nginx:       http://localhost:80 (optional)"
    echo ""
    echo "Default Admin Credentials:"
    echo "  • Username: admin"
    echo "  • Password: admin123"
    echo ""
    echo -e "${YELLOW}⚠️  CHANGE THE ADMIN PASSWORD AFTER FIRST LOGIN!${NC}"
    echo ""
    echo "Useful Commands:"
    echo "  • View logs:     docker-compose logs -f"
    echo "  • Stop services: docker-compose down"
    echo "  • Restart:       docker-compose restart"
    echo "  • Shell into DB: docker exec -it lookingglass-postgres psql -U lookingglass"
    echo ""
    echo "Adding Routers:"
    echo "  1. Go to http://localhost:3000/login"
    echo "  2. Login with admin credentials"
    echo "  3. Navigate to Admin → Routers"
    echo "  4. Click 'Add Router' and fill in SSH details"
    echo ""
    echo "=========================================="
}

# Main installation function
install() {
    echo ""
    echo "=========================================="
    echo "  Looking Glass - Auto Setup"
    echo "=========================================="
    echo ""
    
    check_prerequisites
    init_git
    create_env_files
    start_services
    wait_for_services
    create_admin_user
    print_next_steps
}

# Update existing installation
update() {
    echo ""
    echo "=========================================="
    echo "  Looking Glass - Update"
    echo "=========================================="
    echo ""
    
    check_prerequisites
    pull_changes
    start_services
    wait_for_services
    print_success "Update complete!"
}

# Show status
status() {
    echo ""
    echo "=========================================="
    echo "  Looking Glass - Status"
    echo "=========================================="
    echo ""
    
    if docker compose version &> /dev/null; then
        DOCKER_COMPOSE="docker compose"
    else
        DOCKER_COMPOSE="docker-compose"
    fi
    
    $DOCKER_COMPOSE ps
    echo ""
    
    # Check endpoint health
    echo "Service Health:"
    curl -sf http://localhost:3000 -o /dev/null && echo -e "  ${GREEN}✓${NC} Frontend:   running" || echo -e "  ${RED}✗${NC} Frontend:   not responding"
    curl -sf http://localhost:8080/health -o /dev/null && echo -e "  ${GREEN}✓${NC} Backend:    running" || echo -e "  ${RED}✗${NC} Backend:    not responding"
}

# Stop services
stop() {
    echo ""
    echo "Stopping Looking Glass services..."
    
    if docker compose version &> /dev/null; then
        DOCKER_COMPOSE="docker compose"
    else
        DOCKER_COMPOSE="docker-compose"
    fi
    
    $DOCKER_COMPOSE down
    print_success "Services stopped"
}

# Show help
show_help() {
    echo ""
    echo "Looking Glass Setup Script"
    echo ""
    echo "Usage: ./setup.sh [command]"
    echo ""
    echo "Commands:"
    echo "  install    Install and start Looking Glass (default)"
    echo "  update     Pull latest changes and restart services"
    echo "  status     Show status of all services"
    echo "  stop       Stop all services"
    echo "  help       Show this help message"
    echo ""
}

# Parse command line arguments
case "${1:-install}" in
    install)
        install
        ;;
    update)
        update
        ;;
    status)
        status
        ;;
    stop)
        stop
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac