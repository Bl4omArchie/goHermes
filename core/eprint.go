package core

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var (
	baseURL           = "https://eprint.iacr.org"
	endPointComplete  = "/complete"
	endPointByYear    = "/byyear"
)

type EprintSource struct {
	Docs []*EprintDoc
	TotalDocuments int
	SourceStoragePath string
	PapersByYear map[string]int
}

type EprintDoc struct {
	UrlMetadata string
	UrlDownload string
	DocId int
	Filepath string
	Title string
	Hash string
	Release string
	License string
	Authors []Author
}


// SDP : for a given source, target the required papers in a Scope.
func InitEprint(errChannel *ErrorChannel) (*EprintSource) {
	eprint := &EprintSource{
		TotalDocuments: 0,
		Docs: []*EprintDoc{},
		SourceStoragePath: "pdf/eprint/",
		PapersByYear: make(map[string]int),
	}
	body, _ := GetPageContent(path.Join(baseURL, endPointByYear), errChannel)

	re_years := regexp.MustCompile(`>(\d{4})</a> \((\d+) papers\)`)
	matches_years := re_years.FindAllStringSubmatch(string(body), -1)

	sum := 0
	// Fill the struct with years
	for _, match := range matches_years {
		if len(match) == 3 {
			docCount, _ := strconv.Atoi(match[2])
			eprint.PapersByYear[match[1]] = docCount
			sum += docCount
		}
	}
	eprint.TotalDocuments = sum

	return eprint
}

// CUP : craft every requested urls to download documents of one source
func DownloadEprint(eprint *EprintSource, errChannel *ErrorChannel) {
	downloadPool := StartDownloadPool(15, errChannel)

	go func() {
		docId := 0
		for year, papersYears := range eprint.PapersByYear {
			yearUrl := path.Join(baseURL, year)
			fmt.Println(year)
			
			for docCount := range papersYears {
				docIdYear := fmt.Sprintf("%03d", docCount)

				doc := &EprintDoc {
					UrlMetadata: path.Join(yearUrl, docIdYear),
					UrlDownload: path.Join(yearUrl, docIdYear+".pdf"),
					Filepath: path.Join("pdf", "eprint", docIdYear+".pdf"),
				}
				fmt.Println(doc.UrlDownload)

				if (GetMetadataEprint(doc, errChannel) != nil) {
					docId++
					doc.DocId = docId
					downloadPool.tasks <- DownloadTask{doc.UrlDownload, doc.Filepath, doc.DocId}
					eprint.Docs = append(eprint.Docs, doc)
				}
			}
		}
	}()

	for result := range downloadPool.results {
		fmt.Println("Download status:", result.status, "Hash:", result.hash)
		eprint.Docs[result.taskId].Hash = result.hash
	}
}

func GetMetadataEprint(docTodo *EprintDoc, errChannel *ErrorChannel) (error) {
	data, err := GetPageContent(docTodo.UrlMetadata, errChannel)

	if (err != nil) {
		return err
	}
	
	reTitle := regexp.MustCompile(`<title>(.*?)</title>`)
	reAuthor := regexp.MustCompile(`<meta name="author" content="(.*?)">`)
	reLicense := regexp.MustCompile(`<meta name="license" content="(.*?)">`)

	matchTitle := reTitle.FindStringSubmatch(data)
	if len(matchTitle) > 1 {
		docTodo.Title = matchTitle[1]
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
			docTodo.Authors = append(docTodo.Authors, Author{
				FirstName: firstName,
				LastName:  lastName,
			})
		}
	}

	matchLicense := reLicense.FindStringSubmatch(data)
	if len(matchLicense) > 1 {
		docTodo.License = matchLicense[1]
	}
	return nil
}
