package core

import (
	"net/http"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	baseURL           = "https://eprint.iacr.org"
	endpointComplete  = "/complete"
	endpointDays      = "/days"
	endpointCompact   = "/complete/compact"
	endpointByYear    = "/byyear"
)

var yearLimit = map[int]int{
	2024: 500,
	2025: 450,
}

type DownloadPool struct {
	channel_d chan string
	go_rate int
	wait_time int
}

func createUrl(endpoint string, additional string) (string) {
	if (additional != "") {
		return baseURL + "/" + endpoint + "/" + additional
	}
	return baseURL + "/" + endpoint
}

func createDownloadPool(go_rate int, wait_time int) (*DownloadPool) {
	return &DownloadPool {
		channel_d: make(chan string),
		go_rate: go_rate,
		wait_time: wait_time,
	}
}

func getDocument(url string, storage_path string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download document: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fileName := storage_path + "/" + url[strings.LastIndex(url, "/")+1:]
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func GetDocsPerYears(years []int, storage_folder string) {
    for _, year := range years {
		yearFolder := storage_folder + "/" + fmt.Sprint(year)
		if err := os.MkdirAll(yearFolder, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", yearFolder, err)
			continue
		}

		for i := 1; i <= yearLimit[year]; i++ {
			url := createUrl(fmt.Sprint(year), fmt.Sprintf("%03d.pdf", i))
			
			_, err := getDocument(url, yearFolder)
			if err != nil {
				fmt.Printf("Failed to download %s: %v\n", url, err)
			}
		}
    }
}
