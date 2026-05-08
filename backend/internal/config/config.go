package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Security SecurityConfig
	SSH      SSHConfig
	Public   PublicConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host              string
	Port              string
	User              string
	Password          string
	Name              string
	MaxConnections    int32
	IdleConnections   int32
	ConnectionLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

type JWTConfig struct {
	Secret         string
	Expiry         time.Duration
	RefreshExpiry  time.Duration
}

type SecurityConfig struct {
	CORSOrigins        string
	RateLimitPerIP     int
	RateLimitWindow    time.Duration
	CommandTimeout     time.Duration
	MaxConcurrentQueries int
}

type SSHConfig struct {
	Timeout         time.Duration
	MaxConnections  int
}

type PublicConfig struct {
	Mode       bool
	ReadOnly   bool
}

type LoggingConfig struct {
	Level  string
	Format string
}

// Load reads configuration from environment variables and config files
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config

	// Server
	cfg.Server = ServerConfig{
		Host:         viper.GetString("SERVER_HOST"),
		Port:         viper.GetString("SERVER_PORT"),
		Mode:         viper.GetString("SERVER_MODE"),
		ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
		WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30 * time.Second
	}

	// Database
	cfg.Database = DatabaseConfig{
		Host:              viper.GetString("DB_HOST"),
		Port:              viper.GetString("DB_PORT"),
		User:              viper.GetString("DB_USER"),
		Password:          viper.GetString("DB_PASSWORD"),
		Name:              viper.GetString("DB_NAME"),
		MaxConnections:    viper.GetInt32("DB_MAX_CONNECTIONS"),
		IdleConnections:   viper.GetInt32("DB_IDLE_CONNECTIONS"),
		ConnectionLifetime: viper.GetDuration("DB_CONNECTION_LIFETIME"),
	}
	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}
	if cfg.Database.Port == "" {
		cfg.Database.Port = "5432"
	}
	if cfg.Database.MaxConnections == 0 {
		cfg.Database.MaxConnections = 25
	}
	if cfg.Database.IdleConnections == 0 {
		cfg.Database.IdleConnections = 5
	}

	// Redis
	cfg.Redis = RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
		PoolSize: viper.GetInt("REDIS_POOL_SIZE"),
	}
	if cfg.Redis.Host == "" {
		cfg.Redis.Host = "localhost"
	}
	if cfg.Redis.Port == "" {
		cfg.Redis.Port = "6379"
	}
	if cfg.Redis.PoolSize == 0 {
		cfg.Redis.PoolSize = 10
	}

	// JWT
	cfg.JWT = JWTConfig{
		Secret:        viper.GetString("JWT_SECRET"),
		Expiry:        viper.GetDuration("JWT_EXPIRY"),
		RefreshExpiry: viper.GetDuration("JWT_REFRESH_EXPIRY"),
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "change-this-secret-in-production"
	}
	if cfg.JWT.Expiry == 0 {
		cfg.JWT.Expiry = 24 * time.Hour
	}
	if cfg.JWT.RefreshExpiry == 0 {
		cfg.JWT.RefreshExpiry = 168 * time.Hour
	}

	// Security
	cfg.Security = SecurityConfig{
		CORSOrigins:         viper.GetString("CORS_ORIGINS"),
		RateLimitPerIP:      viper.GetInt("RATE_LIMIT_PER_IP"),
		RateLimitWindow:     viper.GetDuration("RATE_LIMIT_WINDOW"),
		CommandTimeout:      viper.GetDuration("COMMAND_TIMEOUT"),
		MaxConcurrentQueries: viper.GetInt("MAX_CONCURRENT_QUERIES"),
	}
	if cfg.Security.RateLimitPerIP == 0 {
		cfg.Security.RateLimitPerIP = 100
	}
	if cfg.Security.RateLimitWindow == 0 {
		cfg.Security.RateLimitWindow = time.Minute
	}
	if cfg.Security.CommandTimeout == 0 {
		cfg.Security.CommandTimeout = 30 * time.Second
	}
	if cfg.Security.MaxConcurrentQueries == 0 {
		cfg.Security.MaxConcurrentQueries = 50
	}

	// SSH
	cfg.SSH = SSHConfig{
		Timeout:        viper.GetDuration("SSH_TIMEOUT"),
		MaxConnections: viper.GetInt("SSH_MAX_CONNECTIONS"),
	}
	if cfg.SSH.Timeout == 0 {
		cfg.SSH.Timeout = 30 * time.Second
	}
	if cfg.SSH.MaxConnections == 0 {
		cfg.SSH.MaxConnections = 100
	}

	// Public
	cfg.Public = PublicConfig{
		Mode:     viper.GetBool("PUBLIC_MODE"),
		ReadOnly: true,
	}

	// Logging
	cfg.Logging = LoggingConfig{
		Level:  viper.GetString("LOG_LEVEL"),
		Format: viper.GetString("LOG_FORMAT"),
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "json"
	}

	return &cfg, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

// GetAddr returns the Redis address
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// GetServerAddr returns the server address
func (c *ServerConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}