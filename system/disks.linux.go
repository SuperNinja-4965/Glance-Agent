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
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// ignoredMountpoints defines default filesystem mount points to ignore
var ignoredMountpoints = []string{
	"/snap",
	"/boot/efi",
	"/dev",
	"/proc",
	"/sys",
	"/run",
	"/tmp",
	"/var/tmp",
	"/dev/shm",
	"/run/lock",
	"/sys/fs/cgroup",
	"/boot/grub",
	"/var/lib/docker",
}

// ignoredFilesystems defines filesystem types to ignore
var ignoredFilesystems = []string{
	"proc",
	"sysfs",
	"devtmpfs",
	"tmpfs",
	"cgroup",
	"cgroup2",
	"pstore",
	"bpf",
	"debugfs",
	"tracefs",
	"securityfs",
	"hugetlbfs",
	"mqueue",
	"fusectl",
	"configfs",
}

// shouldIgnoreMountpoint checks if a mountpoint or filesystem type should be ignored
func shouldIgnoreMountpoint(mountpoint, fstype string) bool {
	// Check if mountpoint starts with any ignored path
	for _, ignored := range GetIgnoredMountpoints() {
		if strings.HasPrefix(mountpoint, ignored) {
			return true
		}
	}

	// Check if filesystem type is in ignored list
	for _, ignored := range ignoredFilesystems {
		if fstype == ignored {
			return true
		}
	}

	return false
}

// getMountPoints reads filesystem mount information and calculates disk usage
func getMountPoints() ([]MountPoint, error) {
	if disabledFeatures.DisableDisk {
		return []MountPoint{}, nil // Skip if disk monitoring is disabled
	}

	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "error closing file: %v\n", cerr)
		}
	}()

	var mountPoints []MountPoint
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			continue // Skip malformed lines
		}

		mountpoint := fields[1] // Mount path
		fstype := fields[2]     // Filesystem type

		// Filter out virtual/temporary filesystems and system partitions
		if shouldIgnoreMountpoint(mountpoint, fstype) {
			continue
		}

		// Get filesystem statistics using syscall
		var stat syscall.Statfs_t
		if err := syscall.Statfs(mountpoint, &stat); err != nil {
			continue // Skip if we can't get stats
		}

		// Use fragment size if available, otherwise fall back to block size
		blockSize := uint64(stat.Frsize)
		if blockSize == 0 {
			blockSize = uint64(stat.Bsize)
		}

		// Calculate disk usage in bytes
		total := stat.Blocks * blockSize     // Total space
		available := stat.Bavail * blockSize // Available space for non-root users
		used := total - available            // Actually used space

		// Convert to megabytes
		totalMB := int(total / (1024 * 1024))
		usedMB := int(used / (1024 * 1024))
		usedPercent := 0
		if total > 0 {
			usedPercent = int((used * 100) / total)
		}

		// Only include filesystems with actual storage capacity
		if totalMB > 0 {
			mountPoint := MountPoint{
				Path:        mountpoint,
				Name:        mountpoint, // Use path as display name
				TotalMB:     totalMB,
				UsedMB:      usedMB,
				UsedPercent: usedPercent,
			}
			mountPoints = append(mountPoints, mountPoint)
		}
	}

	return mountPoints, nil
}
