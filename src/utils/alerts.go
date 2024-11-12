package utils

import (
	"fmt"
	"time"
)


type AlertChannel struct {
	channel chan ErrorReport
	nbAlerts int
	errorReportList []ErrorReport
}

type ErrorReport struct {
	flag int			// See doc below
	customMsg string	// a custom error message
	date time.Time
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

| **Flag type**           | **Action quit program** | **Action continue program** |
|-------------------------|-------------------------|-----------------------------|
| ExitListener		      | 0x81      		   		| 0xc1                        |
| ErrorDownloadingDocument| 0x82      		   		| 0xc2                        |
| ErrorGetPaperData       | 0x83      		   		| 0xc3                        |
| ErrorGetStatistics      | 0x84      		   		| 0xc4                 		  |
| ErrorInsertingDocument  | 0x85      		   		| 0xc5                   	  |
| ErrorConnection         | 0x86      		   		| 0xc6                 	  	  |

Example : I'm downloading a bunch of PDF and one of the url is incorrect (PDF not found). However, I still need to download the remainding PDF.
In this situation my flag is : 0xc2 (0b10100000 in binary).
*/

func CreateAlertChannel() *AlertChannel {
	return &AlertChannel{
		channel: make(chan ErrorReport),
		nbAlerts: 0,
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
    for er := range ac.channel {
		action := er.flag >> 6

		switch action {
		case 0b10:
			fmt.Printf("Currently, this action is not supported")

		case 0b11:
			fmt.Printf("\033[33m[%s/%s/%s] [%s:%s:%s] \033[31m[LOG] Flag raised: %s\033[0m", er.date.Day(), er.date.Month(), er.date.Year(), er.date.Hour(), er.date.Minute(), er.date.Second(), er.customMsg)
		
		default:
			fmt.Printf("\033[31m[LOG]: Invalid action\033[0m")
		}
    }
}