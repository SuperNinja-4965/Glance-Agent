package system

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
