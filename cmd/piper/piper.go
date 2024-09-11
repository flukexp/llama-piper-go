package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/flukexp/llama-piper-go/internal/downloader"
	"github.com/flukexp/llama-piper-go/internal/utils"
)

func voice() {
	// Download piper voices
	downloader.CreateVoiceDestinationFolder()

	if len(os.Args) > 1 {
		langFilter := os.Args[1]
		downloader.DownloadAndExtractVoiceFiles(langFilter)
	} else {
		downloader.DownloadAndExtractVoiceFiles("voice-en-us-amy-low")
	}
}

// GetOSArchitecture determines the OS and architecture, and returns a formatted string.
func GetOSArchitecture() (string, error) {
	var osName, arch string

	switch runtime.GOOS {
	case "darwin":
		osName = "macos"
	case "windows":
		osName = "windows"
	case "linux":
		osName = "linux"
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	switch runtime.GOARCH {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "aarch64"
	case "arm":
		arch = "armv7l"
	default:
		return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	return fmt.Sprintf("%s_%s", osName, arch), nil
}

func main() {
	utils.PrintHeader("Checking Dependencies")

	dependencies := []string{"curl", "node", "npm"}
	for _, dep := range dependencies {
		if !utils.CommandExists(dep) {
			fmt.Printf("%s%s not found. Installing...%s\n", utils.YELLOW, dep, utils.NC)
			utils.InstallPackage(dep)
		} else {
			fmt.Printf("%s%s is already installed.%s\n", utils.GREEN, dep, utils.NC)
		}
	}

	if downloader.IsPiperInstalled() {
		fmt.Println(utils.GREEN + "Piper is already installed." + utils.NC)
	} else {
		architecture, err := GetOSArchitecture()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Determine the download URL based on the architecture and OS
		utils.PrintHeader("Downloading Piper")

		fileName := fmt.Sprintf("piper_%s", architecture)
		fileURL := "https://github.com/rhasspy/piper/releases/latest"
		releaseFileName := downloader.DownloadPiperFile(fileURL, fileName)

		utils.PrintHeader("Extracting Piper")
		// Extract file (piper tar.gz)
		downloader.ExtractTarGz(releaseFileName, "./")
	}

	// Download piper voices
	voice()

	// Installing npm dependencies and starting piper server

	utils.PrintHeader("Installing npm dependencies and starting Piper server")
	// Set working directory to "piper" for npm commands
	workDir := "piper"

	// Run 'npm install' in the "piper" directory
	cmd := exec.Command("npm", "install", ".")
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout // Direct stdout to terminal
	cmd.Stderr = os.Stderr // Direct stderr to terminal
	err := cmd.Run()
	if err != nil {
		utils.PrintError("Failed to run 'npm install': " + err.Error())
		return
	} else {
		utils.PrintStatus("'npm install' completed successfully.")
	}

	// Run 'npm start' in the "piper" directory and display output in terminal
	cmd = exec.Command("npm", "start")
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout // Direct stdout to terminal
	cmd.Stderr = os.Stderr // Direct stderr to terminal
	err = cmd.Run()
	if err != nil {
		utils.PrintError("Failed to run 'npm start': " + err.Error())
		return
	} else {
		utils.PrintStatus("'npm start' completed successfully.")
	}

}
