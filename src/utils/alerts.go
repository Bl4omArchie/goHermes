package utils

import (
	"fmt"
	"time"
)


type AlertChannel struct {
	active int
	nbAlerts int
	channel chan ErrorReport
	errorReportList []ErrorReport
}

type ErrorReport struct {
	flag int			// See doc below
	customMsg string	// A custom error message
	date time.Time 		// Time when the flag was raised
}

/* 
A flag is a 8 bits integer representing an error type and an action
The first 2 bits (from left) is the action. The last 6 bits the flag type.

Flag type :
	ExitListener 		     = 0b000001
	ErrorDownloadingDocument = 0b000010
	ErrorGetPaperData 		 = 0b000011
	ErrorGetStatistics 		 = 0b000100
	ErrorInsertingDocument   = 0b000101
	ErrorConnection			 = 0b000110

Action :
	ExitProgram 	= 0b10
	ContinueProgram = 0b11

In practise, you shall used the flag as an hexadecimal value. Here is a board that convert everything in hexadecimal :

(See errors.md for flags listing)

Example : I'm downloading a bunch of PDF and one of the url is incorrect (PDF not found). However, I still need to download the remainding PDF.
In this situation my flag is : 0xc2 (0b10100000 in binary).
*/

func CreateAlertChannel() *AlertChannel {
	return &AlertChannel{
		active: 0,
		nbAlerts: 0,
		channel: make(chan ErrorReport),
		errorReportList: []ErrorReport{},
	}
}

func CloseChannel(ac *AlertChannel) {
	close(ac.channel)
}

func SendAlert(flag int, msg string, ac *AlertChannel) {
	var er = ErrorReport{
		flag: flag,
		customMsg: msg,
		date: time.Now(),
	}
	ac.errorReportList = append(ac.errorReportList, er)
	ac.channel <- er
	ac.nbAlerts++
}


func ListenerAlertChannel(ac *AlertChannel) {
	fmt.Println("[!] Welcome on the flag listener ...")
	ac.active = 1

    for er := range ac.channel {
		action := er.flag >> 6

		switch action {
		case 0b10:
			fmt.Println("Currently, this action is not supported")

		case 0b11:
			fmt.Printf("\033[33m[%d/%d/%d] [%d:%d:%d] \033[31m[LOG] Flag raised: %s\033[0m\n", er.date.Day(), er.date.Month(), er.date.Year(), er.date.Hour(), er.date.Minute(), er.date.Second(), er.customMsg)
		
		default:
			fmt.Println("\033[31m[LOG]: Invalid action\033[0m")
		}
    }
}