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
	"os/exec"
	"strconv"
	"strings"
)

type MEMORYSTATUSEX struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

// getMemoryInfo reads memory and swap information using Windows APIs
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
		return memoryInfo, nil
	}

	// Get physical memory info
	if !disabledFeatures.DisableMemory {
		cmd := exec.Command("wmic", "OS", "get", "TotalVisibleMemorySize,FreePhysicalMemory", "/format:list")
		output, err := cmd.Output()
		if err != nil {
			return memoryInfo, fmt.Errorf("failed to get memory info via WMI: %w", err)
		}

		var totalKB, freeKB int64
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "TotalVisibleMemorySize=") {
				value := strings.TrimPrefix(line, "TotalVisibleMemorySize=")
				totalKB, _ = strconv.ParseInt(value, 10, 64)
			} else if strings.HasPrefix(line, "FreePhysicalMemory=") {
				value := strings.TrimPrefix(line, "FreePhysicalMemory=")
				freeKB, _ = strconv.ParseInt(value, 10, 64)
			}
		}

		if totalKB > 0 {
			totalMB := int(totalKB / 1024)
			freeMB := int(freeKB / 1024)
			usedMB := totalMB - freeMB
			usedPercent := (usedMB * 100) / totalMB

			memoryInfo.MemoryIsAvailable = true
			memoryInfo.TotalMB = totalMB
			memoryInfo.UsedMB = usedMB
			memoryInfo.UsedPercent = usedPercent
		}
	}

	// Get page file (swap) info
	if !disabledFeatures.DisableSwap {
		cmd := exec.Command("wmic", "pagefile", "get", "Size,CurrentUsage", "/format:list")
		output, err := cmd.Output()
		if err == nil {
			var totalSwapMB, usedSwapMB int64
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "Size=") {
					value := strings.TrimPrefix(line, "Size=")
					if size, err := strconv.ParseInt(value, 10, 64); err == nil {
						totalSwapMB += size
					}
				} else if strings.HasPrefix(line, "CurrentUsage=") {
					value := strings.TrimPrefix(line, "CurrentUsage=")
					if usage, err := strconv.ParseInt(value, 10, 64); err == nil {
						usedSwapMB += usage
					}
				}
			}

			if totalSwapMB > 0 {
				swapUsedPercent := int((usedSwapMB * 100) / totalSwapMB)
				memoryInfo.SwapIsAvailable = true
				memoryInfo.SwapTotalMB = int(totalSwapMB)
				memoryInfo.SwapUsedMB = int(usedSwapMB)
				memoryInfo.SwapUsedPercent = swapUsedPercent
			}
		}
	}

	return memoryInfo, nil
}
