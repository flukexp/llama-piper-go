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

//go:embed llama/*
var llamaFiles embed.FS

// Check if the piper files have already been extracted
func arePiperFilesExtracted(destination string) bool {
	piperFolder := filepath.Join(destination, "piper")
	_, err := os.Stat(piperFolder)
	return !os.IsNotExist(err)
}

// Check if the llama files have already been extracted
func areLlamaFilesExtracted(destination string) bool {
	llamaFolder := filepath.Join(destination, "llama")
	_, err := os.Stat(llamaFolder)
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

// Extract the embedded llama files to a destination directory
func extractLlamaFiles(destination string) error {
	// Ensure the destination directory exists
	llamaFolder := filepath.Join(destination, "llama")
	err := os.MkdirAll(llamaFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create llama folder: %w", err)
	}

	// Read all entries in the embedded "piper" folder
	entries, err := llamaFiles.ReadDir("llama")
	if err != nil {
		return fmt.Errorf("failed to read embedded llama folder: %w", err)
	}

	for _, entry := range entries {
		// For each file in the embedded folder, extract it to the destination
		filePath := filepath.Join(llamaFolder, entry.Name())
		data, err := llamaFiles.ReadFile("llama/" + entry.Name())
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
		// Running the piper-installer.bat script
		cmd := exec.Command("cmd", "/c", "start", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(RED + "Error running installer.bat: " + err.Error() + NC)
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

	// Check if Llama files are already extracted
	if !areLlamaFilesExtracted(executableDir) {
		// Extract the llama files to the directory where the executable was clicked
		fmt.Println(BLUE + "Extracting Llama files..." + NC)
		err = extractLlamaFiles(executableDir)
		if err != nil {
			fmt.Println(RED + "Failed to extract Llama files: " + err.Error() + NC)
			return
		}
		fmt.Println(GREEN + "Llama files extracted successfully." + NC)
	} else {
		fmt.Println(GREEN + "Llama files already extracted." + NC)
	}

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
	llamaPath := filepath.Join(executableDir, "llama", "llama-installer.bat")

	// Check if the batch file exists
	if _, err := os.Stat(llamaPath); os.IsNotExist(err) {
		fmt.Println(RED + "piper-installer.bat not found in the extracted files." + NC)
		return
	}

	llamaCommand := llamaPath

	fmt.Printf("%sExecuting command: %s%s\n", BLUE, llamaCommand, NC)
	openNewTerminal(llamaCommand)
	fmt.Printf("%sStarting Llama Installer...%s\n", GREEN, NC)

	// Construct the path to the piper-installer.bat script
	batFilePath := filepath.Join(executableDir, "piper", "piper-installer.bat")

	// Check if the batch file exists
	if _, err := os.Stat(batFilePath); os.IsNotExist(err) {
		fmt.Println(RED + "piper-installer.bat not found in the extracted files." + NC)
		return
	}

	// Construct the command to run the .bat file
	piperCommand := batFilePath

	fmt.Printf("%sExecuting command: %s%s\n", BLUE, piperCommand, NC)
	openNewTerminal(piperCommand)
	fmt.Printf("%sStarting Piper Installer...%s\n", GREEN, NC)

}
