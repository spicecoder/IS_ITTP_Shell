#!/bin/bash

# IPTP Build Script
# Builds binaries for macOS, Linux, and Windows

set -e

echo "🚀 Building IPTP with DNS Router..."
echo ""

# Create dist directory
mkdir -p dist

# Download dependencies
echo "📦 Getting dependencies..."
go get github.com/miekg/dns

# Build for macOS (Apple Silicon)
echo "🍎 Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o dist/iptp-darwin-arm64 .

# Build for macOS (Intel)
echo "🍎 Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o dist/iptp-darwin-amd64 .

# Build for Linux (amd64)
echo "🐧 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o dist/iptp-linux-amd64 .

# Build for Linux (arm64)
echo "🐧 Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -o dist/iptp-linux-arm64 .

# Build for Windows (amd64)
echo "🪟 Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o dist/iptp-windows-amd64.exe .

echo ""
echo "✓ Build complete!"
echo ""
echo "Binaries created:"
ls -lh dist/
echo ""
echo "To run on your Mac (Apple Silicon):"
echo "  ./dist/iptp-darwin-arm64"
echo ""
echo "To run DNS router (requires sudo):"
echo "  sudo ./dist/iptp-darwin-arm64"
echo "  [IPTP-1] ~$ dns start"
