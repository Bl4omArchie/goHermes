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
			utils.SendAlert(0xc4, "Incorrect input. You must use a valid year.", &app.ac)
			return 0
		}
	}
	return 1
}

// Download the pdf from the given url
func GetPdf(url string, wg *sync.WaitGroup, app *Application) {
	defer wg.Done()

	resp, err := http.Get(url)
	utils.CheckAlertError(err, 0xc2, fmt.Sprintf("Downloading has failed for PDF %d", url), &app.ac)
	defer resp.Body.Close()
  
}

// Retrieve data such as Category and title
func GetPaperData(url string, wg *sync.WaitGroup, app *Application) {
	defer wg.Done()
	paper := db.Papers{}

	resp, err := http.Get(url)
	utils.CheckAlertError(err, 0xc3, fmt.Sprintf("Failed to reach page: %d", url), &app.ac)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	utils.CheckAlertError(err, 0xc3, fmt.Sprintf("Failed to retrieve data for page: %d", url), &app.ac)

	re := regexp.MustCompile(`@misc{cryptoeprint:[^}]+}`)
	match := re.Find(body)
	if len(match) <= 0 {
		utils.SendAlert(0xc3, fmt.Sprintf("Couldn't find cryptoeprint for PDF %d", url), &app.ac)
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
	
	start := time.Now()

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

	fmt.Println("Temps d'exécution:", time.Since(start))
} 


/* Loading the application consist of :
1- Create the alert channel
2- Get statistics from ePrint website
3- Connection to database
4- Initiate the user input buffer
*/
func LoadApplication() *Application {
	ac := *utils.CreateAlertChannel()
	stats := *GetStatistics(&ac)
	storage := *db.ConnectDatabase(&ac)

	return &Application{
		ac: ac,
		stats: stats,
		storage: storage,
		userInput: []string{},
	}
}


func CloseApplication(app *Application) {
	db.DisconnectDatabase(&app.ac, &app.storage)
	utils.CloseChannel(&app.ac)
}


func StartApplication() {
	// Welcome message
	fmt.Println("\033[34m============================================")
	fmt.Println("=== Welcome to ePrint PDF download tool ===")
	fmt.Println("============================================\033[0m")

	//Load the app : database connection, alert listener, update the statistics...
	app := LoadApplication()
	
	// Options you have
	fmt.Println("=======================================================")
	fmt.Println("= -> Write what years you want to be downloaded below")
	fmt.Println("= -> Write 'all' to download every PDF")
	fmt.Println("=======================================================")

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