package core


import (
	"crypto/sha256"
	"path/filepath"
	"net/http"
	"fmt"
	"io"
	"os"
)


func GetPageContent(url string, errChann *LogChannel) (string, error) {
	resp, err := http.Get(url)
	if (err != nil) {
		CreateLogReport(fmt.Sprintf("Log fetching URL %s: %v", url, err), errChann)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		CreateLogReport(fmt.Sprintf("Failed to download document %s: status code %d", url, resp.StatusCode), errChann)
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		CreateLogReport(fmt.Sprintf("Log reading response body for URL %s: %v", url, err), errChann)
		return "", err
	}

	return string(data), nil
}

func DownloadDocumentReturnHash(url string, filePath string, errChann *LogChannel) (string, error) {
	data, _ := GetPageContent(url, errChann)

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		CreateLogReport(fmt.Sprintf("Log creating directories for %s: %v", filePath, err), errChann)
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		CreateLogReport(fmt.Sprintf("Log creating file %s: %v", filePath, err), errChann)
		return "", err
	}
	defer file.Close()

	_, err = file.Write([]byte(data))
	if err != nil {
		CreateLogReport(fmt.Sprintf("Log writing to file %s: %v", filePath, err), errChann)
		return "", err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		CreateLogReport(fmt.Sprintf("Log seeking file %s: %v", filePath, err), errChann)
		return "", err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		CreateLogReport(fmt.Sprintf("failed to compute hash: %v", err), errChann)
		return "", err
	}

	// Convert byte to string
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
