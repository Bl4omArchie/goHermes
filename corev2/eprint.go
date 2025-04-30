package corev2


import (
	"fmt"
	"strconv"
	"regexp"
)

var (
	baseURL           = "https://eprint.iacr.org"
	endpointComplete  = "complete"
	endpointByYear    = "byyear"
)

type SourceEprint struct {
	Name  string
	SourceLink string
	Years []string
	Urls []EprintDocumentToDownload
	StoragePath string
}

type EprintDocumentToDownload {
	UrlMetadata string
	UrlDownload string
	Filepath string
}

// AAP : acquire years for ePrint source
func (eprint *SourceEprint) ArgsAcquisitionProcess(errChannel *ErrorChannel) (*SourceEprint) {
	eprintSource := CreateEprint()
	eprintSource.Years = []string{"2005"}
	return eprintSource
}

// CUP : craft every requested urls to download documents of one source
func (eprint *SourceEprint) CraftUrlProcess(errChannel *ErrorChannel) {
	body := GetPageContent(endpointByYear, errChannel)

	// Seek for years and the number of papers per year
	re_years := regexp.MustCompile(`>(\d{4})</a> \((\d+) papers\)`)
	matches_years := re_years.FindAllStringSubmatch(string(body), -1)

	urlsResult := map[int]string {}
	for _, match := range matches_years {
		if len(match) == 3{
			urlsResult[strconv.Atoi(match[1])] = match[2]
		}
	}

	for _, year := range eprint.Years {
		_, exists := urlsResult[year]
		if !(exists) {
			delete(urlsResult, year)
		}
	}
	eprint.Urls = urlsResult
}

// DAP : receive every urls from CUP and fill data structure
func (eprint *SourceEprint) DocumentAcquisitionProcess(errChannel ErrorChannel) {
	StartDownloadPool(1000, 15, errChannel)
	
	for url := range eprint.Urls {
		metadata := GetMetadataEprint()
		go DownloadDocumentReturnHash(&DownloadResult{url, "pdf/eprint/212.pdf"})
	} 
}

func GetMetadataEprint() {
	reTitle := regexp.MustCompile(`<title>(.*?)</title>`)
	reAuthor := regexp.MustCompile(`<meta name="author" content="(.*?)">`)
	reLicense := regexp.MustCompile(`<meta name="license" content="(.*?)">`)

	title := ""
	authors := []Author{}
	license := ""

	matchTitle := reTitle.FindStringSubmatch(data)
	if len(matchTitle) > 1 {
		title = matchTitle[1]
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
		license = matchLicense[1]
	}
}