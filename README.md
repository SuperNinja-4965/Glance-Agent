# Glance Agent

A lightweight system monitoring agent for Linux and Windows that provides system metrics via a secure HTTP API. This agent is designed to work with Glance dashboard for real-time system monitoring.

## Features

- **System Information**: CPU load, memory usage, disk usage, and host details
- **Security**: Bearer token authentication, local IP restriction and rate limited
- **Configurable**: Customizable ignored mountpoints, flexible configuration, and selective feature monitoring
- **Feature Toggles**: Enable or disable specific monitoring features (CPU, memory, disk, temperature, swap, host info)

## Requirements

- **Linux/Windows operating system**
- **Go 1.24+** for building from source

## Known Issues

- When running within docker OS information is incorrect. On linux systems this can be resolved by adding a volume passing in the host os-release info `/etc/os-release:/etc/os-release:ro`
- On Windows or within docker CPU temperature readings often do not work.

## Installation

### From Source

```bash
git clone <repository-url>
cd glance-agent
go build -o glance-agent
```

## Configuration

The application supports multiple configuration methods with the following precedence order:

1. **Command line flags** (highest priority)
2. **Environment variables**
3. **.env file** (lowest priority)

### Required Configuration

```bash
export SECRET_TOKEN="your-secure-token-here"
```

### Optional Configuration

```bash
# Server port (default: 9012)
export PORT="9012"

# Additional mountpoints to ignore (comma-separated)
export IGNORE_MOUNTPOINTS="/mnt/backup,/media,/opt/custom"

# Additional IPs whitelist (comma-separated)
export WHITELIST_IPS="100.64.0.0/10,fd7a:115c:a1e0::/48"
# Limit access to whats defined on the whitelist alone
export WHITELIST_ONLY="false"

# Override default ignored mountpoints completely
export OVERRIDE_IGNORED_MOUNTPOINTS="/snap,/boot/efi,/custom"

# Feature toggles (default: all features enabled)
export DISABLE_CPU_LOAD="false"
export DISABLE_TEMPERATURE="false"
export DISABLE_MEMORY="false"
export DISABLE_SWAP="false"
export DISABLE_DISK="false"
export DISABLE_HOST="false"
```

### .env File Configuration

The application automatically loads a `.env` file from the same directory as the binary. Create a `.env` file:

```bash
# .env file format
SECRET_TOKEN=your-secret-token-here
PORT=9012
IGNORE_MOUNTPOINTS=/mnt/backup,/media,/opt/custom
OVERRIDE_IGNORED_MOUNTPOINTS=/snap,/boot/efi

# Disable specific features (optional)
DISABLE_TEMPERATURE=true
DISABLE_SWAP=true
```

### Command Line Flags

```bash
./glance-agent -help
```

Available flags:

- `-token`: Bearer token for API authentication (required)
- `-port`: Server port number (default: 9012)
- `-ignore-mounts`: Comma-separated list of additional mountpoints to ignore
- `-override-mounts`: Comma-separated list to override default ignored mountpoints
- `-disable-cpu`: Disable CPU load monitoring
- `-disable-temp`: Disable temperature monitoring
- `-disable-memory`: Disable memory monitoring
- `-disable-swap`: Disable swap monitoring
- `-disable-disk`: Disable disk monitoring
- `-disable-host`: Disable host information
- `-whitelist-only`: Disables the default IP local connection whitelist
- `-help`: Show help message

## Usage

### Starting the Server

```bash
# Using .env file (recommended)
echo "SECRET_TOKEN=my-secure-token" > .env
./glance-agent

# Using environment variables
SECRET_TOKEN="my-secure-token" ./glance-agent

# Using command line flags
./glance-agent -token my-secure-token -port 8080

# Disable specific features
./glance-agent -token my-secure-token -disable-temp -disable-swap

# Combining methods (CLI flags override env vars and .env file)
./glance-agent -token cli-token -ignore-mounts "/custom/mount" -disable-disk
```

### Feature-Specific Examples

```bash
# Monitor only CPU and memory (disable everything else)
./glance-agent -token my-token -disable-temp -disable-swap -disable-disk -disable-host

# Environment variable approach
export SECRET_TOKEN="my-token"
export DISABLE_TEMPERATURE="true"
export DISABLE_SWAP="true"
export DISABLE_DISK="true"
export DISABLE_HOST="true"
./glance-agent

# Minimal monitoring (only host info)
./glance-agent -token my-token -disable-cpu -disable-temp -disable-memory -disable-swap -disable-disk
```

### API Endpoints

#### Get System Information

```bash
curl -H "Authorization: Bearer your-secret-token" \
     http://localhost:9012/api/sysinfo/all
```

**Response Example:**

