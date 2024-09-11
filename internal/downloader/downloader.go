package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/flukexp/llama-piper-go/internal/constants"
	"github.com/flukexp/llama-piper-go/internal/utils"
)

// Create the destination folder if it doesn't exist
func CreateDestinationFolder() {
	utils.PrintHeader("Checking Destination Folder")
	if _, err := os.Stat(constants.DestinationFolder); os.IsNotExist(err) {
		utils.PrintStatus("Creating destination folder: " + constants.DestinationFolder)
		if err := os.Mkdir(constants.DestinationFolder, os.ModePerm); err != nil {
			utils.PrintError("Failed to create destination folder: " + err.Error())
		}
	} else {
		utils.PrintStatus("Destination folder " + constants.DestinationFolder + " already exists.")
	}
}

// Check if a file is already installed
func isVoiceInstalled(fileBaseName string) bool {
	onnxPath := filepath.Join(constants.DestinationFolder, fileBaseName+".onnx")
	jsonPath := filepath.Join(constants.DestinationFolder, fileBaseName+".onnx.json")
	if _, err := os.Stat(onnxPath); err == nil {
		return true
	}
	if _, err := os.Stat(jsonPath); err == nil {
		return true
	}
	return false
}

// Check if Piper is already installed
func IsPiperInstalled() bool {
	if _, err := os.Stat("./piper"); err == nil {
		return true
	}
	if _, err := os.Stat("./piper.exe"); err == nil {
		return true
	}
	return false
}

// Download file using curl
func DownloadPiperFile(url, filename string) {
	cmd := exec.Command("curl", "-L", "-o", filename, url)
	err := cmd.Run()
	utils.PrintError(err.Error())
}

// Download a file
func downloadVoiceFile(fileName string) error {
	fileURL := fmt.Sprintf("%s/%s", constants.URL, fileName)
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return nil
}

// Download and extract files
func DownloadAndExtractVoiceFiles(langFilter string) {
	for _, file := range constants.Files {
		if strings.Contains(file, langFilter) {
			fileBaseName := strings.TrimPrefix(strings.TrimSuffix(file, ".tar.gz"), "voice-")

			if isVoiceInstalled(fileBaseName) {
				utils.PrintStatus(file + " is already installed.")
			} else {
				utils.PrintHeader("Downloading " + file)
				err := downloadVoiceFile(file)
				if err != nil {
					utils.PrintError("Failed to download " + file + ": " + err.Error())
					continue
				}

				utils.PrintHeader("Extracting " + file)
				err = ExtractTarGz(file, constants.DestinationFolder)
				if err != nil {
					utils.PrintError("Failed to extract " + file + ": " + err.Error())
				} else {
					utils.PrintStatus("Extracted " + file)
				}

				// Clean up
				os.Remove(file)
				utils.PrintStatus("Cleaned up " + file)
			}
		}
	}
}
