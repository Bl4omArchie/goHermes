package engine

import (
	"fmt"
	_ "fmt"
	"regexp"
	"strconv"
	"strings"
	"github.com/Bl4omArchie/eprint-DB/core/database"
	"github.com/Bl4omArchie/eprint-DB/core/utility"
)

var (
	baseURL           = "https://eprint.iacr.org"
	endPointComplete  = "/complete"
	endPointByYear    = "/byyear"
)

type SourceEprint struct {
	Name  string
	SourceStoragePath string
	Docs []EprintDocumentToDownload
	Scope EprintScope
}

type EprintScope struct {
	TotalDocuments int
	PapersByYear map[string]int
	Categories []string
}

type EprintDocumentToDownload struct {
	UrlMetadata string
	UrlDownload string
	Filepath string
	Doc database.Document
	Authors database.Author
}


// SDP : for a given source, target the required papers in a Scope.
func (eprint *SourceEprint) ScopeDefinitionProcess(errChannel *utility.ErrorChannel) {
	body, _ := utility.GetPageContent(baseURL + endPointByYear, errChannel)

	re_years := regexp.MustCompile(`>(\d{4})</a> \((\d+) papers\)`)
	matches_years := re_years.FindAllStringSubmatch(string(body), -1)
	
	// Seek for categories
	re_categories := regexp.MustCompile(`<a href="/search\?category=[^"]+">([^<]+)</a>`)
	matches_categories := re_categories.FindAllStringSubmatch(string(body), -1)

	sum := 0
	// Fill the struct with years
	for _, match := range matches_years {
		if len(match) == 3 {
			docCount, _ := strconv.Atoi(match[2])
			eprint.Scope.PapersByYear[match[1]] = docCount
			sum += docCount
		}
	}
	eprint.Scope.TotalDocuments = sum

	// Fill the struct with categories
	for _, match := range matches_categories {
		if len(match) == 2 {
			eprint.Scope.Categories = append(eprint.Scope.Categories, match[1])
		}
	}
}

// CUP : craft every requested urls to download documents of one source
func (eprint *SourceEprint) CraftUrlProcess(errChannel *utility.ErrorChannel) {
	for year, papersYears := range eprint.Scope.PapersByYear {
		docCountYear := 0
		for i := 0; i < papersYears; i++ {
			docCountYear++
			eprint.Docs = append(eprint.Docs, *CreateDocumentToDownload(baseURL, year, fmt.Sprintf("%03d", docCountYear)))
		}
	}
}

// DAP : receive every urls from CUP and fill data structure
func (eprint *SourceEprint) DocumentAcquisitionProcess(errChannel *utility.ErrorChannel) {
	dp := StartDownloadPool(15, errChannel)
	
	go func() {
		for _, docTodo := range eprint.Docs {
			GetMetadataEprint(&docTodo, errChannel)
			dp.tasks <- DownloadTask{docTodo.UrlDownload, url.Filepath}
		}
		close(dp.tasks)
	}()

	for result := range dp.results {
		fmt.Println("Download status:", result.status, "Hash:", result.hash)
	}
}

func GetMetadataEprint(docTodo *EprintDocumentToDownload, errChannel *utility.ErrorChannel) {
	data, _ := utility.GetPageContent(docTodo.UrlMetadata, errChannel)
	
	reTitle := regexp.MustCompile(`<title>(.*?)</title>`)
	reAuthor := regexp.MustCompile(`<meta name="author" content="(.*?)">`)
	reLicense := regexp.MustCompile(`<meta name="license" content="(.*?)">`)

	matchTitle := reTitle.FindStringSubmatch(data)
	if len(matchTitle) > 1 {
		docTodo.Doc.Title = matchTitle[1]
	}

	matchAuthors := reAuthor.FindAllStringSubmatch(data, -1)
	for _, match := range matchAuthors {
		if len(match) > 1 {
			names := strings.Split(match[1], " ")
			firstName := ""
			lastName := ""
			if len(names) > 1 {
				firstName = names[0]
				lastName = names[len(names)-1]
			} else {
				lastName = names[0]
			}
			docTodo.Doc.Authors = append(docTodo.Doc.Authors, database.Author{
				FirstName: firstName,
				LastName:  lastName,
			})
		}
	}

	matchLicense := reLicense.FindStringSubmatch(data)
	if len(matchLicense) > 1 {
		docTodo.Doc.License = matchLicense[1]
	}
}

func CreateEprint() (*SourceEprint) {
	return &SourceEprint {
		Name: "eprint",
		SourceStoragePath: "pdf/eprint/",
		Docs: []EprintDocumentToDownload{},
		Scope: *CreateEprintScope(),
	}
}

func CreateEprintScope() (*EprintScope) {
	return &EprintScope {
		TotalDocuments: 0,
		PapersByYear: map[string]int{},
		Categories: []string{},
	}
}

func CreateDocumentToDownload(url string, year string, paperId string) (*EprintDocumentToDownload) {
	endpoint := "/" + year + "/" + paperId + ".pdf"
	return &EprintDocumentToDownload{
		UrlMetadata: url,
		UrlDownload: url + endpoint,
		Filepath: "pdf/eprint" +  endpoint,
	}
}