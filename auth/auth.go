package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AuthMiddleware validates Bearer token authentication for protected routes
func Middleware(secretToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			expectedToken := fmt.Sprintf("Bearer %s", secretToken)

			// Check if the Authorization header matches the expected Bearer token
			if authHeader != expectedToken {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(map[string]string{
					"error": "Unauthorized: Invalid token",
				}); err != nil {
					http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
				}
				return
			}

			// Token is valid, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
