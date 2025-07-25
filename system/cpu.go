package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// getCPUTemperature reads CPU temperature from thermal zone
// Returns temperature in Celsius, or 0 if unavailable
func getCPUTemperature() int {
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return 0 // Temperature not available
	}

	tempStr := strings.TrimSpace(string(data))
	temp, err := strconv.Atoi(tempStr)
	if err != nil {
		return 0
	}

	// Convert from millidegrees to degrees Celsius
	return temp / 1000
}

// getLoadAverage reads system load averages from /proc/loadavg
// Returns 1-minute and 15-minute load averages
func getLoadAverage() (float64, float64, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, 0, err
	}

	// /proc/loadavg format: "0.52 0.58 0.59 1/467 12345"
	// Fields: 1min 5min 15min running/total_processes last_pid
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return 0, 0, fmt.Errorf("invalid loadavg format")
	}

	// Parse 1-minute load average
	load1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, 0, err
	}

	// Parse 15-minute load average (skip 5-minute for this application)
	load15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return 0, 0, err
	}

	return load1, load15, nil
}
