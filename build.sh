#!/bin/bash
# test-build.sh
MODULE=$(go list -m)  # This gets your module name

echo "Module: $MODULE"
echo "Testing build with version info..."

# Build with test values
go build  -o iwashere-test ./cmd/

# Run version command
./iwashere-test version

# Clean up
# rm iwashere-test