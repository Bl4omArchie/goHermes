package hermes

import (
	"os"
	"fmt"
	"path/filepath"

)

type HalSource struct {
	Name           string
	Path           string
	BaseUrl        string
	Endpoint	   string
	TotalDocuments int
	Documents      []*Document
}

func NewHalSource() *HalSource {
	return &HalSource{
		Name:      "HAL",
		Path:      "pdf/hal/",
		BaseUrl:   "https://hal.science/search/index/?q=%2A&rows=30&level0_domain_s=info",
		Endpoint:  "https://hal.science/search/index/?q=*&rows=30&level0_domain_s=info&page=800",
		Documents: make([]*Document, 0),
	}
}

func (f *HalSource) Init(engine *Engine) error {
	if err := os.MkdirAll(filepath.Dir(f.Path), os.ModePerm); err != nil {
		CreateLogReport(fmt.Sprintf("Error while creating directories for %s: %v", f.Path, err), engine.Log)
		return err
	}

	return nil
}

func (f *HalSource) Fetch(engine *Engine) error {
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

