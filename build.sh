#!/bin/bash
# build.sh - Run this to build for all platforms

echo "🔨 Building iwashere v0.1.0..."

# Clean build directory
rm -rf builds
mkdir builds

# Function to build for a platform
build() {
    GOOS=$1 GOARCH=$2 go build -o "builds/iwashere-$1-$2$3" ./cmd/iwashere
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