# Glance Agent

A lightweight Linux system monitoring agent that provides system metrics via a secure HTTP API. This agent is designed to work with Glance dashboard for real-time system monitoring.

## Features

- **System Information**: CPU load, memory usage, disk usage, and host details
- **Security**: Bearer token authentication, local IP restriction and rate limited
- **Configurable**: Customizable ignored mountpoints and flexible configuration

## Requirements

- **Linux operating system** (application only builds on Linux)
- **Go 1.24+** for building from source

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

# Override default ignored mountpoints completely
export OVERRIDE_IGNORED_MOUNTPOINTS="/snap,/boot/efi,/custom"
```

### .env File Configuration

The application automatically loads a `.env` file from the same directory as the binary. Create a `.env` file:

```bash
# .env file format
SECRET_TOKEN=your-secret-token-here
PORT=9012
IGNORE_MOUNTPOINTS=/mnt/backup,/media,/opt/custom
OVERRIDE_IGNORED_MOUNTPOINTS=/snap,/boot/efi
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

# Combining methods (CLI flags override env vars and .env file)
./glance-agent -token cli-token -ignore-mounts "/custom/mount"
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
    "total_mb": 8192,
    "used_mb": 4096,
    "used_percent": 50,
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

## Ignored Mountpoints

### Default Ignored Mountpoints

The following mountpoints are ignored by default:

- `/snap`, `/boot/efi`, `/dev`, `/proc`, `/sys`, `/run`
- `/tmp`, `/var/tmp`, `/dev/shm`, `/run/lock`
- `/sys/fs/cgroup`, `/boot/grub`, `/var/lib/docker`

### Default Ignored Filesystems

- `proc`, `sysfs`, `devtmpfs`, `tmpfs`, `cgroup`, `cgroup2`
- `pstore`, `bpf`, `debugfs`, `tracefs`, `securityfs`
- `hugetlbfs`, `mqueue`, `fusectl`, `configfs`

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
