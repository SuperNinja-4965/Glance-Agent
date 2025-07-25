package auth

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
)

// localIPMiddleware restricts access to local IP addresses only
func LocalIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP address
		clientIP := getClientIP(r)

		// Check if IP is local
		if !isLocalIP(clientIP) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "Access denied: Only local connections allowed",
			}); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			log.Printf("Access denied for IP: %s", clientIP)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the real client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		// Take the first IP in the chain
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// isLocalIP checks if an IP address is a local/private address
func isLocalIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check for IPv4 loopback (127.0.0.0/8)
	if ip.IsLoopback() {
		return true
	}

	// Check for IPv6 loopback (::1)
	if ip.Equal(net.IPv6loopback) {
		return true
	}

	// Check for private IPv4 ranges
	privateRanges := []string{
		"10.0.0.0/8",     // Class A private
		"172.16.0.0/12",  // Class B private
		"192.168.0.0/16", // Class C private
		"169.254.0.0/16", // Link-local
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}

	// Check for private IPv6 ranges
	if ip.To4() == nil { // IPv6
		// Check for link-local (fe80::/10)
		if len(ip) >= 2 && ip[0] == 0xfe && (ip[1]&0xc0) == 0x80 {
			return true
		}
		// Check for unique local (fc00::/7)
		if len(ip) >= 1 && (ip[0]&0xfe) == 0xfc {
			return true
		}
	}

	return false
}
