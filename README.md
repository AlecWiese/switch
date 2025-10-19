# OS Reboot Switcher

A simple Go program to reboot between Linux and Windows in a dual-boot system.

## Features

- Automatically detects current OS
- Sets next boot entry to the other OS
- Works with both UEFI and legacy BIOS systems
- Supports Linux (using efibootmgr or grub-reboot)
- Supports Windows (using bcdedit)

## Building

### Build for Linux:
```bash
GOOS=linux GOARCH=amd64 go build -o reboot-switch
```

### Build for Windows:
```bash
GOOS=windows GOARCH=amd64 go build -o reboot-switch.exe
```

### Build for both:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o reboot-switch-linux
# Windows
GOOS=windows GOARCH=amd64 go build -o reboot-switch.exe
```

## Installation

### On Linux:
1. Build the binary for Linux
2. Copy it to a directory in your PATH:
   ```bash
   sudo cp reboot-switch-linux /usr/local/bin/reboot-switch
   sudo chmod +x /usr/local/bin/reboot-switch
   ```

### On Windows:
1. Build the binary for Windows
2. Copy `reboot-switch.exe` to a directory in your PATH (e.g., `C:\Windows\System32` or create a custom directory)

## Usage

```bash
# On Linux (to reboot to Windows)
sudo reboot-switch

# On Windows (to reboot to Linux) - Run as Administrator
reboot-switch
```

Show help:
```bash
reboot-switch -h
```

## Requirements

### Linux:
- For UEFI systems: `efibootmgr` package
  ```bash
  sudo apt install efibootmgr
  ```
- For legacy BIOS: GRUB bootloader
- Root/sudo privileges

### Windows:
- Administrator privileges
- bcdedit (included with Windows)

## Notes

- Always run with administrative privileges (sudo on Linux, Administrator on Windows)
- The program will ask for confirmation before rebooting
- If automatic detection fails, you may need to manually configure boot entries
- Make sure both operating systems are properly configured in your boot manager

## Troubleshooting

### Linux:
- If Windows entry is not found, check: `sudo grep -i windows /boot/grub/grub.cfg`
- For UEFI systems, check boot entries: `efibootmgr`
- For BIOS systems, check GRUB menu entries: `sudo grep menuentry /boot/grub/grub.cfg`

### Windows:
- If Linux entry is not found, check: `bcdedit /enum firmware` (as Administrator)
- Make sure you're running Command Prompt or PowerShell as Administrator
- You may need to add Linux to Windows Boot Manager if it's not visible
