package handlers

import (
	"net/http"

	"github.com/lookingglass/backend/internal/models"
	"github.com/lookingglass/backend/internal/services"
)

// QueryHandler handles looking glass query endpoints
type QueryHandler struct {
	queryService *services.QueryService
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(queryService *services.QueryService) *QueryHandler {
	return &QueryHandler{queryService: queryService}
}

// Ping handles ping queries
func (h *QueryHandler) Ping(w http.ResponseWriter, r *http.Request) {
	var req models.PingRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	options := map[string]interface{}{
		"count": req.Count,
		"size":  req.Size,
	}

	result, err := h.queryService.Ping(r.Context(), req.RouterID, req.Target, options)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// Traceroute handles traceroute queries
func (h *QueryHandler) Traceroute(w http.ResponseWriter, r *http.Request) {
	var req models.TracerouteRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.queryService.Traceroute(r.Context(), req.RouterID, req.Target)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// BGPRoute handles BGP route lookup
func (h *QueryHandler) BGPRoute(w http.ResponseWriter, r *http.Request) {
	var req models.BGPRouteRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.queryService.BGPRoute(r.Context(), req.RouterID, req.Prefix)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// BGPSummary handles BGP summary queries
func (h *QueryHandler) BGPSummary(w http.ResponseWriter, r *http.Request) {
	var req models.BGPSummaryRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.queryService.BGPSummary(r.Context(), req.RouterID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}