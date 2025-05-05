package core

import "fmt"


func StartEngine() {
	errChannel := CreateErrorChannel()
	go ListenerLogFile(errChannel)

	eprint := InitEprint(errChannel)
	fmt.Println(eprint.PapersByYear)
	DownloadEprint(eprint, errChannel)
}