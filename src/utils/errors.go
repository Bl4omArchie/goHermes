package utils


import "log"


//print specific error
func CheckErrorCustom(msg string) {
    log.Printf("\033[31mError: ", msg, "\033[0m")
}

// print error
func CheckError(err error) {
    if err != nil {
        log.Printf("\033[31mError: ", err, "\033[0m")
    }
}

// print error and quit the program
func CheckErrorQuit(err error) {
    if err != nil {
        log.Fatalf("\033[31mError: ", err, "\033[0m")
    }
}

// raise a flag
func RaiseFlag(err []byte) int {
    if err == nil {
        return log.Flags()
    }
    return 1
}