package utility

import (
	"crypto/sha256"
	"net/http"
	"fmt"
	"io"
	"os"
)


func GetPageContent(url string, errChann *ErrorChannel) string {
	resp, err := http.Get(url)
	if (err != nil) {
		CreateErrorReport(fmt.Sprintf("Error fetching URL %s: %v", url, err), errChann)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		CreateErrorReport(fmt.Sprintf("Failed to download document %s: status code %d", url, resp.StatusCode), errChann)
		return ""
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		CreateErrorReport(fmt.Sprintf("Error reading response body for URL %s: %v", url, err), errChann)
		return ""
	}

	return fmt.Sprintf("%x", data)
}

func DownloadDocumentReturnHash(url string, filepath string, errChann *ErrorChannel) string {
	data := GetPageContent(url, errChann)
	file, err := os.Create(filepath)
	if err != nil {
		CreateErrorReport(fmt.Sprintf("Error creating file %s: %v", filepath, err), errChann)
		return ""
	}
	defer file.Close()

	_, err = file.Write([]byte(data))
	if err != nil {
		CreateErrorReport(fmt.Sprintf("Error writing to file %s: %v", filepath, err), errChann)
		return ""
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		CreateErrorReport(fmt.Sprintf("Error seeking file %s: %v", filepath, err), errChann)
		return ""
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		CreateErrorReport(fmt.Sprintf("failed to compute hash: %v", err), errChann)
		return ""
	}

	// Convert byte to string
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
