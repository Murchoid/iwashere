#!/bin/bash
# test-build.sh
MODULE=$(go list -m)  

echo "Building iwashere ..."
# Build with test values
go build  -o iwashere-test ./cmd/

# Run version command
./iwashere-test version

# Clean up
# rm iwashere-test