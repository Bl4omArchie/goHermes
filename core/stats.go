package core

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	mapset "github.com/deckarep/golang-set/v2"
)


type EprintStatistics struct {
	totalDocuments int
	papersYear map[string]int
	categories mapset.Set[string]
	years mapset.Set[string]
}


func CreateStats() *EprintStatistics {
	return &EprintStatistics {
		totalDocuments: 0,
		papersYear: make(map[string]int),
		categories: mapset.NewSet[string](),
		years: mapset.NewSet[string](),
	}
}


func GetStatistics() *EprintStatistics {
	// get the page where you can find stats we want
	resp, _ := http.Get(createUrl(endpointByYear, ""))
	defer resp.Body.Close()

	// Read the body page
	body, _ := io.ReadAll(resp.Body)

	// Seek for years and the number of papers per year
	re_years := regexp.MustCompile(`>(\d{4})</a> \((\d+) papers\)`)
	matches_years := re_years.FindAllStringSubmatch(string(body), -1)
	
	// Seek for categories
	re_categories := regexp.MustCompile(`<a href="/search\?category=[^"]+">([^<]+)</a>`)
	matches_categories := re_categories.FindAllStringSubmatch(string(body), -1)
	
	stats := CreateStats()

	sum := 0
	// Fill the struct with years
	for _, match := range matches_years {
		if len(match) == 3 {
			docCount, _ := strconv.Atoi(match[2])
			stats.years.Add(match[1])
			stats.papersYear[match[1]] = docCount
			sum += docCount
		}
	}
	stats.totalDocuments = sum

	// Fill the struct with categories
	for _, match := range matches_categories {
		if len(match) == 2 {
			stats.categories.Add(match[1])
		}
	}

	return stats
}