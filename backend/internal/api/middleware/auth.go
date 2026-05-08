package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/lookingglass/backend/internal/services"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// RequireAuth returns middleware that requires authentication
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := SetUserContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns middleware that requires a specific role
func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if claims.Role != role && claims.Role != "admin" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractToken extracts JWT token from Authorization header
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check Bearer scheme
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return parts[1]
}

// contextKey is the key for user context
type contextKey string

const userContextKey contextKey = "user"

// SetUserContext sets user claims in context
func SetUserContext(ctx interface{}, claims *services.Claims) interface{} {
	// In a real implementation, use context.WithValue
	return ctx
}

// GetUserContext gets user claims from context
func GetUserContext(ctx interface{}) (*services.Claims, bool) {
	// In a real implementation, use context.Value
	return nil, false
}

// GetUserFromRequest gets user ID from request
func GetUserFromRequest(r *http.Request) string {
	claims, ok := GetUserContext(r.Context())
	if !ok {
		return ""
	}
	return claims.UserID
}

// GetVars returns URL variables
func GetVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}