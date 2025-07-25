package env

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
