package engine


import "github.com/Bl4omArchie/eprint-DB/core/utility"


type Source interface {
	ScopeDefinitionProcess(errChan *utility.ErrorChannel)
	CraftUrlProcess(errChan *utility.ErrorChannel)
	DocumentAcquisitionProcess(errChan *utility.ErrorChannel)
}

// Main function that start the PDF scrapping
func StartEngine() {
	errChannel := utility.CreateErrorChannel()
	go utility.ListenerLogFile(errChannel)

	sources := []Source{CreateEprint()}

	for _, source := range sources {
		source.ScopeDefinitionProcess(errChannel)
		source.CraftUrlProcess(errChannel)
		source.DocumentAcquisitionProcess(errChannel)
	}
}
