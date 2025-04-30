package corev2


type Source interface {
	ArgsAcquisitionProcess() []string
	CraftUrlProcess() []string
	DocumentAcquisitionProcess() []string
}

type SourceNist struct {
	Name string
	SourceLink string
	Category string
	Urls map[int]string
	StoragePath string
}

func CreateNist() (*SourceEprint) {
	return &SourceEprint {
		Name: "NIST",
		SourceLink: "https://www.nist.gov/publications/search?ta%5B0%5D=248731",
		Years: "",
		Urls: {},
		StoragePath: "pdf/nist/",
	}
}

// Main function that start the PDF scrapping
func StartEngine() {
	errChannel = CreateErrorChannel()
	go ListenerLogFile(errChannel)

	sources := []Source{CreateEprint(), CreateNist()}

	for source := range sources {
		source.ArgsAcquisitionProcess(errChannel)
		source.CraftUrlProcess(errChannel)
		go source.DocumentAcquisitionProcess(errChannel)
	}
}
