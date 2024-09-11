package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/flukexp/llama-piper-go/internal/utils"
)

// Determine if running under WSL
func isWSL() bool {
	unameOut, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(unameOut), "Microsoft") || os.Getenv("WSLENV") != ""
}

// Open a new terminal and run commands
func openNewTerminal(command string) {
	switch runtime.GOOS {
	case "darwin":
		// macOS - using osascript to open new Terminal windows
		osascript := fmt.Sprintf(`
		tell application "Terminal"
			activate
			do script "%s"
		end tell
		`, command)
		_, err := exec.Command("osascript", "-e", osascript).Output()
		if err != nil {
			utils.PrintError(err.Error())
		}
	case "linux":
		if isWSL() {
			// WSL - using cmd.exe to run commands in new terminal instances
			_, err := exec.Command("cmd", "/c", "start", "wsl.exe", "-e", "bash", "-c", command).Output()
			if err != nil {
				utils.PrintError(err.Error())
			}
		} else {
			utils.PrintError("Unsupported OS for opening new terminals.")
		}
	default:
		utils.PrintError("Unsupported OS for opening new terminals.")
	}
}

func main() {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		utils.PrintError("Failed to get current working directory: " + err.Error())
		return
	}

	// Construct the command with the current directory and the script path
	command := fmt.Sprintf("cd %s && go run ./cmd/piper/piper.go", dir)

	fmt.Printf("%sExecuting command: %s%s\n", utils.BLUE, command, utils.NC)
	openNewTerminal(command)
	fmt.Printf("%sStarting %s server...%s\n", utils.GREEN, command, utils.NC)
}
