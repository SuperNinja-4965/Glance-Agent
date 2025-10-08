package env

import (
	"log"
	"strings"
)

var WhitelistIParr []string

// configureMountpoints sets up mountpoint ignore lists
func configureWhitelistIPs() {
	// Add extra mountpoints from configuration
	if whitelistedIPs != "" {
		WhitelistIParr = strings.Split(whitelistedIPs, ",")
		for i, mp := range WhitelistIParr {
			WhitelistIParr[i] = strings.TrimSpace(mp)
		}
		log.Printf("Added whitelisted IPs: %v", WhitelistIParr)
	}
}
