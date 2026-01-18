#!/bin/bash
# Build script for creating binaries for both Linux and Windows

echo "Building OS Reboot Switcher..."
echo ""

# Create releases directory
mkdir -p releases

# Build for Linux
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o releases/switch-linux-amd64
if [ $? -eq 0 ]; then
    echo "✓ Linux binary created: releases/switch-linux-amd64"
else
    echo "✗ Failed to build Linux binary"
    exit 1
fi

# Build for Windows
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o releases/switch-windows-amd64.exe
if [ $? -eq 0 ]; then
    echo "✓ Windows binary created: releases/switch-windows-amd64.exe"
else
    echo "✗ Failed to build Windows binary"
    exit 1
fi

echo ""
echo "Build complete!"
echo ""
echo "To install on Linux:"
echo "  sudo cp releases/switch-linux-amd64 /usr/local/bin/switch"
echo "  sudo chmod +x /usr/local/bin/switch"
echo ""
echo "To install on Windows:"
echo "  Copy releases/switch-windows-amd64.exe to a directory in your PATH"
echo "  Rename it to switch.exe"
echo "  (e.g., C:\\Windows\\System32 or a custom directory)"
