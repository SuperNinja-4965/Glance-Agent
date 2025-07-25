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
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// getHostInfo retrieves hostname, platform information, and boot time for Windows
func getHostInfo() (string, string, int64, error) {
	// Get system hostname
	hostname, err := os.Hostname()
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get hostname: %w", err)
	}

	// Get Windows version information
	platform, err := getWindowsVersion()
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get platform information: %w", err)
	}

	// Get system boot time
	bootTime, err := getBootTime()
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get boot time: %w", err)
	}

	return hostname, platform, bootTime, nil
}

// getWindowsVersion retrieves Windows version information
func getWindowsVersion() (string, error) {
	// Try to get version from wmic first
	if version, err := getVersionFromWMIC(); err == nil && version != "" {
		return version, nil
	}

	// Fallback to systeminfo command
	if version, err := getVersionFromSystemInfo(); err == nil && version != "" {
		return version, nil
	}

	// If both methods fail, return error
	return "", fmt.Errorf("unable to determine Windows version")
}

// getVersionFromWMIC uses wmic to get OS information
func getVersionFromWMIC() (string, error) {
	cmd := exec.Command("wmic", "os", "get", "Caption,Version", "/format:list")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("wmic command failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var caption, version string

	// Parse the output to find Caption and Version
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Caption=") {
			caption = strings.TrimPrefix(line, "Caption=")
		} else if strings.HasPrefix(line, "Version=") {
			version = strings.TrimPrefix(line, "Version=")
		}
	}

	if caption != "" && version != "" {
		return caption + " (Version " + version + ")", nil
	} else if caption != "" {
		return caption, nil
	}

	return "", fmt.Errorf("wmic output does not contain expected OS information")
}

// getVersionFromSystemInfo uses systeminfo command as fallback
func getVersionFromSystemInfo() (string, error) {
	cmd := exec.Command("systeminfo")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("systeminfo command failed: %w", err)
	}

	// Parse the output to find the OS Name
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "OS Name:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				osName := strings.TrimSpace(parts[1])
				if osName != "" {
					return osName, nil
				}
			}
		}
	}

	return "", fmt.Errorf("systeminfo output does not contain OS Name")
}

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetTickCount64 = kernel32.NewProc("GetTickCount64")
)

// getBootTime calculates system boot time using GetTickCount64
func getBootTime() (int64, error) {
	// Check if GetTickCount64 is available
	if err := procGetTickCount64.Find(); err != nil {
		return 0, fmt.Errorf("GetTickCount64 not available: %w", err)
	}

	// Get system uptime in milliseconds using GetTickCount64
	ret, _, callErr := procGetTickCount64.Call()
	if callErr != nil && callErr.Error() != "The operation completed successfully." {
		return 0, fmt.Errorf("GetTickCount64 call failed: %w", callErr)
	}

	uptimeMs := int64(ret)
	if uptimeMs == 0 {
		return 0, fmt.Errorf("GetTickCount64 returned invalid uptime")
	}

	// Calculate boot time by subtracting uptime from current time
	now := time.Now().Unix()
	bootTime := now - (uptimeMs / 1000)

	// Sanity check: boot time should be in the past
	if bootTime > now {
		return 0, fmt.Errorf("calculated boot time is in the future")
	}

	return bootTime, nil
}
