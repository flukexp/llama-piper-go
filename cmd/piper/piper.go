package piper

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/flukexp/llama-piper-go/internal/downloader"
	"github.com/flukexp/llama-piper-go/internal/utils"
)

func piper() {
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
		architecture := runtime.GOARCH
		// Determine the download URL based on the architecture and OS
		utils.PrintHeader("Downloading Piper")

		fileURL := fmt.Sprintf("https://github.com/rhasspy/piper/releases/download/latest/piper_%s.tar.gz", architecture)
		downloader.DownloadPiperFile(fileURL, "piper.tar.gz")

		utils.PrintHeader("Extracting Piper")
		// Extract file (assuming tar.gz)
		downloader.ExtractTarGz("piper.tar.gz", "./")
	}

	// Download piper voices
	voice()

	// Installing npm dependencies and starting piper server

	utils.PrintHeader("Installing npm dependencies and starting Piper server")
	cmd := exec.Command("npm", "install", ".")
	err := cmd.Run()
	utils.PrintError(err.Error())

	cmd = exec.Command("npm", "start")
	err = cmd.Run()
	utils.PrintError(err.Error())

}
