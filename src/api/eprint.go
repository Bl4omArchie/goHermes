package api


import (
	"os"
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
	
	tags = mapset.NewSet[string](
		"Applications",
		"Cryptographic protocols",
		"Foundations",
		"Implementation",
		"Secret-key cryptography",
		"Public-key cryptography",
		"Attacks and cryptanalysis",
		"2024", "2023", "2022", "2021", "2020", "2019", "2018", "2017", "2016", "2015", "2014", "2013", "2012", "2010", 
		"2009", "2008", "2007", "2006", "2005", "2004", "2003", "2002", "2001", "2000", "1999", "1998", "1997", "1996")
)

type Papers struct {
    Title string
    Link string
    Publication_year int
	Category string
	File_data string
}

func RetrieveDataPaper(url string) {
	paper := Papers{}

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

func DownloadPapers(input_list []string) {
	url = "https://eprint.iacr.org/2024/"

	start := time.Now()
	for i:=1; i<1799; i++ {
		RetrieveDataPaper(url+strconv.Itoa(i))
	}
	fmt.Println("Temps d'exécution:", time.Since(start))
}

func VerifyInput(input []string) int {
	for _, element := range input {
		if !tags.Contains(element) {
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

	// Start downloading
	DownloadPapers(input_list)

	// Disconnect the DB and quit the program
	defer db.DisconnectDatabase(database)
}