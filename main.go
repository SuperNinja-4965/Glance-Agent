//go:build linux

package main

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
	"glance-agent/auth"
	"glance-agent/env"
	"glance-agent/system"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const Version = "0.1.1"

// init runs before main() and initializes configuration
func init() {
	env.LoadConfig(Version) // Load environment variables from .env file
}

// sysinfoHandler handles requests for system information
func sysinfoHandler(w http.ResponseWriter, _ *http.Request) {
	// Get comprehensive system information
	info, err := system.GetSystemInfo()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		// Generic error message for production
		if encodeErr := json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
		}); encodeErr != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		// Log detailed error server-side only
		log.Printf("System info error: %v", err)
		return
	}

	// Return system information as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// main initializes and starts the HTTP server
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Throttle(10)) // 10 requests per second
	r.Use(auth.LocalIPMiddleware)  // Restrict to local IPs only
	r.Use(auth.SecurityMiddleware) // Add security middleware

	// Protected API routes for system information
	r.Route("/api/sysinfo", func(r chi.Router) {
		r.Use(auth.Middleware(env.GetSecretToken())) // Pass the secret token
		r.Get("/all", sysinfoHandler)
	})

	// Catch-all handler for undefined routes - drops connection
	r.NotFound(auth.DropHandler)

	log.Printf("Server starting on port %s", env.GetPort())
	log.Printf("Configuration: token=%s", maskToken(env.GetSecretToken()))
	log.Fatal(http.ListenAndServe(":"+env.GetPort(), r))
}

// maskToken masks a token for logging purposes
func maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "***" + token[len(token)-4:]
}
