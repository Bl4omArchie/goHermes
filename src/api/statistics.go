package api

import (
	_ "fmt"
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

	// Seek for years a	number of papers per year
	re := regexp.MustCompile(`>(\d{4})</a> \((\d+) papers\)`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	sum := 0
	
	// Create the stat struct
	stats = CreateStats()

	// Fill the struct
	for _, match := range matches {
		if len(match) == 3 {
			docCount, err := strconv.Atoi(match[2])
			utils.CheckError(err)

			stats.years.Add(match[1])
			stats.papersYear[match[1]] = docCount
			sum += docCount
		}
	}

	return stats
}