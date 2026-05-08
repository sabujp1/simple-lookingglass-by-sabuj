package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/lookingglass/backend/internal/models"
)

// UserRepository handles user data access
type UserRepository struct {
	db *Database
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, role, api_token, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`
	user.ID = uuid.New()
	return r.db.Pool.QueryRow(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.Role, user.APIToken, user.IsActive,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, api_token, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.APIToken, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, api_token, is_active, created_at, updated_at
		FROM users WHERE username = $1
	`
	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.APIToken, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, api_token, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`
	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.APIToken, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET username = $1, email = $2, role = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.Pool.Exec(ctx, query, user.Username, user.Email, user.Role, user.IsActive, user.ID)
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

// List retrieves all users
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]models.User, int, error) {
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM users`
	var total int
	r.db.Pool.QueryRow(ctx, countQuery).Scan(&total)

	query := `
		SELECT id, username, email, password_hash, role, api_token, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Role, &user.APIToken, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// Count counts all users
func (r *UserRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// Database wraps the database connection
type Database struct {
	Pool interface {
		QueryRow(ctx context.Context, sql string, args ...interface{}) interface {
			Scan(dest ...interface{}) error
		}
		Query(ctx context.Context, sql string, args ...interface{}) (interface {
			Close() error
			Next() bool
			Scan(dest ...interface{}) error
		}, error)
		Exec(ctx context.Context, sql string, args ...interface{}) (interface{}, error)
	}
}