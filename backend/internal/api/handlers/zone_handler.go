package handlers

import (
	"net/http"

	"github.com/lookingglass/backend/internal/models"
	"github.com/lookingglass/backend/internal/repository"
)

// ZoneHandler handles zone endpoints
type ZoneHandler struct {
	zoneRepo *repository.ZoneRepository
}

// NewZoneHandler creates a new zone handler
func NewZoneHandler(zoneRepo *repository.ZoneRepository) *ZoneHandler {
	return &ZoneHandler{zoneRepo: zoneRepo}
}

// List handles listing zones
func (h *ZoneHandler) List(w http.ResponseWriter, r *http.Request) {
	page := getIntQuery(r, "page", 1)
	pageSize := getIntQuery(r, "page_size", 20)

	zones, total, err := h.zoneRepo.List(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	respondJSON(w, http.StatusOK, models.PaginatedResponse{
		Data:       zones,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// Get handles getting a zone by ID
func (h *ZoneHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := getPathVar(r, "id")
	zone, err := h.zoneRepo.GetByID(r.Context(), parseUUID(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "zone not found")
		return
	}

	respondJSON(w, http.StatusOK, zone)
}

// Create handles creating a zone
func (h *ZoneHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateZoneRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	zone := &models.Zone{
		Name:        req.Name,
		Code:        req.Code,
		Location:    req.Location,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := h.zoneRepo.Create(r.Context(), zone); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, zone)
}

// Update handles updating a zone
func (h *ZoneHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := getPathVar(r, "id")
	zone, err := h.zoneRepo.GetByID(r.Context(), parseUUID(id))
	if err != nil {
		respondError(w, http.StatusNotFound, "zone not found")
		return
	}

	var req models.UpdateZoneRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name != nil {
		zone.Name = *req.Name
	}
	if req.Code != nil {
		zone.Code = *req.Code
	}
	if req.Location != nil {
		zone.Location = *req.Location
	}
	if req.Description != nil {
		zone.Description = *req.Description
	}
	if req.IsActive != nil {
		zone.IsActive = *req.IsActive
	}

	if err := h.zoneRepo.Update(r.Context(), zone); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, zone)
}

// Delete handles deleting a zone
func (h *ZoneHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := getPathVar(r, "id")
	if err := h.zoneRepo.Delete(r.Context(), parseUUID(id)); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}