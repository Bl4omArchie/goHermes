package core

import (
	"fmt"
	"os"
	_ "github.com/schollz/progressbar/v3"
)

const (
	baseURL           = "https://eprint.iacr.org"
	endpointComplete  = "/complete"
	endpointDays      = "/days"
	endpointCompact   = "/complete/compact"
	endpointByYear    = "/byyear"
)

func createUrl(endpoint string, additional string) (string) {
	if (additional != "") {
		return baseURL + "/" + endpoint + "/" + additional
	}
	return baseURL + "/" + endpoint
}


func GetDocsPerYears(years []string, storage_folder string) {
	totalDoc := 0
	stats := GetStatistics()

	for _, year := range years {
		yearFolder := storage_folder + "/" + year
		if err := os.MkdirAll(yearFolder, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", yearFolder, err)
			continue
		}
		totalDoc += stats.PapersYear[year]
	}

	err_workers := CreateWorkerErrorChannel()
	d_pool := CreateDownloadPool(totalDoc, 15)

    for i := 1; i <= d_pool.numWorkers; i++ {
        go EprintWorker(i, d_pool, err_workers)
    }

	for _, year := range years {
		for counter := 1; counter <= stats.PapersYear[year]; counter++ {
			url := createUrl(endpointByYear+"/"+year, fmt.Sprintf("%03d.pdf", counter))
			d_pool.tasks <- CreateDownloadTask(url, url, storage_folder+"/"+year+"/"+fmt.Sprintf("%03d.pdf", counter))
		}
	}
	results := make(chan int, d_pool.numTasks)
	for k := 1; k <= totalDoc; k++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}

	GenerateLog(err_workers, "report.txt")
}
