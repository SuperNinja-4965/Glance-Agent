package system

// CPUInfo contains CPU-related system metrics
type CPUInfo struct {
	LoadIsAvailable        bool `json:"load_is_available"`        // Whether load average data is available
	Load1Percent           int  `json:"load1_percent"`            // 1-minute load average as percentage of CPU capacity
	Load15Percent          int  `json:"load15_percent"`           // 15-minute load average as percentage of CPU capacity
	TemperatureIsAvailable bool `json:"temperature_is_available"` // Whether CPU temperature data is available
	TemperatureC           int  `json:"temperature_c"`            // CPU temperature in Celsius
}

// MemoryInfo contains memory and swap usage metrics
type MemoryInfo struct {
	MemoryIsAvailable bool `json:"memory_is_available"` // Whether memory data is available
	TotalMB           int  `json:"total_mb"`            // Total system memory in megabytes
	UsedMB            int  `json:"used_mb"`             // Used memory in megabytes
	UsedPercent       int  `json:"used_percent"`        // Memory usage as percentage
	SwapIsAvailable   bool `json:"swap_is_available"`   // Whether swap data is available
	SwapTotalMB       int  `json:"swap_total_mb"`       // Total swap space in megabytes
	SwapUsedMB        int  `json:"swap_used_mb"`        // Used swap space in megabytes
	SwapUsedPercent   int  `json:"swap_used_percent"`   // Swap usage as percentage
}

// MountPoint represents a filesystem mount point with usage statistics
type MountPoint struct {
	Path        string `json:"path"`         // Filesystem mount path
	Name        string `json:"name"`         // Display name (same as path)
	TotalMB     int    `json:"total_mb"`     // Total filesystem size in megabytes
	UsedMB      int    `json:"used_mb"`      // Used space in megabytes
	UsedPercent int    `json:"used_percent"` // Disk usage as percentage
}

// SystemInfo is the main structure containing all system metrics
type SystemInfo struct {
	HostInfoIsAvailable bool         `json:"host_info_is_available"` // Whether host information is available
	BootTime            int64        `json:"boot_time"`              // System boot time as Unix timestamp
	Hostname            string       `json:"hostname"`               // System hostname
	Platform            string       `json:"platform"`               // Operating system platform/distribution
	CPU                 CPUInfo      `json:"cpu"`                    // CPU metrics and information
	Memory              MemoryInfo   `json:"memory"`                 // Memory and swap usage information
	MountPoints         []MountPoint `json:"mountpoints"`            // List of filesystem mount points with usage
}
