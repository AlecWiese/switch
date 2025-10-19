#!/bin/bash
# Build script for creating binaries for both Linux and Windows

echo "Building OS Reboot Switcher..."
echo ""

# Build for Linux
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o reboot-switch-linux
if [ $? -eq 0 ]; then
    echo "✓ Linux binary created: reboot-switch-linux"
else
    echo "✗ Failed to build Linux binary"
    exit 1
fi

# Build for Windows
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o reboot-switch.exe
if [ $? -eq 0 ]; then
    echo "✓ Windows binary created: reboot-switch.exe"
else
    echo "✗ Failed to build Windows binary"
    exit 1
fi

echo ""
echo "Build complete!"
echo ""
echo "To install on Linux:"
echo "  sudo cp reboot-switch-linux /usr/local/bin/reboot-switch"
echo "  sudo chmod +x /usr/local/bin/reboot-switch"
echo ""
echo "To install on Windows:"
echo "  Copy reboot-switch.exe to a directory in your PATH"
echo "  (e.g., C:\\Windows\\System32 or a custom directory)"
