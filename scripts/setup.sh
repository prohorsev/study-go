#!/usr/bin/env bash
set -euo pipefail

GO_VERSION=1.24.1

if ! command -v go >/dev/null 2>&1; then
    echo "Installing Go ${GO_VERSION}..."
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    rm go${GO_VERSION}.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
fi

if [ ! -f go.mod ]; then
    echo "go.mod not found in $(pwd). Run this script from the repository root."
    exit 1
fi

# Download module dependencies
if command -v go >/dev/null 2>&1; then
    echo "Downloading Go module dependencies..."
    go mod download
fi

echo "Setup complete."
