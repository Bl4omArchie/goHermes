package api


import (
	"os"
	"sync"
	"strconv"
	"strings"
	"fmt"
	"io"
	"time"
	"bufio"
	"regexp"
	"net/http"
	"github.com/Bl4omArchie/ePrint-DB/src/db"
	"github.com/Bl4omArchie/ePrint-DB/src/utils"
)


type Application struct {
	ac utils.AlertChannel
	stats EprintStatistics
	storage db.Database
	userInput []string
}


func VerifyInput(app *Application) int {
	for _, element := range app.userInput {
		if !app.stats.years.Contains(element) {
			utils.CheckErrorCustom("Category or year not found")
			return 0
		}
	}
	return 1
}

// Download the pdf from the given url
func GetPdf(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	utils.CheckError(err)
	defer resp.Body.Close()
  
}

// Retrieve data such as Category and title
func GetPaperData(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	paper := db.Papers{}

	resp, err := http.Get(url)
	utils.CheckError(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	utils.CheckError(err)

	re := regexp.MustCompile(`@misc{cryptoeprint:[^}]+}`)
	match := re.Find(body)
	utils.RaiseFlag(match)
	
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
	
	start := time.Now()

	for n_year :=0; n_year<len(app.userInput); n_year++ {
		
		for i := 1; i <= app.stats.papersYear[app.userInput[n_year]]; i++ {
			wg_retrieve.Add(1)
			wg_download.Add(1)
			
			go GetPaperData(Url + app.userInput[n_year] + "/" + strconv.Itoa(i), &wg_retrieve)
			go GetPdf(Url + app.userInput[n_year] + "/" + strconv.Itoa(i) + ".pdf", &wg_download)
		}
		wg_retrieve.Wait()
		wg_download.Wait()
	}

	fmt.Println("Temps d'exécution:", time.Since(start))
} 


func LoadApplication() *Application{
	return &Application{
		ac: *utils.CreateAlertChannel(),
		stats: GetStatistics(),
		storage: *db.ConnectDatabase(),
		userInput: []string{},
	}
}


func CloseApplication(app *Application) {
	db.DisconnectDatabase(&app.storage)
	utils.CloseChannel(&app.ac)
}


/*
This function start every features required to make the app works
1) Get statistics
2) Connect database
3) Read command input
4) Launch features
*/
func StartApplication() {
	// Welcome message
	fmt.Println("\033[34m============================================")
	fmt.Println("=== Welcome to ePrint PDF download tool ===")
	fmt.Println("============================================\033[0m")

	//Load the app
	app := LoadApplication()
	
	// Options you have
	fmt.Println("=====================================================================")
	fmt.Println("= -> Write what years or categories you want to be downloaded below")
	fmt.Println("= -> Write 'all' to download every PDF")
	fmt.Println("=====================================================================")

	// Read the user input and clear it
	reader := bufio.NewReader(os.Stdin)
	download_ready := 0

	// Loop until the input is correct
	for download_ready == 0 {
		fmt.Print("Enter option: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		app.userInput = strings.Fields(text)

		download_ready = VerifyInput(app)
	}

	// Start downloading papers
	DownloadPapers(app)


	// Close application
	defer CloseApplication(app)
}