```
███████╗██╗    ██╗██╗████████╗ ██████╗██╗  ██╗
██╔════╝██║    ██║██║╚══██╔══╝██╔════╝██║  ██║
███████╗██║ █╗ ██║██║   ██║   ██║     ███████║
╚════██║██║███╗██║██║   ██║   ██║     ██╔══██║
███████║╚███╔███╔╝██║   ██║   ╚██████╗██║  ██║
╚══════╝ ╚══╝╚══╝ ╚═╝   ╚═╝    ╚═════╝╚═╝  ╚═╝
```                                              
A simple Go program to reboot between Linux and Windows in a dual-boot system.

## Features

- Automatically detects current OS
- Sets next boot entry to the other OS (one-time boot by default)
- Optional persistent mode to change default boot order
- Works with both UEFI and legacy BIOS systems
- Supports Linux (using efibootmgr or grub-reboot)
- Supports Windows (using bcdedit)

## Download

Download pre-built binaries from the [latest release](https://git.moosefeelings.com/9093730/Switcher/releases/latest):

- **Linux (amd64):** [switch-linux-amd64](https://git.moosefeelings.com/9093730/Switcher/releases/download/v1.0.0/switch-linux-amd64)
- **Windows (amd64):** [switch-windows-amd64.exe](https://git.moosefeelings.com/9093730/Switcher/releases/download/v1.0.0/switch-windows-amd64.exe)

## Installation

### Quick Install on Linux:
```bash
# Download the binary
wget https://git.moosefeelings.com/9093730/Switcher/releases/download/v1.0.0/switch-linux-amd64

# Install to /usr/local/bin
sudo cp switch-linux-amd64 /usr/local/bin/switch
sudo chmod +x /usr/local/bin/switch

# Done! Now you can use it
sudo switch
```

### Quick Install on Windows:
1. Download [switch-windows-amd64.exe](https://git.moosefeelings.com/9093730/Switcher/releases/download/v1.0.0/switch-windows-amd64.exe)
2. Rename it to `switch.exe`
3. Copy it to a directory in your PATH (e.g., `C:\Windows\System32`)
4. Run as Administrator: `switch`

## Building from Source

If you prefer to build from source, use the included build script:

```bash
./build.sh
```

This will create:
- `releases/switch-linux-amd64` - Linux binary
- `releases/switch-windows-amd64.exe` - Windows binary

### Manual Build:

**For Linux:**
```bash
GOOS=linux GOARCH=amd64 go build -o switch
```

**For Windows:**
```bash
GOOS=windows GOARCH=amd64 go build -o switch.exe
```

## Usage

### Basic Usage (One-Time Boot)

By default, the tool sets the boot entry for the **next reboot only**. After that single reboot, your system will return to its previous default boot entry.

```bash
# On Linux (to reboot to Windows once)
sudo switch

# On Windows (to reboot to Linux once) - Run as Administrator
switch
```

### Persistent Boot (Change Default)

Use the `-p` flag to **permanently change** your default boot entry:

```bash
# On Linux - make Windows the default boot entry
sudo switch -p

# On Windows - make Linux the default boot entry
switch -p
```

⚠️ **Warning**: Persistent mode changes your BIOS/UEFI boot order permanently. The selected OS will remain the default until you change it again.

### Skip Confirmation

Use the `-y` flag to skip the confirmation prompt (useful for scripts):

```bash
# Automatically confirm and reboot
sudo switch -y

# Combine with persistent mode
sudo switch -y -p
```

### Verbose Mode

Use the `-v` flag to see detailed boot configuration information:

```bash
# Show current and new boot order
sudo switch -v

# Combine with other flags
sudo switch -y -p -v
```

### Available Flags

```bash
switch           # Reboot to other OS (one-time, with confirmation)
switch -y        # Skip confirmation prompt
switch -p        # Make boot selection persistent (change default)
switch -v        # Show verbose output with boot order details
switch -y -p     # Skip confirmation and make persistent
switch -y -p -v  # All flags combined
switch -h        # Show help message
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
