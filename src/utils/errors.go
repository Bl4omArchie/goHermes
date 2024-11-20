package utils


import "log"



// check error and quit
func CheckErrorQuit(err error) {
    if err != nil {
        log.Fatalf("\033[31mError: %s\033[0m", err)
    }
}


// check error and raise flag
func CheckAlertError(err error, flag int, msg string, ac *AlertChannel) {
    if err != nil {
        SendAlert(flag, msg, ac)
    }
}