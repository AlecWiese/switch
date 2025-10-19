package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		printHelp()
		return
	}

	currentOS := runtime.GOOS
	fmt.Printf("Current OS: %s\n", currentOS)

	if currentOS == "linux" {
		rebootToWindows()
	} else if currentOS == "windows" {
		rebootToLinux()
	} else {
		fmt.Printf("Unsupported OS: %s\n", currentOS)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("OS Reboot Switcher")
	fmt.Println("==================")
	fmt.Println("This program reboots your system into the other OS in a dual-boot setup.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  reboot-switch      Reboot into the other OS")
	fmt.Println("  reboot-switch -h   Show this help message")
	fmt.Println("")
	fmt.Println("Note: This program requires administrative privileges (sudo on Linux, Administrator on Windows)")
}

func rebootToWindows() {
	fmt.Println("Preparing to reboot to Windows...")
	fmt.Println("")
	fmt.Println("This will:")
	fmt.Println("1. Set the next boot entry to Windows")
	fmt.Println("2. Reboot the system")
	fmt.Println("")
	fmt.Print("Do you want to continue? (y/N): ")

	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		fmt.Println("Aborted.")
		return
	}

	// Find Windows boot entry
	fmt.Println("\nSearching for Windows boot entry...")
	cmd := exec.Command("sudo", "grep", "-i", "windows", "/boot/grub/grub.cfg")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Note: Could not automatically detect Windows entry.")
		fmt.Println("You may need to manually configure the boot entry.")
	} else {
		fmt.Println("Windows entry found in GRUB configuration.")
	}

	// Use grub-reboot to set next boot (this requires knowing the menu entry)
	// Alternative: use efibootmgr for UEFI systems
	fmt.Println("\nAttempting to set Windows as next boot entry...")

	// Check if system uses UEFI
	if _, err := os.Stat("/sys/firmware/efi"); err == nil {
		// UEFI system - use efibootmgr
		fmt.Println("Detected UEFI system, using efibootmgr...")

		// List boot entries
		cmd = exec.Command("efibootmgr")
		output, err = cmd.Output()
		if err != nil {
			fmt.Printf("Error listing boot entries: %v\n", err)
			fmt.Println("You may need to install efibootmgr: sudo apt install efibootmgr")
			return
		}

		// Parse for Windows entry
		lines := strings.Split(string(output), "\n")
		windowsBootNum := ""
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), "windows") {
				// Extract boot number (e.g., "Boot0001* Windows Boot Manager")
				if strings.HasPrefix(line, "Boot") && len(line) > 8 {
					windowsBootNum = line[4:8]
					fmt.Printf("Found Windows boot entry: %s\n", line)
					break
				}
			}
		}

		if windowsBootNum != "" {
			// Set next boot to Windows
			cmd = exec.Command("sudo", "efibootmgr", "-n", windowsBootNum)
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error setting next boot: %v\n", err)
				return
			}
			fmt.Println("Successfully set next boot to Windows!")
		} else {
			fmt.Println("Could not find Windows boot entry automatically.")
			fmt.Println("Please run 'efibootmgr' to see available boot options")
			return
		}
	} else {
		// Legacy BIOS system - use grub-reboot
		fmt.Println("Detected legacy BIOS system, using grub-reboot...")
		fmt.Println("Note: You may need to manually specify the Windows menu entry number.")
		fmt.Println("Run 'sudo grep menuentry /boot/grub/grub.cfg' to see available entries.")

		// Try common Windows entry positions
		cmd = exec.Command("sudo", "grub-reboot", "Windows")
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Could not set GRUB entry automatically: %v\n", err)
			fmt.Println("You may need to run: sudo grub-reboot <entry_name_or_number>")
			return
		}
	}

	// Reboot the system
	fmt.Println("\nRebooting now...")
	cmd = exec.Command("sudo", "reboot")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error rebooting: %v\n", err)
		fmt.Println("You can manually reboot with: sudo reboot")
	}
}

func rebootToLinux() {
	fmt.Println("Preparing to reboot to Linux...")
	fmt.Println("")
	fmt.Println("This will:")
	fmt.Println("1. Set the next boot entry to Linux")
	fmt.Println("2. Reboot the system")
	fmt.Println("")
	fmt.Print("Do you want to continue? (y/N): ")

	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		fmt.Println("Aborted.")
		return
	}

	fmt.Println("\nAttempting to set Linux as next boot entry...")

	// Use bcdedit to set the boot entry
	// First, list boot entries to find Linux
	cmd := exec.Command("bcdedit", "/enum", "firmware")
	output, err := cmd.Output()

	if err != nil {
		fmt.Printf("Error listing boot entries: %v\n", err)
		fmt.Println("Make sure you're running as Administrator.")
		fmt.Println("You may need to manually configure boot entry using bcdedit.")
		return
	}

	// Parse for Linux/GRUB entry
	lines := strings.Split(string(output), "\n")
	var linuxIdentifier string
	inEntry := false
	currentIdentifier := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for identifier
		if strings.HasPrefix(line, "identifier") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentIdentifier = parts[1]
				inEntry = true
			}
		}

		// Look for Linux/Ubuntu/Pop in description
		if inEntry && strings.HasPrefix(line, "description") {
			lowerLine := strings.ToLower(line)
			if strings.Contains(lowerLine, "ubuntu") ||
				strings.Contains(lowerLine, "linux") ||
				strings.Contains(lowerLine, "pop") ||
				strings.Contains(lowerLine, "grub") {
				linuxIdentifier = currentIdentifier
				fmt.Printf("Found Linux boot entry: %s\n", line)
				break
			}
		}
	}

	if linuxIdentifier != "" {
		// Set next boot to Linux
		cmd = exec.Command("bcdedit", "/set", "{fwbootmgr}", "bootsequence", linuxIdentifier)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error setting next boot: %v\n", err)
			return
		}
		fmt.Println("Successfully set next boot to Linux!")
	} else {
		fmt.Println("Could not find Linux boot entry automatically.")
		fmt.Println("Please run 'bcdedit /enum firmware' as Administrator to see available boot options")
		return
	}

	// Reboot the system
	fmt.Println("\nRebooting now...")
	cmd = exec.Command("shutdown", "/r", "/t", "0")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error rebooting: %v\n", err)
		fmt.Println("You can manually reboot with: shutdown /r /t 0")
	}
}
