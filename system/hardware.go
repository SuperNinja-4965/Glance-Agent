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
	"runtime"
	"strconv"
	"strings"
)

type FeatureToggleStruct struct {
	// Feature toggles
	DisableCPULoad     bool // disable CPU load monitoring
	DisableTemperature bool // disable temperature monitoring
	DisableMemory      bool // disable memory monitoring
	DisableSwap        bool // disable swap monitoring
	DisableDisk        bool // disable disk monitoring
	DisableHost        bool // disable host information
}

var disabledFeatures FeatureToggleStruct

func SetFeatureToggles(t FeatureToggleStruct) {
	disabledFeatures = t
}

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

// GetSystemInfo collects and returns comprehensive system information
func GetSystemInfo() (*SystemInfo, error) {
	var hostname, platform string
	var bootTime int64
	if !disabledFeatures.DisableHost {
		// Gather host information
		var err error
		hostname, platform, bootTime, err = getHostInfo()
		if err != nil {
			return nil, err
		}
	}

	// Get CPU load averages
	load1, load15, err := getLoadAverage()
	if err != nil {
		return nil, err
	}

	// Gather memory information
	memInfo, err := getMemoryInfo()
	if err != nil {
		return nil, err
	}

	// Get filesystem mount points and usage
	mountPoints, err := getMountPoints()
	if err != nil {
		return nil, err
	}

	load1Percent := 0
	load15Percent := 0
	if !disabledFeatures.DisableCPULoad {
		// Get number of CPU cores for load percentage calculation
		cpuCount := runtime.NumCPU()
		// Calculate load percentages based on CPU count
		// Load average of 1.0 = 100% utilization on single-core system
		load1Percent = int((load1 / float64(cpuCount)) * 100)
		if load1Percent > 100 {
			load1Percent = 100 // Cap at 100%
		}

		load15Percent = int((load15 / float64(cpuCount)) * 100)
		if load15Percent > 100 {
			load15Percent = 100 // Cap at 100%
		}
	}

	CPUTempIsAvailable := false
	CPUTemp := 0
	if !disabledFeatures.DisableTemperature {
		CPUTempIsAvailable = true
		CPUTemp := getCPUTemperature()
		if CPUTemp < 0 {
			CPUTempIsAvailable = false // If temperature is negative, assume not available
		}
	}

	// Assemble complete system information
	info := &SystemInfo{
		HostInfoIsAvailable: !disabledFeatures.DisableHost,
		BootTime:            bootTime,
		Hostname:            hostname,
		Platform:            platform,
		CPU: CPUInfo{
			LoadIsAvailable:        !disabledFeatures.DisableCPULoad,
			Load1Percent:           load1Percent,
			Load15Percent:          load15Percent,
			TemperatureIsAvailable: CPUTempIsAvailable,
			TemperatureC:           CPUTemp,
		},
		Memory:      memInfo,
		MountPoints: mountPoints,
	}

	return info, nil
}
