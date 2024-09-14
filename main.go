package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Define colors for output
const (
	RED    = "\033[0;31m"
	GREEN  = "\033[0;32m"
	YELLOW = "\033[1;33m"
	BLUE   = "\033[0;34m"
	NC     = "\033[0m" // No Color
)

//go:embed piper/*
var piperFiles embed.FS

// Check if the piper files have already been extracted
func arePiperFilesExtracted(destination string) bool {
	piperFolder := filepath.Join(destination, "piper")
	_, err := os.Stat(piperFolder)
	return !os.IsNotExist(err)
}

// Extract the embedded piper files to a destination directory
func extractPiperFiles(destination string) error {
	// Ensure the destination directory exists
	piperFolder := filepath.Join(destination, "piper")
	err := os.MkdirAll(piperFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create piper folder: %w", err)
	}

	// Read all entries in the embedded "piper" folder
	entries, err := piperFiles.ReadDir("piper")
	if err != nil {
		return fmt.Errorf("failed to read embedded piper folder: %w", err)
	}

	for _, entry := range entries {
		// For each file in the embedded folder, extract it to the destination
		filePath := filepath.Join(piperFolder, entry.Name())
		data, err := piperFiles.ReadFile("piper/" + entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", entry.Name(), err)
		}

		// Write the file to the destination directory
		err = os.WriteFile(filePath, data, 0755)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}
	return nil
}

// Open a new terminal and run commands
func openNewTerminal(command string) {
	switch runtime.GOOS {
	case "windows":
		// Directly running the .bat script in Windows
		cmd := exec.Command("cmd", "/c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(RED + "Error running batch file: " + err.Error() + NC)
		}

	default:
		fmt.Println(RED + "Unsupported OS for opening new terminals." + NC)
	}
}

func main() {
	// Get the directory where the executable is located
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println(RED + "Failed to get executable path: " + err.Error() + NC)
		return
	}
	executableDir := filepath.Dir(executablePath)

	// Check if Piper files are already extracted
	if !arePiperFilesExtracted(executableDir) {
		// Extract the piper files to the directory where the executable was clicked
		fmt.Println(BLUE + "Extracting Piper files..." + NC)
		err = extractPiperFiles(executableDir)
		if err != nil {
			fmt.Println(RED + "Failed to extract Piper files: " + err.Error() + NC)
			return
		}
		fmt.Println(GREEN + "Piper files extracted successfully." + NC)
	} else {
		fmt.Println(GREEN + "Piper files already extracted." + NC)
	}

	// Construct the path to the piper-installer.bat script
	batFilePath := filepath.Join(executableDir, "piper", "piper-installer.bat")

	// Check if the batch file exists
	if _, err := os.Stat(batFilePath); os.IsNotExist(err) {
		fmt.Println(RED + "piper-installer.bat not found in the extracted files." + NC)
		return
	}

	// Construct the command to run the .bat file
	command := batFilePath

	fmt.Printf("%sExecuting command: %s%s\n", BLUE, command, NC)
	openNewTerminal(command)
	fmt.Printf("%sStarting Piper Installer...%s\n", GREEN, NC)
}
