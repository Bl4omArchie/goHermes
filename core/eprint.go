package core

import (
	"os"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"path/filepath"
	"golang.org/x/net/html"
)

var (
	baseURL           = "https://eprint.iacr.org"
	endPointByYear    = "/byyear"
)

type EprintSource struct {
	Path string
	TotalDocuments int
	PapersByYear map[string]int
}

type EprintDoc struct {
	Url string
	Doc Document
}

func NewEprintSource() *EprintSource {
	return &EprintSource{
		Path:           "pdf/eprint/",
		TotalDocuments: 0,
		PapersByYear:   make(map[string]int),
	}
}

func (f *EprintSource) Init(engine *Engine) error  {	
	if err := os.MkdirAll(filepath.Dir(f.Path), os.ModePerm); err != nil {
		CreateLogReport(fmt.Sprintf("Error while creating directories for %s: %v", f.Path, err), engine.Log)
		return err
	}

	body, _ := GetPageContent(baseURL + endPointByYear, engine.Log)

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

	f.PapersByYear = make(map[string]int)
	f.PapersByYear["2024"] = 500
	return nil
}

func (f *EprintSource) Fetch(engine *Engine) error {
	downloadPool := StartDownloadPool(engine.NumWorkersPools, engine)

	go func() {
		for year, papersYears := range f.PapersByYear {
			yearUrl := baseURL + "/" + year + "/"

			for docCount := 1; docCount < papersYears; docCount++ {
				docIdYear := fmt.Sprintf("%03d", docCount)
				doc := EprintDoc{
					Url: yearUrl + docIdYear,
					Doc: Document{
						Source: "Cryptology {ePrint} Archive",
					},
				}
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
			InsertTable(engine, &result.toIngest.Doc)
		}
	}
	return nil
}

func FetchMetadata(docTodo *EprintDoc, log *Log) error {
	doc, err := GetPageContent(docTodo.Url, log)
	if err != nil {
		CreateLogReport(fmt.Sprintf("Failed to get metadata for %s: %v", docTodo.Url, err), log)
		return err
	}

	tkn := html.NewTokenizer(strings.NewReader(doc))

	var license, bibtex, pdfHref string
	foundLicense := false
	foundBibtex := false
	foundPDF := false

	for {
		tt := tkn.Next()
		if tt == html.ErrorToken {
			if !foundBibtex {
				CreateLogReport(fmt.Sprintf("Couldn't find bibtex [%s]. Error: %v", docTodo.Url, tkn.Err()), log)
			}
			break
		}

		token := tkn.Token()

		if tt == html.StartTagToken {
			switch token.Data {
			case "small":
				if !foundLicense && tkn.Next() == html.TextToken {
					license = strings.TrimSpace(tkn.Token().Data)
					foundLicense = true
				}

			case "pre":
				if !foundBibtex && tkn.Next() == html.TextToken {
					bibtex = strings.TrimSpace(tkn.Token().Data)
					foundBibtex = true
				}

			case "a":
				if !foundPDF {
					for _, attr := range token.Attr {
						if attr.Key == "class" && attr.Val == "btn btn-sm btn-outline-dark" {
							for _, a := range token.Attr {
								if a.Key == "href" {
									pdfHref = a.Val
									foundPDF = true
									break
								}
							}
						}
					}
				}
			}
		}

		// Stop if all data has been found
		if foundLicense && foundBibtex && foundPDF {
			break
		}
	}

	docTodo.Doc.Filetype = pdfHref
	docTodo.Doc.Url = baseURL + pdfHref
	docTodo.Doc.Filepath = "pdf/eprint" + pdfHref
	docTodo.Doc.License = license

	fields := parseBibTeXSimple(bibtex)
	docTodo.Doc.Title = fields["title"]
	docTodo.Doc.Release = fields["year"]

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
