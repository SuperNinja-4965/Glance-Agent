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
	"os"
	"strconv"
	"strings"
)

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
