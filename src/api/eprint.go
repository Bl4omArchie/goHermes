package api

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	_ "strconv"
	"sync"
	"time"
	"github.com/Bl4omArchie/ePrint-DB/src/db"
	"github.com/Bl4omArchie/ePrint-DB/src/utils"
)

// Download the pdf from the given url
func GetPdf(url string, wg *sync.WaitGroup, app *Application) {
	defer wg.Done()

	resp, err := http.Get(url)
	utils.CheckAlertError(err, utils.Error_downloading_document_continue, fmt.Sprintf("Downloading has failed for PDF %s", url), &app.ac)
	defer resp.Body.Close()
}

// Retrieve data such as Category and title
func GetPaperData(url string, id int, year string, wg *sync.WaitGroup, app *Application) *db.Papers{
	defer wg.Done()
	paper := db.Papers{}

	// Connect to the page
	resp, err := http.Get(url)
	utils.CheckAlertError(err, utils.Error_reach_url_continue, fmt.Sprintf("Failed to reach page for document: %d", id), &app.ac)
	defer resp.Body.Close()
	
	// Read the content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SendAlert(utils.Error_read_page_content, fmt.Sprintf("Failed to retrieve data for document: %d", id), &app.ac)
	} else {
		//Get url page and publication year
		paper.Page_url = url
		paper.Publication_year = year

		// get title
		re := regexp.MustCompile(`<title>(.*?)</title>`)
		matchTitle := re.FindStringSubmatch(string(body))
		if len(matchTitle) > 1 {
			paper.Title = matchTitle[1]
		} else {
			utils.SendAlert(utils.Error_get_paper_data_continue, fmt.Sprintf("Couldn't find title for document n°%d", id), &app.ac)
		}
		
		// get category
		re = regexp.MustCompile(`<small class="[^"]+">([^<]+)</small>`)
		matchCategory := re.FindStringSubmatch(string(body))
		if len(matchCategory) > 1 {
			paper.Category = matchCategory[1]
		} else {
			utils.SendAlert(utils.Error_get_paper_data_continue, fmt.Sprintf("Couldn't find category for document n°%d", id), &app.ac)
		}
		
		// get doc url
	}
	return &paper
}


/* 
This function allows you to download papers for each year you want (ie: 2024, 2023)

Algorithm :
1- Launch a new goroutine with GetPaperData
2- Retrieve data such as : authors, category, document extension, document download link
3- Download the PDF
4- Store the binary with the retrieved data into the database
*/
func DownloadPapers(app *Application) {
	var wg_retrieve sync.WaitGroup
	var wg_download sync.WaitGroup

	for n_year :=0; n_year<len(app.userInput); n_year++ {
		for i := 1; i <= app.stats.papersYear[app.userInput[n_year]]; i++ {
			wg_retrieve.Add(1)
			wg_download.Add(1)
			
			go GetPaperData(Url, i, app.userInput[n_year], &wg_retrieve, app)
			time.Sleep(500)
			//go GetPdf(Url + app.userInput[n_year] + "/" + strconv.Itoa(i) + ".pdf", &wg_download, app)
		}
		wg_retrieve.Wait()
		wg_download.Wait()
	}
}