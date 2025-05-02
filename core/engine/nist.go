package engine


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
		SourceStoragePath: "pdf/nist/",
	}
}