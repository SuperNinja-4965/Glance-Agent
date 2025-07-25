//go:build windows

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
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ignoredMountpoints defines default volumes to ignore (e.g., floppy drives)
var ignoredMountpoints = []string{
	"A:\\",
	"B:\\",
}

// shouldIgnoreMountpoint checks if a mountpoint should be ignored
func shouldIgnoreMountpoint(mountpoint string) bool {
	for _, ignored := range GetIgnoredMountpoints() {
		if strings.HasPrefix(strings.ToUpper(mountpoint), strings.ToUpper(ignored)) {
			return true
		}
	}
	return false
}

// getMountPoints gathers disk usage info using `wmic logicaldisk`
// and returns a parsed list of MountPoint structs representing each drive
func getMountPoints() ([]MountPoint, error) {
	if disabledFeatures.DisableDisk {
		return []MountPoint{}, nil // Skip if disk monitoring is disabled
	}

	// Use WMIC to query all logical disks with their drive letter, free space, and total size
	cmd := exec.Command("wmic", "logicaldisk", "get", "Caption,Size,FreeSpace")
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the WMIC command and capture output
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("wmic logicaldisk failed: %w", err)
	}

	// Split the output into lines
	lines := strings.Split(out.String(), "\n")
	if len(lines) <= 1 {
		return nil, fmt.Errorf("unexpected WMIC output") // At least one header and one data line expected
	}

	var mountPoints []MountPoint

	// Skip the header and iterate over each line of disk info
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue // Skip empty lines
		}

		// Split the line by whitespace to extract fields: Caption (e.g., C:), FreeSpace, Size
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue // Skip malformed or incomplete lines
		}

		drive := fields[0]    // e.g., "C:"
		freeStr := fields[1]  // e.g., "1234567890"
		totalStr := fields[2] // e.g., "2345678901"

		// Ensure the drive ends with ":" and append "\\" for consistency (e.g., "C:\\" for display)
		if !strings.HasSuffix(drive, ":") {
			drive += ":"
		}
		drive += "\\"

		// Skip ignored mountpoints (e.g., system-reserved volumes)
		if shouldIgnoreMountpoint(drive) {
			continue
		}

		// Parse total and free space from strings to integers
		total, err1 := strconv.ParseInt(totalStr, 10, 64)
		free, err2 := strconv.ParseInt(freeStr, 10, 64)
		if err1 != nil || err2 != nil || total <= 0 {
			continue // Skip if parsing fails or total size is zero/negative
		}

		// Compute used space and convert sizes to megabytes
		used := total - free
		totalMB := int(total / (1024 * 1024))
		usedMB := int(used / (1024 * 1024))
		usedPercent := int((used * 100) / total)

		// Add the drive info to the result list
		mountPoints = append(mountPoints, MountPoint{
			Path:        drive,
			Name:        drive,
			TotalMB:     totalMB,
			UsedMB:      usedMB,
			UsedPercent: usedPercent,
		})
	}

	return mountPoints, nil
}
