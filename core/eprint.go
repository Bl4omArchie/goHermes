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
	
	// eprint.PapersByYear = map[string]int{"1997": 14}

	for year, papersYears := range eprint.PapersByYear {
		go func() {
			yearUrl := baseURL + "/" + year + "/"
			
			for docCount := 1; docCount < papersYears; docCount++ {
				docIdYear = fmt.Sprintf("%03d", docCount)

				// Instantiate Document and EprintDoc (add url, filepath and source)
				doc := EprintDoc {
					Url: yearUrl + docIdYear,
					Doc: Document{Url: yearUrl + docIdYear +".pdf", Filepath: path.Join(eprint.SourceStoragePath, year, docIdYear+".pdf"), Source: "Cryptology {ePrint} Archive", },
				}
				downloadPool.tasks <- doc
			}
		}()
	}

	for result := range downloadPool.results {
		if (result.status == 1) {
			/*
			fmt.Println(result.toIngest.Doc.Title)
			fmt.Println(result.toIngest.Doc.Authors)
			fmt.Println(result.toIngest.Doc.Release)
			fmt.Println(result.toIngest.Doc.License)
			fmt.Println("=========================")
			*/
			continue
		}
	}
}


func processNode(n *html.Node) {
	// TODO
    switch n.Data {
		case "h2":
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				name := n.FirstChild.Data
				fmt.Println("Name:", name)
			}
    }

    for c := n.FirstChild; c != nil; c = c.NextSibling {
        processNode(c)
    }
}

func GetMetadataEprint(docTodo *EprintDoc, log *Log) (error) {
	doc, err := GetParsedPageContent(docTodo.Url, log)
    if err != nil {
        CreateLogReport(fmt.Sprintf("Failed to get metadata for %s: %v", docTodo.Url, err), log)
        return err
    }

	var processAllProduct func(*html.Node)
    processAllProduct = func(n *html.Node) {
        if n.Type == html.ElementNode {
            processNode(n)
 
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            processAllProduct(c)
        }
    }
    processAllProduct(doc)

	return nil
}
