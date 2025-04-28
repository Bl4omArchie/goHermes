package core

/*
Each pool is a step of the pipeline

=== Download pool ===

A download pool is a set of workers that download documents from a website. 
You can use one or several DownloadPool for one or several websites. Warning : take care of limit rate or banishment from the website.

=====================
*/

import (
	"net/http"
	"io"
	"os"
	"fmt"
	"crypto/sha256"
)

type DownloadTask struct {
	urlMetadata string
	urlDownload string
	storagePath string
}

// TODO : the result struct has to be rework
type DownloadResult struct {
	status int
}

type DownloadPool struct {
	numWorkers int
	numTasks int
	tasks chan *DownloadTask
	results chan *DownloadResult
}

func getPage(url string, wid int, wec *ErrorChannel) string {
	resp, err := http.Get(url)
	if (err != nil) {
		CreateErrorReport(fmt.Sprintf("Error fetching URL %s: %v", url, err), wid, wec)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		CreateErrorReport(fmt.Sprintf("Failed to download document %s: status code %d", url, resp.StatusCode), wid, wec)
		return ""
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		CreateErrorReport(fmt.Sprintf("Error reading response body for URL %s: %v", url, err),wid, wec)
		return ""
	}

	return fmt.Sprintf("%x", data)
}

func downloadPage(url string, filepath string, wid int, wec *ErrorChannel) string {
	data := getPage(url, wid, wec)
	file, err := os.Create(filepath)
	if err != nil {
		CreateErrorReport(fmt.Sprintf("Error creating file %s: %v", filepath, err), wid, wec)
		return ""
	}
	defer file.Close()

	_, err = file.Write([]byte(data))
	if err != nil {
		CreateErrorReport(fmt.Sprintf("Error writing to file %s: %v", filepath, err), wid, wec)
		return ""
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		CreateErrorReport(fmt.Sprintf("Error seeking file %s: %v", filepath, err), wid, wec)
		return ""
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		CreateErrorReport(fmt.Sprintf("failed to compute hash: %v", err), wid, wec)
		return ""
	}

	// Convert byte to string
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func EprintWorker(wid int, dp *DownloadPool, wec *ErrorChannel) {
	for task := range dp.tasks {
		data := getPage(task.urlMetadata, wid, wec)
		if data != "" {
			hash_pdf := downloadPage(task.urlDownload, task.storagePath, wid, wec)
			if hash_pdf != "" {
				dp.results <- &DownloadResult{status: 1}
			} else {
				dp.results <- &DownloadResult{status: 0}
			}
		} else {
			dp.results <- &DownloadResult{status: 0}
		}
	}
}

func CreateDownloadPool(numTasks int, numWorkers int) (*DownloadPool) {
	return &DownloadPool {
		numTasks: numTasks,
		numWorkers: numWorkers,
		tasks: make(chan *DownloadTask, numTasks),
		results: make(chan *DownloadResult, numTasks),
	}
}

func CreateDownloadTask(urlMetadata string, urlDownload string, storagePath string) (*DownloadTask) {
	return &DownloadTask {
		urlMetadata: urlMetadata,
		urlDownload: urlDownload,
		storagePath: storagePath,
	}
}
