package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Define colors
const (
	RED    = "\033[0;31m"
	GREEN  = "\033[0;32m"
	YELLOW = "\033[1;33m"
	BLUE   = "\033[0;34m"
	NC     = "\033[0m" // No Color
)

// Print section headers
func PrintHeader(header string) {
	fmt.Printf("%s==================== %s ====================%s\n", BLUE, header, NC)
}

// Print errors
func PrintError(err string) {
	fmt.Printf("%sError: %s%s\n", RED, err, NC)
}

// Print status
func PrintStatus(status string) {
	fmt.Printf("%s%s%s\n", GREEN, status, NC)
}

// Check if a command exists
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Install a package
func InstallPackage(pkg string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// macOS installation with brew
		if !CommandExists("brew") {
			fmt.Printf("%sHomebrew not found. Installing Homebrew...%s\n", YELLOW, NC)
			cmd = exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
		} else {
			cmd = exec.Command("brew", "install", pkg)
		}
	case "linux":
		cmd = exec.Command("sudo", "apt-get", "install", "-y", pkg)
	case "windows":
		fmt.Printf("%sPackage installation not supported on Windows in this script%s\n", RED, NC)
		return
	default:
		fmt.Printf("%sUnsupported OS%s\n", RED, NC)
		return
	}

	err := cmd.Run()
	if err != nil {
		PrintError("Failed to install " + pkg + ": " + err.Error())
	} else {
		PrintStatus(pkg + " installed successfully!")
	}
}
