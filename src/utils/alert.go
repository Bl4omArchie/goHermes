package utils

import "log"

type AlertChannel struct {
	channel chan int
	nbAlerts int
}


/* Errors message :
	QuitLister 				 = 0x1
	ErrorDownloadingDocument = 0x1337
	ErrorGetPaperData 		 = 0x1338
	ErrorGetStatistics 		 = 0x1339
	ErrorInsertingDocument   = 0x1340
*/

func CreateAlertChannel() *AlertChannel {
	return &AlertChannel{
		channel: make(chan int),
		nbAlerts: 0,
	}
}

func CloseChannel(ac *AlertChannel) {
	close(ac.channel)
}

func SendAlert(errorFlag int, ac *AlertChannel) {
	ac.channel <- errorFlag
	ac.nbAlerts++
}

func ListenerAlertChannel(ac *AlertChannel) {
    for flag := range ac.channel {
        switch flag {
			case 0x1:
				log.Printf("Quitting flag listener...")
				CloseChannel(ac)
				return

			case 0x1337:
				log.Printf("Error downloading a document")

			case 0x1338:
				log.Printf("Error retrieving papers data")

			case 0x1339:
				log.Printf("Error retrieving statistics")

			case 0x1340:
				log.Printf("Error inserting a document")

			default:
				log.Printf("Incorrect flag")
		}
    }
}