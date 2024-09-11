package downloader

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/flukexp/llama-piper-go/internal/utils"
)

// ExtractTarGz extracts a tar.gz file to the specified destination folder
func ExtractTarGz(fileName string, destination string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar file: %w", err)
		}

		target := filepath.Join(destination, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			outFile.Close()
		}
	}
	// Clean up
	os.Remove(fileName)
	utils.PrintStatus("Cleaned up " + fileName)

	return nil
}

// Unzip files
func unzipFile(zipFile, dest string) {
	r, err := zip.OpenReader(zipFile)
	utils.PrintError(err.Error())
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				utils.PrintError(err.Error())
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			utils.PrintError(err.Error())

			rc, err := f.Open()
			utils.PrintError(err.Error())

			_, err = io.Copy(outFile, rc)
			utils.PrintError(err.Error())

			outFile.Close()
			rc.Close()
		}
	}
}
