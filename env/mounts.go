package env

import (
	"glance-agent/system"
	"log"
	"strings"
)

// configureMountpoints sets up mountpoint ignore lists
func configureMountpoints() {
	// Add extra mountpoints from configuration
	if ignoreMountpoints != "" {
		mountpoints := strings.Split(ignoreMountpoints, ",")
		for i, mp := range mountpoints {
			mountpoints[i] = strings.TrimSpace(mp)
		}
		system.AddIgnoredMountpoints(mountpoints)
		log.Printf("Added ignored mountpoints: %v", mountpoints)
	}

	// Override ignored mountpoints if specified
	if overrideIgnoreMountpoints != "" {
		mountpoints := strings.Split(overrideIgnoreMountpoints, ",")
		for i, mp := range mountpoints {
			mountpoints[i] = strings.TrimSpace(mp)
		}
		system.SetExtraIgnoredMountpoints(mountpoints)
		log.Printf("Override ignored mountpoints: %v", mountpoints)
	}
}
