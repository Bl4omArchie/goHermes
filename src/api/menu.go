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

var (
	welcome_message = "\033[34m===========================================\n=== Welcome to ePrint PDF download tool ===\n===========================================\033[0m"
	menu_option = "1- Create database\n2- Download papers\n"
	download_option = "=======================================================\n= -> Write what years you want to be downloaded below\n= -> Write 'all' to download every PDF\n======================================================="
)


func verifyInput(app *Application) int {
	for _, element := range app.userInput {
		if !app.stats.years.Contains(element) {
			utils.SendAlert(utils.Error_user_input_continue, "Incorrect input. You must use a valid year.", &app.ac)
			return 0
		}
	}
	return 1
}

func getInput(app *Application) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your input:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	stringArray := strings.Split(input, " ")
	fmt.Println(stringArray)
}

func Menu(app *Application) {
	// Welcome message
	fmt.Println(welcome_message)
	// Options you have
	fmt.Println(menu_option)

	// Read the user input and clear it
	getInput(app)

	fmt.Println(app.userInput)

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