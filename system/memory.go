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

	// Calculate memory usage statistics
	totalMB := int(memInfo["MemTotal"] / (1024 * 1024))
	availableMB := int(memInfo["MemAvailable"] / (1024 * 1024))
	usedMB := totalMB - availableMB
	usedPercent := 0
	if totalMB > 0 {
		usedPercent = (usedMB * 100) / totalMB
	}

	// Calculate swap usage statistics
	swapTotalMB := int(memInfo["SwapTotal"] / (1024 * 1024))
	swapFreeMB := int(memInfo["SwapFree"] / (1024 * 1024))
	swapUsedMB := swapTotalMB - swapFreeMB
	swapUsedPercent := 0
	if swapTotalMB > 0 {
		swapUsedPercent = (swapUsedMB * 100) / swapTotalMB
	}

	return MemoryInfo{
		MemoryIsAvailable: true,
		TotalMB:           totalMB,
		UsedMB:            usedMB,
		UsedPercent:       usedPercent,
		SwapIsAvailable:   true,
		SwapTotalMB:       swapTotalMB,
		SwapUsedMB:        swapUsedMB,
		SwapUsedPercent:   swapUsedPercent,
	}, nil
}
