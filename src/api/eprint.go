package api

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"github.com/Bl4omArchie/ePrint-DB/src/db"
	"github.com/Bl4omArchie/ePrint-DB/src/utils"
)

// Download the pdf from the given url
func GetPdf(paper *db.Papers, app *Application) {
	resp, err := http.Get(paper.Doc_url)
	utils.CheckAlertError(err, utils.Error_downloading_document_continue, fmt.Sprintf("Downloading has failed for PDF %d", paper.Id), &app.ac)
	
	defer resp.Body.Close()
}

// Retrieve data such as Category and title
func GetPaperData(paper *db.Papers, app *Application) {
	// Connect to the page
	resp, err := http.Get(paper.Page_url)
	utils.CheckAlertError(err, utils.Error_reach_url_continue, fmt.Sprintf("Failed to reach page for document: %d", paper.Id), &app.ac)
	defer resp.Body.Close()
	
	// Read the content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SendAlert(utils.Error_read_page_content, fmt.Sprintf("Failed to retrieve data for document: %d", paper.Id), &app.ac)
	} else {
		// get title
		re := regexp.MustCompile(`<title>(.*?)</title>`)
		matchTitle := re.FindStringSubmatch(string(body))
		if len(matchTitle) > 1 {
			paper.Title = matchTitle[1]
		} else {
			utils.SendAlert(utils.Error_get_paper_data_continue, fmt.Sprintf("Couldn't find title for document n°%d", paper.Id), &app.ac)
		}
		
		// get category
		re = regexp.MustCompile(`<small class="[^"]+">([^<]+)</small>`)
		matchCategory := re.FindStringSubmatch(string(body))
		if len(matchCategory) > 1 {
			paper.Category = matchCategory[1]
		} else {
			utils.SendAlert(utils.Error_get_paper_data_continue, fmt.Sprintf("Couldn't find category for document n°%d", paper.Id), &app.ac)
		}
		
		// get doc url
	}
}

func GetPaper(id int, year string, app *Application, wg *sync.WaitGroup) {
	defer wg.Done()

	paper := db.Papers{}
	paper.Id = id
	paper.Publication_year = year
	paper.Page_url = Url + year + "/" + strconv.Itoa(id)
	
	GetPaperData(&paper, app)
	//GetPdf(&paper, app)
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
	var wg_download sync.WaitGroup

	for n_year :=0; n_year<len(app.userInput); n_year++ {
		for id := 1; id <= app.stats.papersYear[app.userInput[n_year]]; id++ {
			wg_download.Add(1)
			go GetPaper(id, app.userInput[n_year], app, &wg_download)
		}
		wg_download.Wait()
	}
}