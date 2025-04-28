package core

import (
	"fmt"
	"os"
	_ "github.com/schollz/progressbar/v3"
)

const (
	baseURL           = "https://eprint.iacr.org"
	endpointComplete  = "complete"
	endpointDays      = "days"
	endpointCompact   = "complete/compact"
	endpointByYear    = "byyear"
)

func createUrl(endpoints []string) string {
	base := baseURL
	for _, endpoint := range endpoints {
		base += "/" + endpoint
	}
	return base
}

func GetDocsPerYears(years []string, storage_folder string) {
	totalDoc := 0
	stats := GetStatistics()

	// Create folders for each year
	for _, year := range years {
		yearFolder := storage_folder + year
		if err := os.MkdirAll(yearFolder, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", yearFolder, err)
			continue
		}
		totalDoc += stats.PapersYear[year]
	}

	// Instantiate error channel and download pool
	err_workers := CreateErrorChannel()
	d_pool := CreateDownloadPool(totalDoc, 10)
	go GenerateLog(err_workers)

	// Launch workers
    for i := 1; i <= d_pool.numWorkers; i++ {
        go EprintWorker(i, d_pool, err_workers)
    }

	// Send tasks to workers
	for _, year := range years {
		for counter := 1; counter <= stats.PapersYear[year]; counter++ {
			url_metadata := createUrl([]string{year, fmt.Sprintf("%03d", counter)})
			url_download := createUrl([]string{year, fmt.Sprintf("%03d.pdf", counter)})
			d_pool.tasks <- CreateDownloadTask(url_metadata, url_download, storage_folder+"/"+year+"/"+fmt.Sprintf("%03d.pdf", counter))
		}
	}
	close(d_pool.tasks)

    for k := 1; k <= d_pool.numTasks; k++ {
        _ = <-d_pool.results
    }
}
