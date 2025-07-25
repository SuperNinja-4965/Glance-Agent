//go:build linux

package system

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
	"os"
	"strconv"
	"strings"
)

// getHostInfo retrieves hostname, platform information, and boot time
func getHostInfo() (string, string, int64, error) {
	// Get system hostname
	hostname, err := os.Hostname()
	if err != nil {
		return "", "", 0, err
	}

	// Get platform/OS information from /etc/os-release
	platform := "Linux" // Default fallback
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				// Extract the pretty name and remove quotes
				platform = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				break
			}
		}
	}

	// Get system boot time from /proc/stat
	bootTime := int64(0)
	if data, err := os.ReadFile("/proc/stat"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "btime ") {
				// Parse boot time as Unix timestamp
				if bt, err := strconv.ParseInt(strings.TrimPrefix(line, "btime "), 10, 64); err == nil {
					bootTime = bt
				}
				break
			}
		}
	}

	return hostname, platform, bootTime, nil
}
