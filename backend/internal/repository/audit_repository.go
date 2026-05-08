package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/lookingglass/backend/internal/models"
)

// AuditLogRepository handles audit log data access
type AuditLogRepository struct {
	db *Database
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *Database) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create creates a new audit log entry
func (r *AuditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, user_id, action, resource_type, resource_id, ip_address, user_agent, details)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at
	`
	log.ID = uuid.New()
	return r.db.Pool.QueryRow(ctx, query,
		log.ID, log.UserID, log.Action, log.ResourceType,
		log.ResourceID, log.IPAddress, log.UserAgent, log.Details,
	).Scan(&log.CreatedAt)
}

// List retrieves audit logs with pagination
func (r *AuditLogRepository) List(ctx context.Context, page, pageSize int) ([]models.AuditLog, int, error) {
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM audit_logs`
	var total int
	r.db.Pool.QueryRow(ctx, countQuery).Scan(&total)

	query := `
		SELECT id, user_id, action, resource_type, resource_id, ip_address, user_agent, details, created_at
		FROM audit_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.ResourceType,
			&log.ResourceID, &log.IPAddress, &log.UserAgent, &log.Details, &log.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

// ListByUser retrieves audit logs for a specific user
func (r *AuditLogRepository) ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.AuditLog, int, error) {
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1`
	var total int
	r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&total)

	query := `
		SELECT id, user_id, action, resource_type, resource_id, ip_address, user_agent, details, created_at
		FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.ResourceType,
			&log.ResourceID, &log.IPAddress, &log.UserAgent, &log.Details, &log.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

// QueryHistoryRepository handles query history data access
type QueryHistoryRepository struct {
	db *Database
}

// NewQueryHistoryRepository creates a new query history repository
func NewQueryHistoryRepository(db *Database) *QueryHistoryRepository {
	return &QueryHistoryRepository{db: db}
}

// Create creates a new query history entry
func (r *QueryHistoryRepository) Create(ctx context.Context, history *models.QueryHistory) error {
	query := `
		INSERT INTO query_history (id, user_id, router_id, command_type, target, parameters, result_summary, execution_time_ms, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at
	`
	history.ID = uuid.New()
	return r.db.Pool.QueryRow(ctx, query,
		history.ID, history.UserID, history.RouterID, history.CommandType,
		history.Target, history.Parameters, history.ResultSummary,
		history.ExecutionTime, history.Status,
	).Scan(&history.CreatedAt)
}

// List retrieves query history with pagination
func (r *QueryHistoryRepository) List(ctx context.Context, page, pageSize int) ([]models.QueryHistory, int, error) {
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM query_history`
	var total int
	r.db.Pool.QueryRow(ctx, countQuery).Scan(&total)

	query := `
		SELECT id, user_id, router_id, command_type, target, parameters, result_summary, execution_time_ms, status, created_at
		FROM query_history ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var histories []models.QueryHistory
	for rows.Next() {
		var h models.QueryHistory
		if err := rows.Scan(
			&h.ID, &h.UserID, &h.RouterID, &h.CommandType, &h.Target,
			&h.Parameters, &h.ResultSummary, &h.ExecutionTime, &h.Status, &h.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		histories = append(histories, h)
	}

	return histories, total, nil
}

// CountRecent counts queries in the last minute
func (r *QueryHistoryRepository) CountRecent(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM query_history WHERE created_at > NOW() - INTERVAL '1 minute'`
	err := r.db.Pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}