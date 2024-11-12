package utils


import "log"



// check error and quit
func CheckErrorQuit(err error) {
    if err != nil {
        log.Fatalf("\033[31mError: ", err, "\033[0m")
    }
}



// check alert with an error
func CheckAlertError(err error, flag int, msg string, ac *AlertChannel) {
    if err != nil {
        SendAlert(flag, msg, ac)
    }
}