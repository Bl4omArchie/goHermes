package core

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"golang.org/x/net/html"
)

type FreeHavenSource struct {
	Name string
	Path string
	BaseUrl string
	TotalDocuments int
	Documents []*Document
}


func NewFreeHavenSource() *FreeHavenSource {
	return &FreeHavenSource{
		Name: "FreeHaven",
		Path: "pdf/freeHaven/",
		BaseUrl: "https://www.freehaven.net/anonbib",
		Documents: make([]*Document, 0),
	}
}

func (f *FreeHavenSource) Init(engine *Engine) error {
	if err := os.MkdirAll(filepath.Dir(f.Path), os.ModePerm); err != nil {
		CreateLogReport(fmt.Sprintf("Error while creating directories for %s: %v", f.Path, err), engine.Log)
		return err
	}

	body, _ := GetPageContent(f.BaseUrl, engine.Log)

	root, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			var title string
			var url string
			var filetype string
			var foundTitle bool

			var collect func(*html.Node)
			collect = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "a" {
					for _, attr := range n.Attr {
						if attr.Key == "name" && n.FirstChild != nil {
							title = strings.TrimSpace(n.FirstChild.Data)
							foundTitle = true
						}
						if attr.Key == "href" && (strings.HasSuffix(attr.Val, ".pdf") || strings.HasSuffix(attr.Val, ".ps")) {
							url = attr.Val
							if strings.HasSuffix(attr.Val, ".pdf") {
								filetype = "pdf"
							} else {
								filetype = "ps"
							}
						}
					}
				}
				if n.Type == html.TextNode && !foundTitle {
					s := strings.TrimSpace(n.Data)
					if title == "" && s != "" {
						title = s
					}
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					collect(c)
				}
			}

			collect(n)

			if title != "" && url != "" {
				fullUrl := ""
				if !strings.HasPrefix(url, "http") {
					fullUrl = strings.TrimRight(f.BaseUrl, "/") + "/" + strings.TrimLeft(url, "./")
				} else {
					fullUrl = url
				}

				doc := &Document{
					Title:    title,
					Filetype: filetype,
					Url:      fullUrl,
					Filepath: f.Path + title + "." + filetype,
					Hash:     "",
					Release:  "",
					Source:   f.Name,
				}
				f.Documents = append(f.Documents, doc)
				f.TotalDocuments++
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)

	return nil
}

func (f *FreeHavenSource) Fetch(engine *Engine) error {
	downloadPool := StartDownloadPool(engine.NumWorkersPools, engine)

	go func() {
		for _, doc := range f.Documents {
			downloadPool.tasks <- doc
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
