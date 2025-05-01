package engine


import "github.com/Bl4omArchie/eprint-DB/core/utility"

type Source interface {
	ArgsAcquisitionProcess() []string
	CraftUrlProcess() []string
	DocumentAcquisitionProcess() []string
}

// Main function that start the PDF scrapping
func StartEngine() {
	errChannel := utility.CreateErrorChannel()
	go utility.ListenerLogFile(errChannel)

	/*
	for source := range sources {
		source.ArgsAcquisitionProcess(errChannel)
		source.CraftUrlProcess(errChannel)
		go source.DocumentAcquisitionProcess(errChannel)
	}
	*/
}
