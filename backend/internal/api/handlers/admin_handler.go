package handlers

import (
	"net/http"

	"github.com/lookingglass/backend/internal/repository"
)

// AdminHandler handles admin endpoints
type AdminHandler struct {
	userRepo *repository.UserRepository
	logRepo  *repository.AuditLogRepository
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(userRepo *repository.UserRepository, logRepo *repository.AuditLogRepository) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
		logRepo:  logRepo,
	}
}

// ListUsers handles listing users
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := getIntQuery(r, "page", 1)
	pageSize := getIntQuery(r, "page_size", 20)

	users, total, err := h.userRepo.List(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":        users,
		"page":        page,
		"page_size":   pageSize,
		"total_count": total,
		"total_pages": totalPages,
	})
}

// GetLogs handles getting audit logs
func (h *AdminHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	page := getIntQuery(r, "page", 1)
	pageSize := getIntQuery(r, "page_size", 50)

	logs, total, err := h.logRepo.List(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":        logs,
		"page":        page,
		"page_size":   pageSize,
		"total_count": total,
		"total_pages": totalPages,
	})
}

// DashboardStats handles dashboard statistics
func (h *AdminHandler) DashboardStats(w http.ResponseWriter, r *http.Request) {
	// This would normally fetch real stats from repositories
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"total_routers":     10,
		"online_routers":     8,
		"offline_routers":    2,
		"queries_per_minute": 15,
		"active_sessions":    5,
		"total_users":        3,
		"total_zones":        4,
	})
}

// ZoneHealth handles zone health status
func (h *AdminHandler) ZoneHealth(w http.ResponseWriter, r *http.Request) {
	// This would normally fetch real health data
	respondJSON(w, http.StatusOK, []map[string]interface{}{
		{"zone_id": "1", "zone_name": "Dhaka", "routers_online": 3, "routers_total": 3, "health_percent": 100},
		{"zone_id": "2", "zone_name": "Chittagong", "routers_online": 2, "routers_total": 3, "health_percent": 66.6},
		{"zone_id": "3", "zone_name": "Singapore", "routers_online": 2, "routers_total": 2, "health_percent": 100},
		{"zone_id": "4", "zone_name": "BDIX", "routers_online": 1, "routers_total": 2, "health_percent": 50},
	})
}