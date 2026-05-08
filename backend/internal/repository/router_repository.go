package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/lookingglass/backend/internal/models"
)

// RouterRepository handles router data access
type RouterRepository struct {
	db *Database
}

// NewRouterRepository creates a new router repository
func NewRouterRepository(db *Database) *RouterRepository {
	return &RouterRepository{db: db}
}

// Create creates a new router
func (r *RouterRepository) Create(ctx context.Context, router *models.Router) error {
	query := `
		INSERT INTO routers (id, hostname, ip_address, vendor, model, asn, ssh_port, 
			ssh_username, ssh_password_encrypted, ssh_key_path, zone_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at
	`
	router.ID = uuid.New()
	return r.db.Pool.QueryRow(ctx, query,
		router.ID, router.Hostname, router.IPAddress, router.Vendor, router.Model,
		router.ASN, router.SSHPort, router.SSHUsername, router.SSHPasswordEnc,
		router.SSHKeyPath, router.ZoneID, router.IsActive,
	).Scan(&router.CreatedAt, &router.UpdatedAt)
}

// GetByID retrieves a router by ID
func (r *RouterRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Router, error) {
	query := `
		SELECT id, hostname, ip_address, vendor, model, asn, ssh_port, ssh_username, 
			ssh_password_encrypted, ssh_key_path, zone_id, is_active, is_online, 
			last_seen, created_at, updated_at
		FROM routers WHERE id = $1
	`
	var router models.Router
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&router.ID, &router.Hostname, &router.IPAddress, &router.Vendor, &router.Model,
		&router.ASN, &router.SSHPort, &router.SSHUsername, &router.SSHPasswordEnc,
		&router.SSHKeyPath, &router.ZoneID, &router.IsActive, &router.IsOnline,
		&router.LastSeen, &router.CreatedAt, &router.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("router not found")
		}
		return nil, err
	}
	return &router, nil
}

// Update updates a router
func (r *RouterRepository) Update(ctx context.Context, router *models.Router) error {
	query := `
		UPDATE routers SET hostname = $1, ip_address = $2, vendor = $3, model = $4,
			asn = $5, ssh_port = $6, ssh_username = $7, ssh_password_encrypted = $8,
			ssh_key_path = $9, zone_id = $10, is_active = $11, updated_at = NOW()
		WHERE id = $12
	`
	_, err := r.db.Pool.Exec(ctx, query,
		router.Hostname, router.IPAddress, router.Vendor, router.Model,
		router.ASN, router.SSHPort, router.SSHUsername, router.SSHPasswordEnc,
		router.SSHKeyPath, router.ZoneID, router.IsActive, router.ID,
	)
	return err
}

// Delete deletes a router
func (r *RouterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM routers WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

// List retrieves all routers
func (r *RouterRepository) List(ctx context.Context, page, pageSize int) ([]models.Router, int, error) {
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM routers`
	var total int
	r.db.Pool.QueryRow(ctx, countQuery).Scan(&total)

	query := `
		SELECT id, hostname, ip_address, vendor, model, asn, ssh_port, ssh_username,
			ssh_password_encrypted, ssh_key_path, zone_id, is_active, is_online,
			last_seen, created_at, updated_at
		FROM routers ORDER BY hostname ASC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var routers []models.Router
	for rows.Next() {
		var router models.Router
		if err := rows.Scan(
			&router.ID, &router.Hostname, &router.IPAddress, &router.Vendor, &router.Model,
			&router.ASN, &router.SSHPort, &router.SSHUsername, &router.SSHPasswordEnc,
			&router.SSHKeyPath, &router.ZoneID, &router.IsActive, &router.IsOnline,
			&router.LastSeen, &router.CreatedAt, &router.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		routers = append(routers, router)
	}

	return routers, total, nil
}

// ListByZone retrieves routers by zone
func (r *RouterRepository) ListByZone(ctx context.Context, zoneID uuid.UUID) ([]models.Router, error) {
	query := `
		SELECT id, hostname, ip_address, vendor, model, asn, ssh_port, ssh_username,
			ssh_password_encrypted, ssh_key_path, zone_id, is_active, is_online,
			last_seen, created_at, updated_at
		FROM routers WHERE zone_id = $1 ORDER BY hostname ASC
	`
	rows, err := r.db.Pool.Query(ctx, query, zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routers []models.Router
	for rows.Next() {
		var router models.Router
		if err := rows.Scan(
			&router.ID, &router.Hostname, &router.IPAddress, &router.Vendor, &router.Model,
			&router.ASN, &router.SSHPort, &router.SSHUsername, &router.SSHPasswordEnc,
			&router.SSHKeyPath, &router.ZoneID, &router.IsActive, &router.IsOnline,
			&router.LastSeen, &router.CreatedAt, &router.UpdatedAt,
		); err != nil {
			return nil, err
		}
		routers = append(routers, router)
	}

	return routers, nil
}

// UpdateStatus updates router online status
func (r *RouterRepository) UpdateStatus(ctx context.Context, id uuid.UUID, isOnline bool) error {
	query := `UPDATE routers SET is_online = $1, last_seen = NOW() WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, isOnline, id)
	return err
}

// CountOnline counts online routers
func (r *RouterRepository) CountOnline(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM routers WHERE is_online = true AND is_active = true`
	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// CountTotal counts total routers
func (r *RouterRepository) CountTotal(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM routers WHERE is_active = true`
	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}