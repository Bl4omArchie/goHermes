package core

import (
	"net/http"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"github.com/schollz/progressbar/v3"
)

const (
	baseURL           = "https://eprint.iacr.org"
	endpointComplete  = "/complete"
	endpointDays      = "/days"
	endpointCompact   = "/complete/compact"
	endpointByYear    = "/byyear"
)

type DownloadPool struct {
	channel_d chan string
	semaphore chan struct{}
	semaphore_count int
	wg sync.WaitGroup
}


func createUrl(endpoint string, additional string) (string) {
	if (additional != "") {
		return baseURL + "/" + endpoint + "/" + additional
	}
	return baseURL + "/" + endpoint
}

func createDownloadPool(limit_rate int) (*DownloadPool) {
	return &DownloadPool {
		channel_d: make(chan string),
		semaphore_count: limit_rate,
		semaphore: make(chan struct{}, limit_rate),
		wg: sync.WaitGroup{},
	}
}

func getDocument(url string, storage_path string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching URL %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//fmt.Printf("Failed to download document %s: status code %d\n", url, resp.StatusCode)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body for URL %s: %v\n", url, err)
		return
	}

	fileName := storage_path + "/" + url[strings.LastIndex(url, "/")+1:]
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", fileName, err)
		return
	}
}

func GetDocsPerYears(years []string, storage_folder string) {
	d_pool := createDownloadPool(15)
	stats := GetStatistics()

	for _, year := range years {
		// Start a loading bar for a year
        progressBar := progressbar.NewOptions(stats.PapersYear[year],
            progressbar.OptionSetDescription(fmt.Sprintf("Downloading papers for year %s", year)),
            progressbar.OptionShowCount(),
            progressbar.OptionSetWidth(20),
        )
		
		// Create folders for each years
		yearFolder := storage_folder + "/" + year
		if err := os.MkdirAll(yearFolder, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", yearFolder, err)
			continue
		}
		
		//Download papers
		for counter := 1; counter <= stats.PapersYear[year]; counter++ {
			d_pool.wg.Add(1)
			d_pool.semaphore <- struct{}{}

			url := createUrl(year, fmt.Sprintf("%03d.pdf", counter))
			go func(url, folder string, bar *progressbar.ProgressBar) {
				defer func() { 
					<-d_pool.semaphore
					bar.Add(1)
				}()
				getDocument(url, folder, &d_pool.wg)
			}(url, yearFolder, progressBar)
		}
	}
	d_pool.wg.Wait()
}
