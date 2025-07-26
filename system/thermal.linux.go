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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var thermalZone = -1 // Default thermal zone for CPU temperature monitoring. autodetect.

func SetCPUThermalZone(zone int) {
	// Set the thermal zone for CPU temperature monitoring
	// This is a no-op on non-Linux systems
	if zone < -1 { // -1 means autodetect
		return // Invalid zone, do nothing
	}

	thermalZone = zone
}

type ThermalZone struct {
	Name        string  // e.g. "thermal_zone0"
	Type        string  // e.g. "x86_pkg_temp"
	Temperature float64 // in Celsius
}

// GetThermalZones returns all thermal zones and their temperature readings
func GetThermalZones() ([]ThermalZone, error) {
	basePath := "/sys/class/thermal"
	zones := []ThermalZone{}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read thermal directory: %w", err)
	}

	// Iterate through all entries in the thermal directory
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "thermal_zone") {
			log.Println("Skipping non-thermal zone entry:", entry.Name())
			continue
		}

		zonePath := filepath.Join(basePath, entry.Name())
		typePath := filepath.Join(zonePath, "type")
		tempPath := filepath.Join(zonePath, "temp")

		// Read type
		typeData, err := os.ReadFile(typePath)
		if err != nil {
			log.Println("Skipping zone due to read error on type:", entry.Name(), err)
			continue // skip if we can't read type
		}
		zoneType := strings.TrimSpace(string(typeData))

		// Read temperature
		tempData, err := os.ReadFile(tempPath)
		if err != nil {
			log.Println("Skipping zone due to read error on temp:", entry.Name(), err)
			continue // skip if we can't read temp
		}
		tempMilli, err := strconv.Atoi(strings.TrimSpace(string(tempData)))
		if err != nil {
			log.Println("Skipping zone due to invalid temperature format:", entry.Name(), err)
			continue // invalid number
		}
		tempC := float64(tempMilli) / 1000.0

		zone := ThermalZone{
			Name:        entry.Name(),
			Type:        zoneType,
			Temperature: tempC,
		}
		zones = append(zones, zone)
	}

	return zones, nil
}

func SelectPrimaryCPUThermalZone() (ThermalZone, error) {
	zones, err := GetThermalZones()
	if err != nil {
		return ThermalZone{}, err
	}

	preferredTypes := []string{
		"x86_pkg_temp",
		"cpu_thermal",
		"coretemp",
		"k10temp",
		"proc_thermal",
		"acpitz",
	}
	fmt.Println("Checking zone:", zones)

	// Check for zones with preferred types
	// return the first one found
	for _, preferred := range preferredTypes {
		for _, zone := range zones {
			if zone.Type == preferred {
				return zone, nil
			}
		}
	}

	// No preferred zone found
	return ThermalZone{}, fmt.Errorf("no preferred CPU thermal zone found")
}

// getCPUTemperature reads CPU temperature from thermal zone
// Returns temperature in Celsius, or 0 if unavailable
func getCPUTemperature() int {
	// Autodetect thermal zone if not set
	if thermalZone < 0 {
		zone, err := SelectPrimaryCPUThermalZone()
		if err != nil {
			fmt.Println(zone, err)
			thermalZone = 0 // Use the first detected zone
		} else {
			fmt.Println("Detected primary CPU thermal zone:", zone.Name)
			strippedName := strings.TrimPrefix(zone.Name, "thermal_zone")
			thermalZone, err = strconv.Atoi(strippedName)
			if err != nil {
				log.Println("ERROR: Invalid thermal zone name:", zone.Name)
				return 0 // Invalid zone, return 0
			}
		}
	}

	fmt.Println("Using thermal zone:", thermalZone)

	data, err := os.ReadFile(fmt.Sprintf("/sys/class/thermal/thermal_zone%d/temp", thermalZone))
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
