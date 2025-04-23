package core

/*
Each pool is a step of the pipeline

=== Download pool ===

A download pool is a set of workers that download documents from a website. 
You can use one or several DownloadPool for one or sereveral websites. Warning : take care of limit rate or banishment from the website.

=====================

*/

import (
	"net/http"
	"io"
	"os"
	_ "strings"
	_ "sync"
	"fmt"
)


type DownloadTask struct {
	url_metadata string
	url_download string
	storage_path string
}

type DownloadPool struct {
	numWorkers int
	numTasks int
	tasks chan DownloadTask
	results chan bool
}


func EprintWorker(wid int, dp *DownloadPool, worker_chan *WorkerErrorChannel) {
	for task := range dp.tasks {
		resp, err := http.Get(task.url_download)
		if (err != nil) {
			CreateWorkerErrorReport(
				fmt.Sprintf("Error fetching URL %s: %v", task.url_download, err),
				wid, worker_chan,
			)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			CreateWorkerErrorReport(
				fmt.Sprintf("Failed to download document %s: status code %d", task.url_download, resp.StatusCode),
				wid, worker_chan,
			)
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			CreateWorkerErrorReport(
				fmt.Sprintf("Error reading response body for URL %s: %v", task.url_download, err),
				wid, worker_chan,
			)
			continue
		}

		file, err := os.Create(task.storage_path)
		if err != nil {
			CreateWorkerErrorReport(
				fmt.Sprintf("Error creating file %s: %v", file, err),
				wid, worker_chan,
			)
			continue
		}
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			CreateWorkerErrorReport(
				fmt.Sprintf("Error writing to file %s: %v", file, err),
				wid, worker_chan,
			)
			continue
		}
	}
	dp.results <- true
}

func CreateDownloadPool(numTasks int, numWorkers int) (*DownloadPool) {
	return &DownloadPool {
		numTasks: numTasks,
		numWorkers: numWorkers,
		tasks: make(chan DownloadTask),
		results: make(chan bool),
	}
}

func CreateDownloadTask(url_metadata string, url_download string, storage_path string) (*DownloadTask) {
	return &DownloadTask {
		url_metadata: url_metadata,
		url_download: url_download,
		storage_path: storage_path,
	}
}