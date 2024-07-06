#!/bin/sh

# Run go test
echo "Running unit tests..."
go test ./...
if [ $? -ne 0 ]; then
     echo "❌ Unit tests failed. Aborting commit."
    exit 1
fi

# Run golangci-lint
echo "Running golangci-lint..."
golangci-lint run
if [ $? -ne 0 ]; then
    echo "❌ golangci-lint found issues. Aborting commit."
    exit 1
fi

echo "🎉 All checks passed. Yay!"
exit 0