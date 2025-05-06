package core

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
    "sync/atomic"
)

var (
	baseURL           = "https://eprint.iacr.org"
	endPointComplete  = "/complete"
	endPointByYear    = "/byyear"
)

type EprintSource struct {
	TotalDocuments int
	SourceStoragePath string
	PapersByYear map[string]int
}

type EprintDoc struct {
	UrlMetadata string
	UrlDownload string
	DocId atomic.Uint64
	Filepath string
	Title string
	Hash string
	Release string
	License string
	Authors []Author
}


func InitEprint(errChannel *ErrorChannel) (*EprintSource) {
	eprint := &EprintSource{
		TotalDocuments: 0,
		SourceStoragePath: "pdf/eprint2/",
		PapersByYear: make(map[string]int),
	}
	body, _ := GetPageContent(baseURL + endPointByYear, errChannel)

	re_years := regexp.MustCompile(`>(\d{4})</a> \((\d+) papers\)`)
	matches_years := re_years.FindAllStringSubmatch(body, -1)

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

func DownloadEprint(eprint *EprintSource, errChannel *ErrorChannel) {
	downloadPool := StartDownloadPool(100, errChannel)

	for year, papersYears := range eprint.PapersByYear {
		go func() {
			yearUrl := baseURL + "/" + year + "/"
			
			for docCount := 1; docCount < papersYears; docCount++ {
				docIdYear := fmt.Sprintf("%03d", docCount)

				doc := EprintDoc {
					UrlMetadata: yearUrl + docIdYear,
					UrlDownload: yearUrl + docIdYear +".pdf",
					DocId: atomic.Uint64{},
					Filepath: path.Join(eprint.SourceStoragePath, year, docIdYear+".pdf"),
				}

				downloadPool.tasks <- doc
			}
		}()
	}

    downloadPool.wg.Wait()
    close(downloadPool.tasks)
}

func GetMetadataEprint(docTodo *EprintDoc, errChannel *ErrorChannel) (error) {
	data, err := GetPageContent(docTodo.UrlMetadata, errChannel)

	if (err != nil) {
		CreateErrorReport(fmt.Sprintf("Failed to get metadata for %s: %v", docTodo.UrlMetadata, err), errChannel)
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
