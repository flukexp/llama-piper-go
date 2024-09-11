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
func CreateVoiceDestinationFolder() {
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
	if _, err := os.Stat("./piper/piper"); err == nil {
		return true
	}
	if _, err := os.Stat("./piper/piper.exe"); err == nil {
		return true
	}
	return false
}

// Determine the file extension based on the filename
func getPiperFileExtension(filename string) string {
	if strings.Contains(filename, "amd64") {
		return ".zip"
	}
	return ".tar.gz"
}

// fetchLatestReleaseURL fetches the redirect URL for the latest release.
func fetchLatestReleaseURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	redirectURL := resp.Request.URL.String()
	return redirectURL, nil
}

// Download file using curl
func DownloadPiperFile(url, filename string) string {

	// Get file extrsnsion
	fileExtension := getPiperFileExtension(filename)

	// Fetching latest release url
	redirectURL, err := fetchLatestReleaseURL(url)
	if err != nil {
		utils.PrintError("Failed to fetch latest release URL")
	}

	// Extract the release tag from the redirect URL
	releaseTag := strings.TrimPrefix(redirectURL, "https://github.com/rhasspy/piper/releases/tag/")
	releaseTag = strings.Trim(releaseTag, "/")

	// Construct the release file name and URL
	releaseFileName := filename + fileExtension
	releaseFileURL := fmt.Sprintf("https://github.com/rhasspy/piper/releases/download/%s/%s", releaseTag, releaseFileName)
	utils.PrintStatus(releaseFileURL)

	cmd := exec.Command("cd", "piper")
	err = cmd.Run()
	if err != nil {
		utils.PrintError("Failed to access piper folder: " + err.Error())
	}

	cmd = exec.Command("curl", "-L", "-o", releaseFileName, releaseFileURL)
	err = cmd.Run()
	if err != nil {
		utils.PrintError("Failed to download file: " + err.Error())
	}

	return releaseFileName
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
