package core


var (
	arxivApi = "export.arxiv.org/api/query?search_query="
)


type ArxivRssFeed struct {
	Query string
	Filepath string
	Version string
	Encoding string
	ArxivDocs []ArxivDoc
}

type ArxivDoc struct {
	Url string
	Doc Document
	Authors []Author
}

func InitArxiv() {

}

func DownloadArxiv() {

}

func GetMetadataArxiv() {
	
}