package utils

import "log"

type AlertChannel struct {
	channel chan ErrorReport
	nbAlerts int
	errorReportList []ErrorReport
}

type ErrorReport struct {
	errorFlag int
	customMsg string
	action string
}

/* Errors flags :
	QuitLister 				 = 0x1
	ErrorDownloadingDocument = 0x1337
	ErrorGetPaperData 		 = 0x1338
	ErrorGetStatistics 		 = 0x1339
	ErrorInsertingDocument   = 0x1340
	ErrorConnection			 = 0x1341
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

func SendAlert(errorFlag int, message string, action string, ac *AlertChannel) {
	var er = ErrorReport{
		errorFlag: errorFlag,
		customMsg: message,
		action: action,
	}
	ac.errorReportList = append(ac.errorReportList, er)
	ac.channel <- er
	ac.nbAlerts++
}

func ListenerAlertChannel(ac *AlertChannel) {
    for er := range ac.channel {
		switch er.action {
		case "continue":
			log.Printf("\033[31m[LOG]: %s\033[0m", er.customMsg)

		case "quit":
			log.Fatalf("\033[31m[LOG]: %s\033[0m", er.customMsg)
		
		default:
			log.Printf("\033[31m[LOG]: Invalid action\033[0m")
		}
    }
}