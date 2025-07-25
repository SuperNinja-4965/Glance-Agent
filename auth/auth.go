package auth

// Copyright (C) Ava Glass <SuperNinja_4965>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

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
