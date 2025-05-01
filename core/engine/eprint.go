package engine


import (
	"fmt"
	"strings"
	"strconv"
	"regexp"
	"github.com/Bl4omArchie/eprint-DB/core/utility"
)

var (
	baseURL           = "https://eprint.iacr.org"
	endpointComplete  = "complete"
	endpointByYear    = "byyear"
)

type SourceEprint struct {
	Name  string
	Years map[string]int
	Urls []EprintDocumentToDownload
	StoragePath string
}

type EprintDocumentToDownload struct {
	UrlMetadata string
	UrlDownload string
	Filepath string
}

// AAP : acquire years for ePrint source
func (eprint *SourceEprint) ArgsAcquisitionProcess(errChannel *utility.ErrorChannel) {
	eprint.Years = map[string]int{"2025": 500}
}

// CUP : craft every requested urls to download documents of one source
func (eprint *SourceEprint) CraftUrlProcess(errChannel *utility.ErrorChannel) {
	body := utility.GetPageContent(endpointByYear, errChannel)

	// Seek for years and the number of papers per year
	re_years := regexp.MustCompile(`>(\d{4})</a> \((\d+) papers\)`)
	matches_years := re_years.FindAllStringSubmatch(string(body), -1)

	urlsResult := map[string]int {}
	for _, match := range matches_years {
		if len(match) == 3{
			stat, _:= strconv.Atoi(match[2])
			urlsResult[match[1]] = stat
		}
	}

	fmt.Println(urlsResult)
}

// DAP : receive every urls from CUP and fill data structure
func (eprint *SourceEprint) DocumentAcquisitionProcess(errChannel *utility.ErrorChannel) {
	receive_channel := StartDownloadPool(1000, 15, errChannel)
	
	for url := range eprint.Urls {
		GetMetadataEprint(url.UrlMetadata, errChannel)
		receive_channel <- *DownloadTask{url, "pdf/eprint/212.pdf"}
	}
}

func GetMetadataEprint(url string, errChannel *utility.ErrorChannel) {
	data := utility.GetPageContent(url, errChannel)
	
	reTitle := regexp.MustCompile(`<title>(.*?)</title>`)
	reAuthor := regexp.MustCompile(`<meta name="author" content="(.*?)">`)
	reLicense := regexp.MustCompile(`<meta name="license" content="(.*?)">`)

	authors := []Author{}

	matchTitle := reTitle.FindStringSubmatch(data)
	if len(matchTitle) > 1 {
		title := matchTitle[1]
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
			authors = append(authors, Author{
				FirstName: firstName,
				LastName:  lastName,
			})
		}
	}

	matchLicense := reLicense.FindStringSubmatch(data)
	if len(matchLicense) > 1 {
		license := matchLicense[1]
	}
}

func CreateEprint() (*SourceEprint) {
	return &SourceEprint {
		Name: "eprint",
		StoragePath: "pdf/eprint/",
	}
}