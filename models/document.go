package models



type DistantDocument struct {
	Url 	string
	Path	string
	Alive 	bool
}


type DistantDocumentMetadata struct {
	Url			string
	Title		string
	Authors 	[]Author
	Filetype	string
	Release		string
	License 	string		
	Source   	string
}


type DownloadedDocument struct {
	Path string
	Hash string
}


type Document struct {
	Source 		DistantDocument
	Metadata	DistantDocumentMetadata
	Download	DownloadedDocument
}



func NewDistantDocument(url, path string, alive bool) DistantDocument {
	return DistantDocument{
		Url:   url,
		Path: path,
		Alive: alive,
	}
}

func NewDistantDocumentMetadata(
	url string,
	title string,
	authors []Author,
	filetype string,
	release string,
	license string,
	source string,
) *DistantDocumentMetadata {

	return &DistantDocumentMetadata{
		Url:      url,
		Title:    title,
		Authors:  authors,
		Filetype: filetype,
		Release:  release,
		License:  license,
		Source:   source,
	}
}


func NewDownloadedDocument(path string, hash string) *DownloadedDocument {
	return &DownloadedDocument{
		Path: path,
		Hash: hash,
	}
}


func NewDocument(
	source DistantDocument,
	metadata DistantDocumentMetadata,
	download DownloadedDocument,
)* Document {

	return& Document{
		Source:    source,
		Metadata:  metadata,
		Download:  download,
	}
}
