package api

import (
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Bl4omArchie/ePrint-DB/src/utils"
	mapset "github.com/deckarep/golang-set/v2"
)

var (
	Url = "https://eprint.iacr.org/"
	Url_by_years = "https://eprint.iacr.org/byyear"
)

type EprintStatistics struct {
	totalDocuments int				//Total amount of documents
	papersYear map[string]int		//for each years, the number of documents
	categories mapset.Set[string]	//an array of every avvailable categories
	years mapset.Set[string]		//an array of every available years
}

func CreateStats() (stats EprintStatistics) {
	stats.totalDocuments = 0
	stats.papersYear = make(map[string]int)
	stats.categories = mapset.NewSet[string]()
	stats.years = mapset.NewSet[string]()

	return stats
}

func GetStatistics() (stats EprintStatistics) {
	// get the page where you can find stats we want
	resp, err := http.Get(Url_by_years)
	utils.CheckError(err)
	defer resp.Body.Close()

	// Read the body page
	body, err := io.ReadAll(resp.Body)
	utils.CheckError(err)

	// Seek for years and the number of papers per year
	re_years := regexp.MustCompile(`>(\d{4})</a> \((\d+) papers\)`)
	matches_years := re_years.FindAllStringSubmatch(string(body), -1)
	
	// Seek for categories
	re_categories := regexp.MustCompile(`<a href="/search\?category=[^"]+">([^<]+)</a>`)
	matches_categories := re_categories.FindAllStringSubmatch(string(body), -1)
	
	// Create the stat struct
	stats = CreateStats()

	sum := 0
	// Fill the struct with years
	for _, match := range matches_years {
		if len(match) == 3 {
			docCount, err := strconv.Atoi(match[2])
			utils.CheckError(err)

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