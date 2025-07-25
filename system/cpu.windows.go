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
	"math"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// getCPUTemperature attempts to get CPU temperature on Windows
// Returns temperature in Celsius, or 0 if unavailable
func getCPUTemperature() int {
	// Try WMI first (most reliable method)
	if temp := getTemperatureFromWMI(); temp > 0 {
		return temp
	}

	// Temperature monitoring on Windows is limited without admin privileges
	// Most methods require WMI or specialized drivers
	return 0 // Temperature not available
}

// getTemperatureFromWMI uses WMI to get CPU temperature
func getTemperatureFromWMI() int {
	// Use wmic to query temperature from WMI
	cmd := exec.Command("wmic", "/namespace:\\\\root\\wmi", "path", "MSAcpi_ThermalZoneTemperature", "get", "CurrentTemperature", "/value")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CurrentTemperature=") {
			tempStr := strings.TrimPrefix(line, "CurrentTemperature=")
			if temp, err := strconv.Atoi(tempStr); err == nil {
				// WMI returns temperature in tenths of Kelvin
				// Convert to Celsius
				celsius := (temp / 10) - 273
				if celsius > 0 && celsius < 150 { // Sanity check
					return celsius
				}
			}
		}
	}

	return 0
}

// LoadTracker stores exponentially weighted averages of CPU load over time
type LoadTracker struct {
	sync.Mutex
	numCPU        int
	oneMinAvg     float64
	fifteenMinAvg float64
	initialized   bool
}

// Constants for decay factor for exponential moving average (like Unix loadavg)
// Assuming we sample every 5 seconds
const (
	sampleInterval = 5 * time.Second
)

var (
	oneMinDecay     = 1 - math.Exp(-float64(sampleInterval.Seconds())/60)
	fifteenMinDecay = 1 - math.Exp(-float64(sampleInterval.Seconds())/900)
)

// Global load tracker instance
var tracker = &LoadTracker{
	numCPU: runtime.NumCPU(),
}

func init() {
	// Start tracking CPU usage on Windows
	StartTrackingCPUUsage()
}

// StartTrackingCPUUsage begins background collection of CPU usage
// Samples every 5 seconds and updates simulated 1-min and 15-min averages
func StartTrackingCPUUsage() {
	go func() {
		ticker := time.NewTicker(sampleInterval)
		defer ticker.Stop()

		for range ticker.C {
			cpuUsage, err := getCPUUsagePercentage()
			if err != nil {
				continue // skip failed readings
			}
			tracker.update(cpuUsage)
		}
	}()
}

// update applies the new sample to the exponential moving averages
func (lt *LoadTracker) update(cpuUsage float64) {
	lt.Lock()
	defer lt.Unlock()

	// Convert CPU usage percentage to load average equivalent
	load := float64(lt.numCPU) * (cpuUsage / 100.0)

	// Initialize on first sample
	if !lt.initialized {
		lt.oneMinAvg = load
		lt.fifteenMinAvg = load
		lt.initialized = true
	} else {
		// Apply exponential moving average
		lt.oneMinAvg = lt.oneMinAvg*(1-oneMinDecay) + load*oneMinDecay
		lt.fifteenMinAvg = lt.fifteenMinAvg*(1-fifteenMinDecay) + load*fifteenMinDecay
	}
}

// getLoadAverage returns tracked 1-minute and 15-minute simulated load averages
// Returns error if not enough data has been collected yet
func getLoadAverage() (float64, float64, error) {
	tracker.Lock()
	defer tracker.Unlock()

	if !tracker.initialized {
		return 0, 0, fmt.Errorf("load tracker not initialized yet")
	}

	return tracker.oneMinAvg, tracker.fifteenMinAvg, nil
}

// getCPUUsagePercentage calculates CPU usage percentage using wmic
func getCPUUsagePercentage() (float64, error) {
	// Get CPU usage from wmic processor
	cmd := exec.Command("wmic", "cpu", "get", "loadpercentage", "/value")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("wmic cpu command failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "LoadPercentage=") {
			usageStr := strings.TrimPrefix(line, "LoadPercentage=")
			if usage, err := strconv.ParseFloat(usageStr, 64); err == nil {
				// Ensure reasonable bounds
				if usage < 0 {
					usage = 0
				} else if usage > 100 {
					usage = 100
				}
				return usage, nil
			}
		}
	}

	return 0, fmt.Errorf("could not parse CPU usage from wmic output")
}
