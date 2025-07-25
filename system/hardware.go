package system

import (
	"os"
	"runtime"
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

// GetSystemInfo collects and returns comprehensive system information
func GetSystemInfo() (*SystemInfo, error) {
	// Gather host information
	hostname, platform, bootTime, err := getHostInfo()
	if err != nil {
		return nil, err
	}

	// Get CPU load averages
	load1, load15, err := getLoadAverage()
	if err != nil {
		return nil, err
	}

	// Get number of CPU cores for load percentage calculation
	cpuCount := runtime.NumCPU()

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

	// Calculate load percentages based on CPU count
	// Load average of 1.0 = 100% utilization on single-core system
	load1Percent := int((load1 / float64(cpuCount)) * 100)
	if load1Percent > 100 {
		load1Percent = 100 // Cap at 100%
	}

	load15Percent := int((load15 / float64(cpuCount)) * 100)
	if load15Percent > 100 {
		load15Percent = 100 // Cap at 100%
	}

	// Assemble complete system information
	info := &SystemInfo{
		HostInfoIsAvailable: true,
		BootTime:            bootTime,
		Hostname:            hostname,
		Platform:            platform,
		CPU: CPUInfo{
			LoadIsAvailable:        true,
			Load1Percent:           load1Percent,
			Load15Percent:          load15Percent,
			TemperatureIsAvailable: true,
			TemperatureC:           getCPUTemperature(),
		},
		Memory:      memInfo,
		MountPoints: mountPoints,
	}

	return info, nil
}
