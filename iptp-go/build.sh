#!/bin/bash
# build.sh - Improved build script for iptp

set -e  # Exit on error

echo "🔨 Building iptp..."
echo ""

# Detect current platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert arch names
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
esac

# Build local binary first (for current platform)
echo "Building for local use ($OS/$ARCH)..."
go build -o iptp
echo "✓ ./iptp (ready to use!)"
echo ""

# Ask if user wants to build for all platforms
read -p "Build for all platforms? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo ""
    echo "Building for all platforms..."
    mkdir -p dist

    # Build for Linux (amd64)
    echo "Building for Linux (amd64)..."
    GOOS=linux GOARCH=amd64 go build -o dist/iptp-linux-amd64
    echo "✓ dist/iptp-linux-amd64"

    # Build for Linux (arm64)
    echo "Building for Linux (arm64)..."
    GOOS=linux GOARCH=arm64 go build -o dist/iptp-linux-arm64
    echo "✓ dist/iptp-linux-arm64"

    # Build for macOS (amd64 - Intel)
    echo "Building for macOS (Intel)..."
    GOOS=darwin GOARCH=amd64 go build -o dist/iptp-darwin-amd64
    echo "✓ dist/iptp-darwin-amd64"

    # Build for macOS (arm64 - Apple Silicon)
    echo "Building for macOS (Apple Silicon)..."
    GOOS=darwin GOARCH=arm64 go build -o dist/iptp-darwin-arm64
    echo "✓ dist/iptp-darwin-arm64"

    # Build for Windows (amd64)
    echo "Building for Windows (amd64)..."
    GOOS=windows GOARCH=amd64 go build -o dist/iptp-windows-amd64.exe
    echo "✓ dist/iptp-windows-amd64.exe"

    echo ""
    echo "✅ All builds complete!"
    echo ""
    ls -lh dist/
else
    echo ""
    echo "✅ Local build complete! Run with: ./iptp"
fi