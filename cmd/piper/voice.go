package piper

import (
	"os"

	"github.com/flukexp/llama-piper-go/internal/downloader"
)

func voice() {

	// Download piper voices
	downloader.CreateDestinationFolder()

	if len(os.Args) > 1 {
		langFilter := os.Args[1]
		downloader.DownloadAndExtractVoiceFiles(langFilter)
	} else {
		downloader.DownloadAndExtractVoiceFiles("voice-en-us-amy-low")
	}

}
