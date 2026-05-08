package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	APIToken     string     `json:"api_token,omitempty"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Zone represents a geographic zone or PoP
type Zone struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Location    string    `json:"location,omitempty"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Routers     []Router  `json:"routers,omitempty"`
}

// Router represents a network router
type Router struct {
	ID               uuid.UUID  `json:"id"`
	Hostname         string     `json:"hostname"`
	IPAddress        string     `json:"ip_address"`
	Vendor           string     `json:"vendor"`
	Model            string     `json:"model,omitempty"`
	ASN              *int       `json:"asn,omitempty"`
	SSHPort          int        `json:"ssh_port"`
	SSHUsername      string     `json:"ssh_username"`
	SSHPasswordEnc   string     `json:"-"`
	SSHKeyPath       string     `json:"-"`
	ZoneID           *uuid.UUID `json:"zone_id,omitempty"`
	IsActive         bool       `json:"is_active"`
	IsOnline         bool       `json:"is_online"`
	LastSeen         *time.Time `json:"last_seen,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id,omitempty"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type,omitempty"`
	ResourceID   *string   `json:"resource_id,omitempty"`
	IPAddress    string    `json:"ip_address,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	Details      string    `json:"details,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// QueryHistory represents a query history entry
type QueryHistory struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id,omitempty"`
	RouterID       uuid.UUID `json:"router_id,omitempty"`
	CommandType    string    `json:"command_type"`
	Target         string    `json:"target"`
	Parameters     string    `json:"parameters,omitempty"`
	ResultSummary  string    `json:"result_summary,omitempty"`
	ExecutionTime  int       `json:"execution_time_ms,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// Vendor represents supported router vendors
type Vendor string

const (
	VendorMikroTik Vendor = "mikrotik"
	VendorJuniper  Vendor = "juniper"
	VendorCisco    Vendor = "cisco"
	VendorHuawei   Vendor = "huawei"
)

// Role represents user roles
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleViewer Role = "viewer"
)

// Query types
type QueryType string

const (
	QueryTypePing       QueryType = "ping"
	QueryTypeTraceroute QueryType = "traceroute"
	QueryTypeBGPRoute   QueryType = "bgp_route"
	QueryTypeBGPSummary QueryType = "bgp_summary"
	QueryTypeASN        QueryType = "asn"
	QueryTypePrefix     QueryType = "prefix"
)

// QueryStatus represents the status of a query
type QueryStatus string

const (
	QueryStatusPending   QueryStatus = "pending"
	QueryStatusRunning   QueryStatus = "running"
	QueryStatusCompleted QueryStatus = "completed"
	QueryStatusFailed    QueryStatus = "failed"
)

// Request/Response models

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateRouterRequest struct {
	Hostname   string     `json:"hostname" validate:"required"`
	IPAddress  string     `json:"ip_address" validate:"required,ip"`
	Vendor     string     `json:"vendor" validate:"required"`
	Model      string     `json:"model,omitempty"`
	ASN        *int       `json:"asn,omitempty"`
	SSHPort    int        `json:"ssh_port,omitempty"`
	SSHUsername string    `json:"ssh_username" validate:"required"`
	SSHPassword string    `json:"ssh_password,omitempty"`
	SSHKeyPath string     `json:"ssh_key_path,omitempty"`
	ZoneID     *uuid.UUID `json:"zone_id,omitempty"`
	IsActive   bool       `json:"is_active,omitempty"`
}

type UpdateRouterRequest struct {
	Hostname   *string    `json:"hostname,omitempty"`
	IPAddress  *string    `json:"ip_address,omitempty"`
	Vendor     *string    `json:"vendor,omitempty"`
	Model      *string    `json:"model,omitempty"`
	ASN        *int       `json:"asn,omitempty"`
	SSHPort    *int       `json:"ssh_port,omitempty"`
	SSHUsername *string   `json:"ssh_username,omitempty"`
	SSHPassword *string   `json:"ssh_password,omitempty"`
	SSHKeyPath *string    `json:"ssh_key_path,omitempty"`
	ZoneID     *uuid.UUID `json:"zone_id,omitempty"`
	IsActive   *bool      `json:"is_active,omitempty"`
}

type CreateZoneRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Location    string `json:"location,omitempty"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active,omitempty"`
}

type UpdateZoneRequest struct {
	Name        *string `json:"name,omitempty"`
	Code        *string `json:"code,omitempty"`
	Location    *string `json:"location,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type QueryRequest struct {
	RouterID string `json:"router_id" validate:"required"`
	Target   string `json:"target" validate:"required"`
}

type PingRequest struct {
	RouterID string `json:"router_id" validate:"required"`
	Target   string `json:"target" validate:"required"`
	Count    int    `json:"count,omitempty"`
	Size     int    `json:"size,omitempty"`
}

type TracerouteRequest struct {
	RouterID string `json:"router_id" validate:"required"`
	Target   string `json:"target" validate:"required"`
}

type BGPRouteRequest struct {
	RouterID string `json:"router_id" validate:"required"`
	Prefix   string `json:"prefix,omitempty"`
}

type BGPSummaryRequest struct {
	RouterID string `json:"router_id" validate:"required"`
}

type QueryResponse struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Output    string `json:"output,omitempty"`
	Error     string `json:"error,omitempty"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at,omitempty"`
}

type DashboardStats struct {
	TotalRouters     int `json:"total_routers"`
	OnlineRouters    int `json:"online_routers"`
	OfflineRouters    int `json:"offline_routers"`
	QueriesPerMinute int `json:"queries_per_minute"`
	ActiveSessions   int `json:"active_sessions"`
	TotalUsers       int `json:"total_users"`
	TotalZones       int `json:"total_zones"`
}

type ZoneHealth struct {
	ZoneID        uuid.UUID `json:"zone_id"`
	ZoneName      string    `json:"zone_name"`
	RoutersOnline int       `json:"routers_online"`
	RoutersTotal  int       `json:"routers_total"`
	HealthPercent float64   `json:"health_percent"`
}

type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}