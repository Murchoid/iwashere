#!/bin/bash
# build.sh - Run this to build for all platforms

echo "🔨 Building iwashere v0.1.0..."

# Clean build directory
rm -rf builds
mkdir builds

$version=v0.2.0
# Function to build for a platform
build() {
    GOOS=$1 GOARCH=$2 go build -o "builds/releases/download/$vesion/iwashere-$1-$2$3" ./cmd
    echo "  ✅ Built for $1/$2"
}

# Linux builds
build linux amd64 ""
build linux 386 ""

# Windows builds  
build windows amd64 ".exe"
build windows 386 ".exe"

# macOS builds
build darwin amd64 ""
build darwin arm64 ""

echo ""
echo "📦 Builds ready in ./builds/"
ls -lh builds/