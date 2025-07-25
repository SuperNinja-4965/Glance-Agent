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
	"strconv"
	"strings"
)

// getMemoryInfo reads memory and swap information from /proc/meminfo
func getMemoryInfo() (MemoryInfo, error) {
	var memoryInfo = MemoryInfo{
		MemoryIsAvailable: false,
		TotalMB:           0,
		UsedMB:            0,
		UsedPercent:       0,
		SwapIsAvailable:   false,
		SwapTotalMB:       0,
		SwapUsedMB:        0,
		SwapUsedPercent:   0,
	}
	if disabledFeatures.DisableMemory && disabledFeatures.DisableSwap {
		return memoryInfo, nil // Skip if both memory and swap monitoring are disabled
	}

	// Read the memory and swap information from /proc/meminfo
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemoryInfo{}, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "error closing file: %v\n", cerr)
		}
	}()

	// Parse /proc/meminfo into a map of key-value pairs
	memInfo := make(map[string]int64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":") // Remove trailing colon
		value, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}
		memInfo[key] = value * 1024 // Convert from kB to bytes
	}

	// Only proceed if we are checking memory
	if !disabledFeatures.DisableMemory {
		// Calculate memory usage statistics
		totalMB := int(memInfo["MemTotal"] / (1024 * 1024))
		memoryIsAvailable := memInfo["MemTotal"] > 0

		// Check if MemAvailable exists (Linux 3.14+), fallback to calculation if not
		var availableMB int
		if memAvailable, exists := memInfo["MemAvailable"]; exists {
			availableMB = int(memAvailable / (1024 * 1024))
		} else {
			// Fallback calculation for older kernels
			freeMB := int(memInfo["MemFree"] / (1024 * 1024))
			buffersMB := int(memInfo["Buffers"] / (1024 * 1024))
			cachedMB := int(memInfo["Cached"] / (1024 * 1024))
			availableMB = freeMB + buffersMB + cachedMB
		}

		usedMB := totalMB - availableMB
		usedPercent := 0
		if totalMB > 0 {
			usedPercent = (usedMB * 100) / totalMB
		}

		memoryInfo.MemoryIsAvailable = memoryIsAvailable
		memoryInfo.TotalMB = totalMB
		memoryInfo.UsedMB = usedMB
		memoryInfo.UsedPercent = usedPercent
	}

	// Only proceed if we are checking swap
	if disabledFeatures.DisableSwap {
		// Calculate swap usage statistics
		swapTotalMB := int(memInfo["SwapTotal"] / (1024 * 1024))
		swapFreeMB := int(memInfo["SwapFree"] / (1024 * 1024))
		swapUsedMB := swapTotalMB - swapFreeMB
		swapUsedPercent := 0
		if swapTotalMB > 0 {
			swapUsedPercent = (swapUsedMB * 100) / swapTotalMB
		}

		// Determine if swap is available (some systems have no swap configured)
		swapIsAvailable := swapTotalMB > 0

		memoryInfo.SwapIsAvailable = swapIsAvailable
		memoryInfo.SwapTotalMB = swapTotalMB
		memoryInfo.SwapUsedMB = swapUsedMB
		memoryInfo.SwapUsedPercent = swapUsedPercent
	}

	return memoryInfo, nil
}
