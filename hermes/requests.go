package hermes

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/html"
)

func GetPageContent(url string, errChan *Log) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		CreateLogReport(fmt.Sprintf("Log fetching URL %s: %v", url, err), errChan)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		CreateLogReport(fmt.Sprintf("Failed to download document %s: status code %d", url, resp.StatusCode), errChan)
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		CreateLogReport(fmt.Sprintf("Log reading response body for URL %s: %v", url, err), errChan)
		return "", err
	}

	return string(data), nil
}

func GetParsedPageContent(url string, errChan *Log) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		CreateLogReport(fmt.Sprintf("Log fetching URL %s: %v", url, err), errChan)
		return &html.Node{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		CreateLogReport(fmt.Sprintf("Failed to download document %s: status code %d", url, resp.StatusCode), errChan)
		return &html.Node{}, err
	}

	return html.Parse(resp.Body)
}

func DownloadDocumentReturnHash(url string, filePath string, errChan *Log) (string, error) {
	data, err := GetPageContent(url, errChan)

	if err != nil {
		CreateLogReport(fmt.Sprintf("Paper has been withdrawn %s: %v", filePath, err), errChan)
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		CreateLogReport(fmt.Sprintf("Error while creating directories for %s: %v", filePath, err), errChan)
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		CreateLogReport(fmt.Sprintf("Error while creating file %s: %v", filePath, err), errChan)
		return "", err
	}
	defer file.Close()

	_, err = file.Write([]byte(data))
	if err != nil {
		CreateLogReport(fmt.Sprintf("Log writing to file %s: %v", filePath, err), errChan)
		return "", err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		CreateLogReport(fmt.Sprintf("Log seeking file %s: %v", filePath, err), errChan)
		return "", err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		CreateLogReport(fmt.Sprintf("failed to compute hash: %v", err), errChan)
		return "", err
	}

	// Convert byte to string
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
