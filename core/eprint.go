package core

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
	"golang.org/x/net/html"
)

var (
	baseURL           = "https://eprint.iacr.org"
	endPointByYear    = "/byyear"
)

type EprintSource struct {
	TotalDocuments int
	SourceStoragePath string
	PapersByYear map[string]int
}

type EprintDoc struct {
	Url string
	Doc Document
	Authors []Author
}


func InitEprint(engineInstance *Engine) (*EprintSource) {
	eprint := &EprintSource{
		TotalDocuments: 0,
		SourceStoragePath: "pdf/eprint/",
		PapersByYear: make(map[string]int),
	}
	body, _ := GetPageContent(baseURL + endPointByYear, engineInstance.Log)

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

func DownloadEprint(eprint *EprintSource, engineInstance *Engine) {
	downloadPool := StartDownloadPool(engineInstance.NumWorkersPools, engineInstance.Log)
	docIdYear := ""
	
	eprint.PapersByYear = map[string]int{"1997": 14}

	for year, papersYears := range eprint.PapersByYear {
		go func() {
			yearUrl := baseURL + "/" + year + "/"
			
			for docCount := 1; docCount < papersYears; docCount++ {
				docIdYear = fmt.Sprintf("%03d", docCount)

				// Instantiate Document and EprintDoc (add url, filepath and source)
				doc := EprintDoc {
					Url: yearUrl + docIdYear,
					Doc: Document{
						Url: yearUrl + docIdYear +".pdf", 
						Filepath: path.Join(eprint.SourceStoragePath, year, docIdYear+".pdf"), 
						Source: "Cryptology {ePrint} Archive", },
				}
				downloadPool.tasks <- doc
			}
		}()
	}

	for result := range downloadPool.results {
		if (result.status == 1) {
			fmt.Println("Title: ", result.toIngest.Doc.Title)
			fmt.Println("Filepath: ", result.toIngest.Doc.Filepath)
			fmt.Println("Release: ", result.toIngest.Doc.Release)
			fmt.Println("Hash: ",result.toIngest.Doc.Hash)
			fmt.Println("License: ",result.toIngest.Doc.License)
			fmt.Println("=========================")
			continue
		}
	}
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

func GetMetadataEprint(docTodo *EprintDoc, log *Log) (error) {
	doc, err := GetPageContent(docTodo.Url, log)
    if err != nil {
        CreateLogReport(fmt.Sprintf("Failed to get metadata for %s: %v", docTodo.Url, err), log)
        return err
    }

    tkn := html.NewTokenizer(strings.NewReader(doc))

	license := ""
	bibtex := ""
	MetaFlag := 0
	
	for MetaFlag < 2 {
		tt := tkn.Next()

		switch {
			case tt == html.ErrorToken:
				CreateLogReport(fmt.Sprintf("Couldn't find bibtex [%s]. Error : %v", docTodo.Url, tkn.Err()), log)
				MetaFlag = 2

			case tt == html.StartTagToken:
				t := tkn.Token()
				
				// Get license
				if t.Data == "small" {
					tt = tkn.Next()

					if tt == html.TextToken {
						t := tkn.Token()
						license = t.Data
						MetaFlag++
						break
					}
				}
				
				// Get bibtex (title, release date, authors)
				if t.Data == "pre" {
					tt = tkn.Next()

					if tt == html.TextToken {
						t := tkn.Token()
						bibtex = t.Data
						MetaFlag++
						break
					}
				}
		}
	}

	docTodo.Doc.License = license
	fields := parseBibTeXSimple(bibtex)
	docTodo.Doc.Title = fields["title"]
	docTodo.Doc.Release = fields["year"]

	return nil
}
