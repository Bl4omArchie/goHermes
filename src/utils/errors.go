package utils


import "log"


// print error
func CheckError(err error) {
    if err != nil {
        log.Printf("Error:", err)
    }
}

// print error and quit the program
func CheckErrorQuit(err error) {
    if err != nil {
        log.Fatalf("Error:", err)
    }
}