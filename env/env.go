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
	"flag"
	"fmt"
	"glance-agent/system"
	"log"
	"os"
	"path/filepath"
)

// Variables to hold configuration
var (
	secretToken               string                     // Bearer token for API authentication
	port                      string                     // Server port number
	ignoreMountpoints         string                     // Comma-separated list of mountpoints to ignore
	overrideIgnoreMountpoints string                     // Comma-separated list to override default ignored mountpoints
	showHelp                  bool                       // Show help message
	appVersion                string                     // Application version, set by build process
	featureToggles            system.FeatureToggleStruct // Feature toggles
)

// GetSecretToken returns the configured secret token
func GetSecretToken() string {
	return secretToken
}

// GetPort returns the configured server port
func GetPort() string {
	return port
}

// showUsage displays help information
func showUsage() {
	fmt.Printf("Glance Agent %s - Linux System Monitoring Agent\n\n", appVersion)
	fmt.Println("USAGE:")
	fmt.Printf("  %s [OPTIONS]\n\n", filepath.Base(os.Args[0]))
	fmt.Println("OPTIONS:")
	flag.PrintDefaults()
	fmt.Println("\nENVIRONMENT VARIABLES:")
	fmt.Println("  SECRET_TOKEN                   Bearer token for API authentication")
	fmt.Println("  PORT                           Server port number (default: 9012)")
	fmt.Println("  IGNORE_MOUNTPOINTS             Comma-separated additional mountpoints to ignore")
	fmt.Println("  OVERRIDE_IGNORED_MOUNTPOINTS   Comma-separated override for default ignored mountpoints")
	fmt.Println("  DISABLE_CPU_LOAD                Disable CPU load monitoring (default: true)")
	fmt.Println("  DISABLE_TEMPERATURE             Disable temperature monitoring (default: true)")
	fmt.Println("  DISABLE_MEMORY                  Disable memory monitoring (default: true)")
	fmt.Println("  DISABLE_SWAP                    Disable swap monitoring (default: true)")
	fmt.Println("  DISABLE_DISK                    Disable disk monitoring (default: true)")
	fmt.Println("  DISABLE_HOST                    Disable host information (default: true)")
	fmt.Println("\nEXAMPLES:")
	fmt.Printf("  %s -token mytoken -port 8080\n", filepath.Base(os.Args[0]))
	fmt.Printf("  SECRET_TOKEN=mytoken %s\n", filepath.Base(os.Args[0]))
	fmt.Printf("  %s -token mytoken -disable-temp -disable-swap\n", filepath.Base(os.Args[0]))
	fmt.Println("\n.ENV FILE:")
	fmt.Println("  The application will automatically load a .env file from the same directory as the binary.")
	fmt.Println("  Format:")
	fmt.Println("  SECRET_TOKEN=your-secret-token")
	fmt.Println("  PORT=9012")
	fmt.Println("  IGNORE_MOUNTPOINTS=/mnt/backup,/media")
	fmt.Println("  OVERRIDE_IGNORED_MOUNTPOINTS=/snap,/boot/efi")
	fmt.Println("  DISABLE_CPU_LOAD=true")
	fmt.Println("  DISABLE_TEMPERATURE=true")
	fmt.Println("")
	fmt.Printf("Glance Agent Copyright (C) Ava Glass <SuperNinja_4965> \nThis program comes as is with ABSOLUTELY NO WARRANTY. \nThis is free software, and you are welcome to redistribute it \nunder certain conditions; For details please visit https://github.com/SuperNinja-4965/Glance-Agent/blob/main/LICENSE.")
}

