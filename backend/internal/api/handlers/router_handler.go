package handlers

import (
	"net/http"

	"github.com/lookingglass/backend/internal/models"
	"github.com/lookingglass/backend/internal/repository"
)

// RouterHandler handles router endpoints
type RouterHandler struct {
	routerRepo *repository.RouterRepository
}

// NewRouterHandler creates a new router handler
func NewRouterHandler(routerRepo *repository.RouterRepository) *RouterHandler {
	return &RouterHandler{routerRepo: routerRepo}
}

// List handles listing routers
func (h *RouterHandler) List(w http.ResponseWriter, r *http.Request) {
	page := getIntQuery(r, "page", 1)
	pageSize := getIntQuery(r, "page_size", 20)

	routers, total, err := h.routerRepo.List(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	respondJSON(w, http.StatusOK, models.PaginatedResponse{
		Data:       routers,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// Get handles getting a router by ID
func (h *RouterHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := getPathVar(r, "id")
	router, err := h.routerRepo.GetByID(r.Context(), parseUUID(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "router not found")
		return
	}

	respondJSON(w, http.StatusOK, router)
}

// Create handles creating a router
func (h *RouterHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRouterRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	router := &models.Router{
		Hostname:    req.Hostname,
		IPAddress:   req.IPAddress,
		Vendor:      req.Vendor,
		Model:       req.Model,
		ASN:         req.ASN,
		SSHPort:     req.SSHPort,
		SSHUsername: req.SSHUsername,
		SSHPasswordEnc: req.SSHPassword,
		SSHKeyPath:  req.SSHKeyPath,
		ZoneID:      req.ZoneID,
		IsActive:    req.IsActive,
	}

	if router.SSHPort == 0 {
		router.SSHPort = 22
	}

	if err := h.routerRepo.Create(r.Context(), router); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, router)
}

// Update handles updating a router
func (h *RouterHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := getPathVar(r, "id")
	router, err := h.routerRepo.GetByID(r.Context(), parseUUID(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "router not found")
		return
	}

	var req models.UpdateRouterRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	applyUpdate(router, &req)

	if err := h.routerRepo.Update(r.Context(), router); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, router)
}

// Delete handles deleting a router
func (h *RouterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := getPathVar(r, "id")
	if err := h.routerRepo.Delete(r.Context(), parseUUID(id)); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// TestConnection tests router SSH connection
func (h *RouterHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	id := getPathVar(r, "id")
	router, err := h.routerRepo.GetByID(r.Context(), parseUUID(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "router not found")
		return
	}

	// TODO: Implement actual SSH connection test
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"router":  router.Hostname,
		"message": "Connection test would be performed here",
	})
}

// GetStatus gets router status
func (h *RouterHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	id := getPathVar(r, "id")
	router, err := h.routerRepo.GetByID(r.Context(), parseUUID(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "router not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":        router.ID,
		"hostname":  router.Hostname,
		"is_online": router.IsOnline,
		"last_seen": router.LastSeen,
	})
}

// Helper functions
func getIntQuery(r *http.Request, key string, defaultVal int) int {
	if val := r.URL.Query().Get(key); val != "" {
		if n, ok := parseInt(val); ok {
			return n
		}
	}
	return defaultVal
}

func parseInt(s string) (int, bool) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

func getPathVar(r *http.Request, key string) string {
	// In real implementation, use mux.Vars(r)
	return r.URL.Path
}

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func applyUpdate(router *models.Router, req *models.UpdateRouterRequest) {
	if req.Hostname != nil {
		router.Hostname = *req.Hostname
	}
	if req.IPAddress != nil {
		router.IPAddress = *req.IPAddress
	}
	if req.Vendor != nil {
		router.Vendor = *req.Vendor
	}
	if req.Model != nil {
		router.Model = *req.Model
	}
	if req.ASN != nil {
		router.ASN = req.ASN
	}
	if req.SSHPort != nil {
		router.SSHPort = *req.SSHPort
	}
	if req.SSHUsername != nil {
		router.SSHUsername = *req.SSHUsername
	}
	if req.SSHKeyPath != nil {
		router.SSHKeyPath = *req.SSHKeyPath
	}
	if req.ZoneID != nil {
		router.ZoneID = req.ZoneID
	}
	if req.IsActive != nil {
		router.IsActive = *req.IsActive
	}
}

func parseUUID(s string) [16]byte {
	// Simplified - real implementation uses google/uuid
	var uuid [16]byte
	return uuid
}

func jsonNewDecoder(r *http.Request) *json.Decoder {
	return json.NewDecoder(r.Body)
}

func jsonNewEncoder(w http.ResponseWriter) *json.Encoder {
	return json.NewEncoder(w)
}