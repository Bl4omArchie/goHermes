package core


import "sync/atomic"

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
	UrlMetadata string
	UrlDownload string
	PubDate string
	Title string
	Hash string
	Summary string
	Authors []Author
	DocId atomic.Uint64
	
}