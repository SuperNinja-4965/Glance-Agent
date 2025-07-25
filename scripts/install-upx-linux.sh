#!/bin/bash

set -e

# Set the installation directory
INSTALL_DIR="$HOME/.local/bin"

# Determine the system architecture
ARCH=$(uname -m)
OS=$(uname | tr '[:upper:]' '[:lower:]')

# Map the architecture to the expected UPX format
case $ARCH in
  x86_64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  armv7l) ARCH="arm" ;;
  *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "🔍 Detecting system: $OS-$ARCH"

# Fetch the latest version number from GitHub API
UPX_VERSION=$(curl -s https://api.github.com/repos/upx/upx/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')
if [[ -z "$UPX_VERSION" ]]; then
  echo "❌ Failed to fetch the latest UPX version."
  exit 1
fi

# Compose the download URL
ARCHIVE_NAME="upx-${UPX_VERSION:1}-${ARCH}_${OS}.tar.xz"
DOWNLOAD_URL="https://github.com/upx/upx/releases/download/${UPX_VERSION}/${ARCHIVE_NAME}"

echo "⬇️  Downloading UPX version $UPX_VERSION from: $DOWNLOAD_URL"

# Create the install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Download the file
curl -L -o "$ARCHIVE_NAME" "$DOWNLOAD_URL"

# Check the file format to ensure it's a valid tar.xz
if ! file "$ARCHIVE_NAME" | grep -q "XZ compressed data"; then
  echo "❌ Downloaded file is not a valid XZ archive. Removing..."
  rm -f "$ARCHIVE_NAME"
  exit 1
fi

# Extract the binary
tar -xJf "$ARCHIVE_NAME" --strip-components=1 -C "$INSTALL_DIR" upx-${UPX_VERSION:1}-${ARCH}_${OS}/upx

# Clean up the archive
rm -f "$ARCHIVE_NAME"

# Add the install directory to PATH if not already present
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "export PATH=\"$INSTALL_DIR:\$PATH\"" >> ~/.bashrc
  source ~/.bashrc
  echo "🔗 Added $INSTALL_DIR to PATH"
fi

echo "✅ UPX installed successfully!"
echo "🛠️  Version: $(upx --version)"
