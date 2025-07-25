package auth

import (
	"log"
	"net/http"
)

// dropHandler immediately closes the connection without response
func DropHandler(w http.ResponseWriter, r *http.Request) {
	// Get the underlying connection and close it immediately
	if hj, ok := w.(http.Hijacker); ok {
		conn, _, err := hj.Hijack()
		if err == nil {
			if cerr := conn.Close(); cerr != nil {
				log.Printf("Error closing connection: %v", cerr)
			}
			return
		}
	}
	// Fallback if hijacking fails - return 404
	http.NotFound(w, r)
}
