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
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/Bl4omArchie/ePrint-DB/src/db"
	"github.com/Bl4omArchie/ePrint-DB/src/utils"
)

var (
	url = "https://eprint.iacr.org/"
	
	categories = mapset.NewSet[string](
		"Applications",
		"Cryptographic protocols",
		"Foundations",
		"Implementation",
		"Secret-key cryptography",
		"Public-key cryptography",
		"Attacks and cryptanalysis",)
	
	years = mapset.NewSet[string](
		"2024", "2023", "2022", "2021", "2020", "2019", "2018", "2017", "2016", "2015", "2014", "2013", "2012", "2010", 
		"2009", "2008", "2007", "2006", "2005", "2004", "2003", "2002", "2001", "2000", "1999", "1998", "1997", "1996")
)

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
func DownloadPapers(input_list []string) {
	var wg1 sync.WaitGroup

	url := "https://eprint.iacr.org/"

	start := time.Now()
	for n_year :=0; n_year<len(input_list); n_year++ {
		
		for i := 1; i <= PapersByYear[input_list[n_year]]; i++ {
			time.Sleep(500)
			
			wg1.Add(1)
			
			go RetrieveData(url + input_list[n_year] + "/" + strconv.Itoa(i), &wg1)
		}
		wg1.Wait()
	}

	fmt.Println("Temps d'exécution:", time.Since(start))
} 


func VerifyInput(input []string) int {
	for _, element := range input {
		if !years.Contains(element) {
			utils.CheckErrorCustom("Category or year not found")
			return 0
		}
	}
	return 1
}


func StartApplication() {
	// Welcome message and database connection
	fmt.Println("\n\033[34m === Welcome to ePrint PDF download tool ===\033[0m\n")
	database := db.ConnectDatabase()
	
	// Option for downloading PDF
	fmt.Println("| -> Write what years or categories you want to be downloaded below")
	fmt.Println("| -> Write 'all' to download every PDF\n")
	
	// Read the user input and clear it
	reader := bufio.NewReader(os.Stdin)
	var input_list []string
	download_ready := 0

	for download_ready == 0 {
		fmt.Print("Enter option: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		input_list = strings.Fields(text)

		download_ready = VerifyInput(input_list)
	}

	//Start downloading papers
	DownloadPapers(input_list)

	//Disconnect the DB
	defer db.DisconnectDatabase(database)
}