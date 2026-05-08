package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/lookingglass/backend/internal/models"
)

// ZoneRepository handles zone data access
type ZoneRepository struct {
	db *Database
}

// NewZoneRepository creates a new zone repository
func NewZoneRepository(db *Database) *ZoneRepository {
	return &ZoneRepository{db: db}
}

// Create creates a new zone
func (r *ZoneRepository) Create(ctx context.Context, zone *models.Zone) error {
	query := `
		INSERT INTO zones (id, name, code, location, description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	zone.ID = uuid.New()
	return r.db.Pool.QueryRow(ctx, query,
		zone.ID, zone.Name, zone.Code, zone.Location, zone.Description, zone.IsActive,
	).Scan(&zone.CreatedAt, &zone.UpdatedAt)
}

// GetByID retrieves a zone by ID
func (r *ZoneRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Zone, error) {
	query := `
		SELECT id, name, code, location, description, is_active, created_at, updated_at
		FROM zones WHERE id = $1
	`
	var zone models.Zone
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&zone.ID, &zone.Name, &zone.Code, &zone.Location,
		&zone.Description, &zone.IsActive, &zone.CreatedAt, &zone.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("zone not found")
		}
		return nil, err
	}
	return &zone, nil
}

// Update updates a zone
func (r *ZoneRepository) Update(ctx context.Context, zone *models.Zone) error {
	query := `
		UPDATE zones SET name = $1, code = $2, location = $3, description = $4,
			is_active = $5, updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.db.Pool.Exec(ctx, query,
		zone.Name, zone.Code, zone.Location, zone.Description, zone.IsActive, zone.ID,
	)
	return err
}

// Delete deletes a zone
func (r *ZoneRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM zones WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

// List retrieves all zones
func (r *ZoneRepository) List(ctx context.Context, page, pageSize int) ([]models.Zone, int, error) {
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM zones`
	var total int
	r.db.Pool.QueryRow(ctx, countQuery).Scan(&total)

	query := `
		SELECT id, name, code, location, description, is_active, created_at, updated_at
		FROM zones ORDER BY name ASC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var zones []models.Zone
	for rows.Next() {
		var zone models.Zone
		if err := rows.Scan(
			&zone.ID, &zone.Name, &zone.Code, &zone.Location,
			&zone.Description, &zone.IsActive, &zone.CreatedAt, &zone.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		zones = append(zones, zone)
	}

	return zones, total, nil
}

// Count counts all zones
func (r *ZoneRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM zones`
	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}