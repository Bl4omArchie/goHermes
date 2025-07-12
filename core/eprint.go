package core

import (
	"os"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
	"path/filepath"
	"golang.org/x/net/html"
)

type EprintSource struct {
	Name string
	Path string
	BaseUrl string
	Endpoint string
	TotalDocuments int
	PapersByYear map[string]int
	Documents []*Document
}


func NewEprintSource() *EprintSource {
	return &EprintSource{
		Name: "Cryptology {ePrint} Archive",
		Path: "pdf/eprint",
		BaseUrl: "https://eprint.iacr.org",
		Endpoint: "/byyear",
		TotalDocuments: 0,
		PapersByYear: make(map[string]int),
		Documents: make([]*Document, 0),
	}
}

func (f *EprintSource) Init(engine *Engine) error  {	
	if err := os.MkdirAll(filepath.Dir(f.Path), os.ModePerm); err != nil {
		CreateLogReport(fmt.Sprintf("Error while creating directories for %s: %v", f.Path, err), engine.Log)
		return err
	}

	body, _ := GetPageContent(f.BaseUrl + f.Endpoint, engine.Log)

	re_years := regexp.MustCompile(`>(\d{4})</a> \((\d+) papers\)`)
	matches_years := re_years.FindAllStringSubmatch(body, -1)

	sum := 0
	// Fill the struct with years
	for _, match := range matches_years {
		if len(match) == 3 {
			docCount, _ := strconv.Atoi(match[2])
			f.PapersByYear[match[1]] = docCount
			sum += docCount
		}
	}
	f.TotalDocuments = sum

	return nil
}

func (f *EprintSource) Fetch(engine *Engine) error {
	downloadPool := StartDownloadPool(engine.NumWorkersPools, engine)

	go func() {
		for year, papersYears := range f.PapersByYear {
			yearUrl := f.BaseUrl + "/" + year + "/"

			for docCount := 1; docCount < papersYears; docCount++ {
				docIdYear := fmt.Sprintf("%03d", docCount)
				doc := &Document{
					Url: yearUrl + docIdYear,
					Source: f.Name,
				}
				errFetchMeta := FetchMetadata(doc, engine)
				if errFetchMeta != nil{
					continue
				}
				f.Documents = append(f.Documents, doc)
				downloadPool.tasks <- doc
			}
		}
		close(downloadPool.tasks)
	}()

	go func() {
		downloadPool.waitgroup.Wait()
		close(downloadPool.results)
	}()

	for result := range downloadPool.results {
		if result.status == 1 {
			InsertTable(engine, &result.toIngest)
		}
	}
	return nil
}

func FetchMetadata(doc *Document, engine *Engine) error {
	body, err := GetPageContent(doc.Url, engine.Log)
	if err != nil {
		CreateLogReport(fmt.Sprintf("Failed to get metadata for %s: %v", doc.Url, err), engine.Log)
		return err
	}

	tkn := html.NewTokenizer(strings.NewReader(body))

	var license, bibtex, pdfUrl, title, year string
	foundLicense := false
	foundBibtex := false
	foundPDF := false
	foundTitle := false
	foundYear := false
	foundBibtexPre := false

	for {
		tt := tkn.Next()
		if tt == html.ErrorToken {
			if !foundBibtex {
				CreateLogReport(fmt.Sprintf("Couldn't find bibtex [%s]. Error: %v", doc.Url, tkn.Err()), engine.Log)
			}
			break
		}

		token := tkn.Token()

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			switch token.Data {
			case "meta":
				var name, content string
				for _, attr := range token.Attr {
					if attr.Key == "name" {
						name = attr.Val
					}
					if attr.Key == "content" {
						content = attr.Val
					}
				}
				switch name {
				case "citation_title":
					title = content
					foundTitle = true
				case "citation_publication_date":
					year = content
					foundYear = true
				case "citation_pdf_url":
					if !foundPDF {
						pdfUrl = content
						foundPDF = true
					}
				}

			case "a":
				var classVal, hrefVal, relVal string
				for _, attr := range token.Attr {
					if attr.Key == "class" {
						classVal = attr.Val
					}
					if attr.Key == "href" {
						hrefVal = attr.Val
					}
					if attr.Key == "rel" {
						relVal = attr.Val
					}
				}
				// Fallback PDF link if meta not found
				if !foundPDF && classVal == "btn btn-sm btn-outline-dark" && strings.HasSuffix(hrefVal, ".pdf") {
					pdfUrl = hrefVal
					foundPDF = true
				}
				if !foundLicense && relVal == "license" {
					license = hrefVal
					foundLicense = true
				}

			case "pre":
				for _, attr := range token.Attr {
					if !foundBibtexPre && attr.Key == "id" && attr.Val == "bibtex" {
						foundBibtexPre = true
						tt = tkn.Next()
						if tt == html.TextToken {
							bibtex = strings.TrimSpace(tkn.Token().Data)
							foundBibtex = true
						}
					}
				}
			}
		}

		if foundLicense && foundBibtex && foundPDF && foundTitle && foundYear {
			break
		}
	}

	fields := parseBibTeXSimple(bibtex)
	if title == "" {
		title = fields["title"]
	}
	if year == "" {
		year = fields["year"]
	}

	doc.Title = title
	doc.Release = year
	doc.License = license

	ext := path.Ext(pdfUrl)
	if ext == "" {
		ext = ".pdf"
	}
	doc.Filetype = ext
	doc.Url = strings.TrimRight(doc.Url, "/") + doc.Filetype

	sanitizedTitle := strings.ReplaceAll(title, "/", "-")
	sanitizedTitle = strings.ReplaceAll(sanitizedTitle, ":", "-")
	sanitizedTitle = strings.ReplaceAll(sanitizedTitle, "\\", "-")
	sanitizedTitle = strings.TrimSpace(sanitizedTitle)

	doc.Filepath = fmt.Sprintf("pdf/eprint/%s/%s%s", year, sanitizedTitle, ext)

	return nil
}

func parseBibTeXSimple(input string) map[string]string {
	fields := make(map[string]string)

	input = strings.Trim(input, "{}")
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		val = strings.TrimSuffix(val, ",")
		val = strings.Trim(val, "{}")

		fields[key] = val
	}

	return fields
}
