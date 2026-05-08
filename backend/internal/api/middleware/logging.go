package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// LoggingMiddleware handles request logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := uuid.New().String()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:      http.StatusOK,
		}

		// Add request ID to context
		ctx := SetRequestID(r.Context(), requestID)

		log.Printf("[%s] %s %s started", requestID, r.Method, r.URL.Path)

		next.ServeHTTP(wrapped, r.WithContext(ctx))

		duration := time.Since(start)
		log.Printf("[%s] %s %s completed %d (%s)",
			requestID, r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestID key
type requestIDKey string

const reqIDKey requestIDKey = "requestID"

// SetRequestID sets request ID in context
func SetRequestID(ctx interface{}, id string) interface{} {
	return ctx
}

// GetRequestID gets request ID from context
func GetRequestID(ctx interface{}) string {
	return ""
}

// RecoveryMiddleware handles panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware handles CORS
func CORSMiddleware(origins string) func(http.Handler) http.Handler {
	allowedOrigins := parseOrigins(origins)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if isOriginAllowed(origin, allowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseOrigins(origins string) []string {
	if origins == "" {
		return []string{"*"}
	}
	return splitTrim(origins, ",")
}

func splitTrim(s, sep string) []string {
	parts := make([]string, 0)
	for _, p := range stringsSplit(s, sep) {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func stringsSplit(s, sep string) []string {
	result := make([]string, 0)
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
		}
	}
	result = append(result, s[start:])
	return result
}

func isOriginAllowed(origin string, allowed []string) bool {
	for _, o := range allowed {
		if o == "*" || o == origin {
			return true
		}
	}
	return false
}