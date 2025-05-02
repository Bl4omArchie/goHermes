package engine


import "github.com/Bl4omArchie/eprint-DB/core/utility"

type Source interface {
	ScopeDefinitionProcess() []string
	CraftUrlProcess() []string
	DocumentAcquisitionProcess() []string
}

// Main function that start the PDF scrapping
func StartEngine() {
	errChannel := utility.CreateErrorChannel()
	go utility.ListenerLogFile(errChannel)

	sources := []Source{CreateEprint(), CreateNist()}

	for source := range sources {
		source.ScopeDefinitionProcess(errChannel)
		source.CraftUrlProcess(errChannel)
		go source.DocumentAcquisitionProcess(errChannel)
	}
}
