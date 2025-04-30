package corev2

type Document struct {
	Title string
	Authors []Author
	Filepath string
	Hash string
	Release string
	License string
	Source string
}

type Author struct {
	FirstName string
	LastName string
}

func CreateDocument(title string, authors []Author, filepath string, hash string, release string, license string, source string) (*Document) {
	return &Document{
		Title:    title,
		Authors:  authors,
		Filepath: filepath,
		Hash:     hash,
		Release:  release,
		License:  license,
		Source:    source,
	}
}

func CreateAuthor(firstName string, lastName string) (*Author) {
	return &Author {
		FirstName: firstName,
		LastName: lastName,
	}
}