```json
{
  "host_info_is_available": true,
  "boot_time": 1640995200,
  "hostname": "myserver",
  "platform": "Ubuntu 22.04",
  "cpu": {
    "load_1": 1.2,
    "load_15": 0.8,
    "load_1_percent": 60,
    "load_15_percent": 40,
    "temperature": 45
  },
  "memory": {
    "memory_is_available": true,
    "total_mb": 8192,
    "used_mb": 4096,
    "used_percent": 50,
    "swap_is_available": true,
    "swap_total_mb": 2048,
    "swap_used_mb": 0,
    "swap_used_percent": 0
  },
  "mountpoints": [
    {
      "path": "/",
      "name": "/",
      "total_mb": 51200,
      "used_mb": 25600,
      "used_percent": 50
    }
  ]
}
```

## Feature Toggle Details

### Available Features

| Feature     | CLI Flag             | Environment Variable    | Description                                    |
| ----------- | -------------------- | ----------------------- | ---------------------------------------------- |
| CPU Load    | `--disable-cpu`    | `DISABLE_CPU_LOAD`    | Disables the CPU load averages and percentages |
| Temperature | `--disable-temp`   | `DISABLE_TEMPERATURE` | Disables the CPU temperature monitoring        |
| Memory      | `--disable-memory` | `DISABLE_MEMORY`      | Disables the RAM usage statistics              |
| Swap        | `--disable-swap`   | `DISABLE_SWAP`        | Disables the Swap usage statistics             |
| Disk        | `--disable-disk`   | `DISABLE_DISK`        | Disables the Disk usage for all mountpoints    |
| Host Info   | `--disable-host`   | `DISABLE_HOST`        | Disables the Hostname, platform, boot time     |

## Ignored Mountpoints

### Default Ignored Mountpoints

The following mountpoints are ignored by default:

- `/snap`, `/boot/efi`, `/dev`, `/proc`, `/sys`, `/run`
- `/tmp`, `/var/tmp`, `/dev/shm`, `/run/lock`
- `/sys/fs/cgroup`, `/boot/grub`, `/var/lib/docker`

#### Docker

The sample .env file contains:

- `/usr/lib/os-release`, `/etc/resolv.conf`
- `/etc/hostname`,`/etc/hosts`

#### Windows

On Windows `A:/` and `B:/` are ignored by default

### Default Ignored Filesystems

- `proc`, `sysfs`, `devtmpfs`, `tmpfs`, `cgroup`, `cgroup2`
- `pstore`, `bpf`, `debugfs`, `tracefs`, `securityfs`
- `hugetlbfs`, `mqueue`, `fusectl`, `configfs`

#### Windows

- `none`

### Custom Configuration

#### Add Extra Mountpoints to Ignore

```bash
# Via environment variable
export IGNORE_MOUNTPOINTS="/mnt/backup,/media,/opt/custom"

# Via .env file
echo "IGNORE_MOUNTPOINTS=/mnt/backup,/media,/opt/custom" >> .env

# Via command line flag
./glance-agent -token mytoken -ignore-mounts "/mnt/backup,/media,/opt/custom"
```

#### Override All Ignored Mountpoints

```bash
# Via environment variable
export OVERRIDE_IGNORED_MOUNTPOINTS="/snap,/boot/efi"

# Via .env file
echo "OVERRIDE_IGNORED_MOUNTPOINTS=/snap,/boot/efi" >> .env

# Via command line flag
./glance-agent -token mytoken -override-mounts "/snap,/boot/efi"
```

## Testing

### Basic Functionality Test

```bash
# Start the server with .env file
echo "SECRET_TOKEN=test-token" > .env
./glance-agent

# Or start with command line
./glance-agent -token test-token

# Test in another terminal
curl -H "Authorization: Bearer test-token" \
     http://localhost:9012/api/sysinfo/all | jq '.'
```

### Testing Disabled Features

```bash
# Test with only CPU monitoring enabled
./glance-agent -token test-token -disable-temp -disable-memory -disable-swap -disable-disk -disable-host

# Verify response contains only CPU data
curl -H "Authorization: Bearer test-token" \
     http://localhost:9012/api/sysinfo/all | jq '.cpu'
```

## Deployment

### Systemd Service

Create `/etc/systemd/system/glance-agent.service`:

```ini
[Unit]
Description=Glance System Monitoring Agent
After=network.target

[Service]
Type=simple
User=glance
Group=glance
Environment=SECRET_TOKEN=your-production-token
Environment=PORT=9012
ExecStart=/opt/glance-agent/glance-agent
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Ensure the glance user exists

```bash
sudo useradd --system --no-create-home --shell /bin/false glance
```

### Copy the binary into the correct location and set the permissions

```bash
sudo mkdir -p /opt/glance-agent
sudo cp ./glance-agent.x86_64 /opt/glance-agent/glance-agent
sudo chown root:root /opt/glance-agent/glance-agent
sudo chmod +x /opt/glance-agent/glance-agent
```

### Enable and Start

```bash
sudo systemctl daemon-reload
sudo systemctl enable glance-agent
sudo systemctl start glance-agent
```

> [!WARNING]
> While the application is designed to be as secure as reasonable it is recommended you restrict what clients can access the api using your system's firewall

## License

This software is licenced with the GPL v3 License.

```
Copyright (C) Ava Glass <SuperNinja_4965>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
```

Some exclusions apply to scripts which may be licenses with the MIT license - these files are marked in their header.
