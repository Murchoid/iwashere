#!/bin/bash
# test-build.sh
MODULE=$(go list -m)  # This gets your module name

echo "Module: $MODULE"
echo "Testing build with version info..."

# Build with test values
go build -ldflags="
    -X $MODULE/internal/commands.Version=v0.2.0
    -X $MODULE/internal/commands.Commit=$(git rev-parse --short HEAD)
    -X $MODULE/internal/commands.BuildDate=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    -X $MODULE/internal/commands.BuiltBy=hand
" -o iwashere-test ./cmd/

# Run version command
./iwashere-test version

# Clean up
rm iwashere-test