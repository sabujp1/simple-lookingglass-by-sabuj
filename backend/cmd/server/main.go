package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/lookingglass/backend/internal/api/handlers"
	"github.com/lookingglass/backend/internal/api/middleware"
	"github.com/lookingglass/backend/internal/api/websocket"
	"github.com/lookingglass/backend/internal/cache"
	"github.com/lookingglass/backend/internal/config"
	"github.com/lookingglass/backend/internal/database"
	"github.com/lookingglass/backend/internal/repository"
	"github.com/lookingglass/backend/internal/services"
	"github.com/lookingglass/backend/internal/ssh"
	"github.com/lookingglass/backend/internal/queue"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v, using defaults", err)
		cfg = &config.Config{}
	}

	// Initialize PostgreSQL
	ctx := context.Background()
	db, err := database.NewDatabase(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis")
	}
	redisCache := cache.NewRedisCache(redisClient)

	// Initialize SSH Pool
	sshPool := ssh.NewSSHPool(cfg.SSH.MaxConnections, cfg.SSH.ConnectionTimeout)
	defer sshPool.Close()
	log.Println("SSH pool initialized")

	// Initialize Queue
	queryQueue := queue.NewQueryQueue(cfg.Queue.Workers, cfg.Queue.MaxQueueSize)
	queryQueue.Start()
	defer queryQueue.Stop()
	log.Println("Query queue started")

	// Initialize Repositories
	dbWrapper := &repository.Database{Pool: db}
	userRepo := repository.NewUserRepository(dbWrapper)
	routerRepo := repository.NewRouterRepository(dbWrapper)
	zoneRepo := repository.NewZoneRepository(dbWrapper)
	auditRepo := repository.NewAuditLogRepository(dbWrapper)

	// Initialize Services
	authService := services.NewAuthService(cfg.JWT)
	queryService := services.NewQueryService(routerRepo)

	// Initialize WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize Middleware
	rateLimiter := middleware.NewRateLimiter(cfg.Security)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService, userRepo)
	routerHandler := handlers.NewRouterHandler(routerRepo)
	zoneHandler := handlers.NewZoneHandler(zoneRepo)
	queryHandler := handlers.NewQueryHandler(queryService)
	adminHandler := handlers.NewAdminHandler(userRepo, auditRepo)

	// Create Router
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware(cfg.Server.CORSOrigins))
	r.Use(rateLimiter.Middleware)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// API v1 routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth routes (public)
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/refresh", authHandler.Refresh).Methods("POST")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.NewAuthMiddleware(authService).RequireAuth)

	// User routes
	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET")

	// Router routes
	protected.HandleFunc("/routers", routerHandler.List).Methods("GET")
	protected.HandleFunc("/routers/{id}", routerHandler.Get).Methods("GET")
	protected.HandleFunc("/routers", routerHandler.Create).Methods("POST")
	protected.HandleFunc("/routers/{id}", routerHandler.Update).Methods("PUT")
	protected.HandleFunc("/routers/{id}", routerHandler.Delete).Methods("DELETE")
	protected.HandleFunc("/routers/{id}/test", routerHandler.TestConnection).Methods("POST")
	protected.HandleFunc("/routers/{id}/status", routerHandler.GetStatus).Methods("GET")

	// Zone routes
	protected.HandleFunc("/zones", zoneHandler.List).Methods("GET")
	protected.HandleFunc("/zones/{id}", zoneHandler.Get).Methods("GET")
	protected.HandleFunc("/zones", zoneHandler.Create).Methods("POST")
	protected.HandleFunc("/zones/{id}", zoneHandler.Update).Methods("PUT")
	protected.HandleFunc("/zones/{id}", zoneHandler.Delete).Methods("DELETE")

	// Query routes
	protected.HandleFunc("/query/ping", queryHandler.Ping).Methods("POST")
	protected.HandleFunc("/query/traceroute", queryHandler.Traceroute).Methods("POST")
	protected.HandleFunc("/query/bgp/route", queryHandler.BGPRoute).Methods("POST")
	protected.HandleFunc("/query/bgp/summary", queryHandler.BGPSummary).Methods("POST")

	// Admin routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.NewAuthMiddleware(authService).RequireRole("admin"))

	admin.HandleFunc("/users", adminHandler.ListUsers).Methods("GET")
	admin.HandleFunc("/logs", adminHandler.GetLogs).Methods("GET")
	admin.HandleFunc("/stats", adminHandler.DashboardStats).Methods("GET")
	admin.HandleFunc("/zones/health", adminHandler.ZoneHealth).Methods("GET")

	// WebSocket route
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(wsHub, w, r)
	})

	// Start server
	srv := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Starting server on %s", cfg.Server.Address())
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

// Config returns server configuration
func (s *Server) Config() *config.Config {
	return s.cfg
}

// Server holds server dependencies
type Server struct {
	cfg          *config.Config
	db           *pgxpool.Pool
	redis        *redis.Client
	sshPool      *ssh.SSHPool
	queryQueue   *queue.QueryQueue
	redisCache   *cache.RedisCache
	wsHub        *websocket.Hub
	authService  *services.AuthService
	queryService *services.QueryService
}
