package api


import (
	"os"
	"fmt"
	"bufio"
	"strings"
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

func Menu(app *Application) {
	// Welcome message
	fmt.Println("\033[34m===========================================")
	fmt.Println("=== Welcome to ePrint PDF download tool ===")
	fmt.Println("===========================================\033[0m")

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

	DownloadPapers(app)

}

/* Loading the application consist of :
1- Create the alert channel
2- Get statistics from ePrint website
3- Connection to database
4- Initiate the user input buffer
*/
func LoadApplication() *Application {
	ac := *utils.CreateAlertChannel()
	go utils.ListenerAlertChannel(&ac)
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
	// Load the app : database connection, alert listener, update the statistics...
	app := LoadApplication()
	
	// Start the menu
	Menu(app)

	// Close application
	defer CloseApplication(app)
}