package api


import (
	"io"
	"fmt"
	"sync"
	"regexp"
	"strconv"
	"net/http"
	"github.com/Bl4omArchie/ePrint-DB/src/db"
	"github.com/Bl4omArchie/ePrint-DB/src/utils"
)


// Download the pdf from the given url
func GetPdf(url string, wg *sync.WaitGroup, app *Application) {
	defer wg.Done()

	resp, err := http.Get(url)
	utils.CheckAlertError(err, 0xc2, fmt.Sprintf("Downloading has failed for PDF %s", url), &app.ac)
	defer resp.Body.Close()
  
}

// Retrieve data such as Category and title
func GetPaperData(url string, wg *sync.WaitGroup, app *Application) {
	defer wg.Done()
	paper := db.Papers{}

	resp, err := http.Get(url)
	utils.CheckAlertError(err, 0xc3, fmt.Sprintf("Failed to reach page: %s", url), &app.ac)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	utils.CheckAlertError(err, 0xc3, fmt.Sprintf("Failed to retrieve data for page: %s", url), &app.ac)

	re := regexp.MustCompile(`@misc{cryptoeprint:[^}]+}`)
	match := re.Find(body)
	if len(match) <= 0 {
		utils.SendAlert(0xc3, fmt.Sprintf("Couldn't find cryptoeprint for PDF %s", url), &app.ac)
	}
	
	re = regexp.MustCompile(`<small class="[^"]+">([^<]+)</small>`)
	matchCategory := re.FindStringSubmatch(string(body))
	
	if len(matchCategory) > 1 {
		paper.Category = matchCategory[1]
	}
}

// Take a list of years (ie: 2024, 2023 ...) and launch the stages of data retrieve and pdf download
func DownloadPapers(app *Application) {
	var wg_retrieve sync.WaitGroup
	var wg_download sync.WaitGroup

	for n_year :=0; n_year<len(app.userInput); n_year++ {
		
		for i := 1; i <= app.stats.papersYear[app.userInput[n_year]]; i++ {
			wg_retrieve.Add(1)
			wg_download.Add(1)

			go GetPaperData(Url + app.userInput[n_year] + "/" + strconv.Itoa(i), &wg_retrieve, app)
			go GetPdf(Url + app.userInput[n_year] + "/" + strconv.Itoa(i) + ".pdf", &wg_download, app)
		}
		wg_retrieve.Wait()
		wg_download.Wait()
	}
}