func LoadConfig(version string) {
	appVersion = version // Set the application version

	// Define command line flags
	flag.StringVar(&secretToken, "token", "", "Bearer token for API authentication (required)")
	flag.StringVar(&port, "port", "9012", "Server port number")
	flag.StringVar(&ignoreMountpoints, "ignore-mounts", "", "Comma-separated list of additional mountpoints to ignore")
	flag.StringVar(&overrideIgnoreMountpoints, "override-mounts", "", "Comma-separated list to override default ignored mountpoints")
	flag.BoolVar(&showHelp, "help", false, "Show the help message")

	flag.BoolVar(&featureToggles.DisableCPULoad, "disable-cpu", false, "Disable CPU load monitoring")
	flag.BoolVar(&featureToggles.DisableTemperature, "disable-temp", false, "Disable temperature monitoring")
	flag.BoolVar(&featureToggles.DisableMemory, "disable-memory", false, "Disable memory monitoring")
	flag.BoolVar(&featureToggles.DisableSwap, "disable-swap", false, "Disable swap monitoring")
	flag.BoolVar(&featureToggles.DisableDisk, "disable-disk", false, "Disable disk monitoring")
	flag.BoolVar(&featureToggles.DisableHost, "disable-host", false, "Disable host information")

	// Custom usage function
	flag.Usage = showUsage

	// Parse command line flags
	flag.Parse()

	// Show help and exit
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Load .env file if it exists
	loadEnvFile()

	// Set configuration from environment variables and command line flags
	configureFromSources()

	// Validate required configuration
	if secretToken == "" {
		log.Fatal("SECRET_TOKEN is required. Set via environment variable, .env file, or -token flag")
	}

	// Configure mountpoints
	configureMountpoints()

	// Set what features are enabled/disabled
	system.SetFeatureToggles(featureToggles)

}

// configureFromSources sets configuration from multiple sources with precedence:
// 1. Command line flags (highest priority)
// 2. Environment variables
// 3. .env file (lowest priority)
func configureFromSources() {
	// Check if flags were actually set by user
	tokenSet := false
	portSet := false
	cpuFlagSet := false
	tempFlagSet := false
	memoryFlagSet := false
	swapFlagSet := false
	diskFlagSet := false
	hostFlagSet := false

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "token":
			tokenSet = true
		case "port":
			portSet = true
		case "disable-cpu":
			cpuFlagSet = true
		case "disable-temp":
			tempFlagSet = true
		case "disable-memory":
			memoryFlagSet = true
		case "disable-swap":
			swapFlagSet = true
		case "disable-disk":
			diskFlagSet = true
		case "disable-host":
			hostFlagSet = true
		}
	})

	// SECRET_TOKEN: CLI flag > env var
	if !tokenSet {
		if envToken := os.Getenv("SECRET_TOKEN"); envToken != "" {
			secretToken = envToken
		}
	}

	// PORT: CLI flag > env var > default
	if !portSet {
		if envPort := os.Getenv("PORT"); envPort != "" {
			port = envPort
		}
	}

	// IGNORE_MOUNTPOINTS: CLI flag > env var
	if ignoreMountpoints == "" {
		ignoreMountpoints = os.Getenv("IGNORE_MOUNTPOINTS")
	}

	// OVERRIDE_IGNORED_MOUNTPOINTS: CLI flag > env var
	if overrideIgnoreMountpoints == "" {
		overrideIgnoreMountpoints = os.Getenv("OVERRIDE_IGNORED_MOUNTPOINTS")
	}

	// Feature toggles: CLI flag > env var > default (true)
	if !cpuFlagSet {
		if envVal := os.Getenv("DISABLE_CPU_LOAD"); envVal != "" {
			featureToggles.DisableCPULoad = envVal == "true"
		}
	}

	if !tempFlagSet {
		if envVal := os.Getenv("DISABLE_TEMPERATURE"); envVal != "" {
			featureToggles.DisableTemperature = envVal == "true"
		}
	}

	if !memoryFlagSet {
		if envVal := os.Getenv("DISABLE_MEMORY"); envVal != "" {
			featureToggles.DisableMemory = envVal == "true"
		}
	}

	if !swapFlagSet {
		if envVal := os.Getenv("DISABLE_SWAP"); envVal != "" {
			featureToggles.DisableSwap = envVal == "true"
		}
	}

	if !diskFlagSet {
		if envVal := os.Getenv("DISABLE_DISK"); envVal != "" {
			featureToggles.DisableDisk = envVal == "true"
		}
	}

	if !hostFlagSet {
		if envVal := os.Getenv("DISABLE_HOST"); envVal != "" {
			featureToggles.DisableHost = envVal == "true"
		}
	}
}
