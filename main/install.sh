#!/bin/bash
# install.sh - One-liner to install iwashere

set -e

echo "Installing iwashere..."

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64) ARCH="arm64" ;;
    i386) ARCH="386" ;;
esac

# Handle macOS naming
if [ "$OS" = "darwin" ]; then
    OS="darwin"
fi

# Get the latest version (you can hardcode v0.1.0 for now)
VERSION="v0.1.0"
REPO="Murchoid/iwashere"  # Change this!

# Construct download URL
FILENAME="iwashere-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    FILENAME="${FILENAME}.exe"
fi

URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

echo "Downloading from: $URL"

# Download with proper headers and follow redirects
if command -v wget >/dev/null 2>&1; then
    wget -q --show-progress -O /tmp/iwashere "$URL"
else
    curl -L --progress-bar -o /tmp/iwashere "$URL"
fi

# Check if we got HTML instead of binary
if file /tmp/iwashere | grep -q "HTML"; then
    echo "Error: Downloaded HTML instead of binary. URL might be wrong."
    cat /tmp/iwashere | head -n 5
    rm /tmp/iwashere
    exit 1
fi

# Make executable and install
chmod +x /tmp/iwashere
sudo mv /tmp/iwashere /usr/local/bin/iwashere

echo "iwashere installed successfully!"
echo "Run 'iwashere --help' to get started"