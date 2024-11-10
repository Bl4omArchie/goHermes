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


func VerifyInput(input []string, stats EprintStatistics) int {
	for _, element := range input {
		if !stats.years.Contains(element) {
			utils.CheckErrorCustom("Category or year not found")
			return 0
		}
	}
	return 1
}

// Download the pdf from the given url
func GetPdf(url string, wg *sync.WaitGroup) {
	defer wg.Done()
}

// Retrieve data such as Category and title
func RetrieveData(url string, wg *sync.WaitGroup) {
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
func DownloadPapers(input_list []string, stats EprintStatistics) {
	var wg1 sync.WaitGroup

	start := time.Now()
	for n_year :=0; n_year<len(input_list); n_year++ {
		
		for i := 1; i <= stats.papersYear[input_list[n_year]]; i++ {
			time.Sleep(500)
			
			wg1.Add(1)
			
			go RetrieveData(Url + input_list[n_year] + "/" + strconv.Itoa(i), &wg1)
		}
		wg1.Wait()
	}

	fmt.Println("Temps d'exécution:", time.Since(start))
} 


/*
This function start every features required to make the app works
1) Get statistics
2) Connect database
3) Read command input
4) Launch features
*/
func StartApplication() {
	// Welcome message and database connection
	fmt.Println("\033[34m============================================")
	fmt.Println("=== Welcome to ePrint PDF download tool ===")
	fmt.Println("============================================\033[0m")
	database := db.ConnectDatabase()
	
	// Options you have
	fmt.Println("=====================================================================")
	fmt.Println("= -> Write what years or categories you want to be downloaded below")
	fmt.Println("= -> Write 'all' to download every PDF")
	fmt.Println("=====================================================================")
	
	// Get statistics about ePrint
	stats := GetStatistics()

	// Read the user input and clear it
	reader := bufio.NewReader(os.Stdin)
	var input_list []string
	download_ready := 0

	// Loop until the input is correct
	for download_ready == 0 {
		fmt.Print("Enter option: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		input_list = strings.Fields(text)

		download_ready = VerifyInput(input_list, stats)
	}

	// Start downloading papers
	DownloadPapers(input_list, stats)

	// Disconnect the DB
	defer db.DisconnectDatabase(database)
}