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
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func SetCPUThermalZone(_ int) {
	log.Println("Thermal zone setting is not applicable on Windows. Ignoring value")
}

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
