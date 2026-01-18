package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	// Define command-line flags
	skipConfirm := flag.Bool("y", false, "Skip confirmation prompt")
	persist := flag.Bool("p", false, "Make boot selection persistent (default is one-time boot)")
	verbose := flag.Bool("v", false, "Show verbose output including boot order details")

	// Set custom usage message
	flag.Usage = printHelp
	flag.Parse()

	currentOS := runtime.GOOS
	fmt.Printf("Current OS: %s\n", currentOS)

	if currentOS == "linux" {
		rebootToWindows(*skipConfirm, *persist, *verbose)
	} else if currentOS == "windows" {
		rebootToLinux(*skipConfirm, *persist, *verbose)
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
	fmt.Println("  switch           Reboot into the other OS (one-time boot)")
	fmt.Println("  switch -y        Skip confirmation prompt")
	fmt.Println("  switch -p        Make boot selection persistent (change default boot order)")
	fmt.Println("  switch -v        Show verbose output with boot order details")
	fmt.Println("  switch -y -p     Skip confirmation and make persistent")
	fmt.Println("  switch -h        Show this help message")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -y    Skip confirmation prompt")
	fmt.Println("  -p    Persist boot selection (changes default boot order permanently)")
	fmt.Println("  -v    Verbose mode (show boot order information)")
	fmt.Println("  -h    Show this help message")
	fmt.Println("")
	fmt.Println("Note: This program requires administrative privileges (sudo on Linux, Administrator on Windows)")
}

func rebootToWindows(skipConfirm, persist, verbose bool) {
	fmt.Println("Preparing to reboot to Windows...")
	fmt.Println("")
	fmt.Println("This will:")
	if persist {
		fmt.Println("1. Set Windows as the DEFAULT boot entry (permanent change)")
	} else {
		fmt.Println("1. Set the next boot entry to Windows (one-time)")
	}
	fmt.Println("2. Reboot the system")
	fmt.Println("")

	if !skipConfirm {
		fmt.Print("Do you want to continue? (y/N): ")

		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))

		if response != "y" && response != "yes" {
			fmt.Println("Aborted.")
			return
		}
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

		if verbose {
			fmt.Println("\nCurrent boot configuration:")
			fmt.Println(string(output))
		}

		// Parse for Windows entry and current boot order
		lines := strings.Split(string(output), "\n")
		windowsBootNum := ""
		currentBootOrder := ""
		allBootNums := []string{}

		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), "windows") {
				// Extract boot number (e.g., "Boot0001* Windows Boot Manager")
				if strings.HasPrefix(line, "Boot") && len(line) > 8 {
					windowsBootNum = line[4:8]
					fmt.Printf("Found Windows boot entry: %s\n", line)
				}
			}
			// Parse BootOrder line
			if strings.HasPrefix(line, "BootOrder:") {
				currentBootOrder = strings.TrimSpace(strings.TrimPrefix(line, "BootOrder:"))
				allBootNums = strings.Split(currentBootOrder, ",")
				if verbose {
					fmt.Printf("Current boot order: %s\n", currentBootOrder)
				}
			}
		}

		if windowsBootNum == "" {
			fmt.Println("Could not find Windows boot entry automatically.")
			fmt.Println("Please run 'efibootmgr' to see available boot options")
			return
		}

		if persist {
			// Persistent boot - reorder boot entries with Windows first
			if len(allBootNums) == 0 {
				fmt.Println("Warning: Could not parse current boot order.")
				fmt.Println("Falling back to one-time boot mode...")
				persist = false
			} else {
				// Build new boot order with Windows first, followed by other entries
				newBootOrder := []string{windowsBootNum}
				for _, num := range allBootNums {
					num = strings.TrimSpace(num)
					if num != windowsBootNum && num != "" {
						newBootOrder = append(newBootOrder, num)
					}
				}

				newBootOrderStr := strings.Join(newBootOrder, ",")
				if verbose {
					fmt.Printf("New boot order: %s\n", newBootOrderStr)
				}

				fmt.Println("\n⚠️  WARNING: This will permanently change your default boot order!")
				fmt.Printf("Windows will become the default boot option.\n\n")

				// Set persistent boot order
				cmd = exec.Command("sudo", "efibootmgr", "-o", newBootOrderStr)
				err = cmd.Run()
				if err != nil {
					fmt.Printf("Error setting boot order: %v\n", err)
					return
				}
				fmt.Println("Successfully set Windows as default boot entry!")
			}
		}

		if !persist {
			// One-time boot to Windows
			cmd = exec.Command("sudo", "efibootmgr", "-n", windowsBootNum)
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error setting next boot: %v\n", err)
				return
			}
			fmt.Println("Successfully set next boot to Windows!")
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

func rebootToLinux(skipConfirm, persist, verbose bool) {
	fmt.Println("Preparing to reboot to Linux...")
	fmt.Println("")
	fmt.Println("This will:")
	if persist {
		fmt.Println("1. Set Linux as the DEFAULT boot entry (permanent change)")
	} else {
		fmt.Println("1. Set the next boot entry to Linux (one-time)")
	}
	fmt.Println("2. Reboot the system")
	fmt.Println("")

	if !skipConfirm {
		fmt.Print("Do you want to continue? (y/N): ")

		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))

		if response != "y" && response != "yes" {
			fmt.Println("Aborted.")
			return
		}
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

	if verbose {
		fmt.Println("\nCurrent boot configuration:")
		fmt.Println(string(output))
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

	if linuxIdentifier == "" {
		fmt.Println("Could not find Linux boot entry automatically.")
		fmt.Println("Please run 'bcdedit /enum firmware' as Administrator to see available boot options")
		return
	}

	if persist {
		// Persistent boot - set Linux as first in displayorder
		fmt.Println("\n⚠️  WARNING: This will permanently change your default boot order!")
		fmt.Printf("Linux will become the default boot option.\n\n")

		cmd = exec.Command("bcdedit", "/set", "{fwbootmgr}", "displayorder", linuxIdentifier, "/addfirst")
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error setting boot order: %v\n", err)
			return
		}
		fmt.Println("Successfully set Linux as default boot entry!")
	} else {
		// One-time boot to Linux
		cmd = exec.Command("bcdedit", "/set", "{fwbootmgr}", "bootsequence", linuxIdentifier)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error setting next boot: %v\n", err)
			return
		}
		fmt.Println("Successfully set next boot to Linux!")
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
