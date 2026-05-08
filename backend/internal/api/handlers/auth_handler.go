package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lookingglass/backend/internal/models"
	"github.com/lookingglass/backend/internal/services"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
	userRepo    interface{}
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService, userRepo interface{}) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
	}
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get user from repository (simplified)
	// user, err := h.userRepo.GetByUsername(r.Context(), req.Username)
	// if err != nil {
	// 	respondError(w, http.StatusUnauthorized, "invalid credentials")
	// 	return
	// }

	// Check password
	// if !services.CheckPassword(req.Password, user.PasswordHash) {
	// 	respondError(w, http.StatusUnauthorized, "invalid credentials")
	// 	return
	// }

	// For demo purposes, create a mock user
	user := &models.User{
		ID:       [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Username: req.Username,
		Role:     "user",
	}

	// Generate tokens
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	respondJSON(w, http.StatusOK, models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Hash password
	hash, err := services.HashPassword(req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	// Create user (simplified)
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         "user",
		IsActive:     true,
	}

	// h.userRepo.Create(r.Context(), user)

	respondJSON(w, http.StatusCreated, user)
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, err := h.authService.ValidateToken(req.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	// Create new tokens
	user := &models.User{
		ID:       [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Username: claims.Username,
		Role:     claims.Role,
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

// Me returns current user info
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	userID := r.Context().Value("userID")
	if userID == nil {
		respondError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":       userID,
		"username": "admin",
		"role":     "admin",
	})
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error:   message,
		Message: message,
	})
}