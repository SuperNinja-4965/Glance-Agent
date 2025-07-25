package env

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
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// loadEnvFile loads environment variables from .env file in the same directory as the binary
func loadEnvFile() {
	// Get the directory where the binary is located
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: Could not determine executable path: %v", err)
		return
	}

	// Get the directory containing the executable
	execDir := filepath.Dir(execPath)
	envFile := filepath.Join(execDir, ".env")

	// Check if .env file exists in the same directory as the binary
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("Warning: Error loading .env file '%s': %v", envFile, err)
		} else {
			log.Printf("Loaded configuration from %s", envFile)
		}
	} else {
		// Silently continue if .env file doesn't exist
		log.Printf("No .env file found at %s", envFile)
	}
}
