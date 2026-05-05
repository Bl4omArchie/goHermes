package hermes

import (
	"context"

	"github.com/Bl4omArchie/goHermes/models"
)


type Source interface {
	FetchDocumentsUrls(ctx context.Context, network *HermesNetwork) ([]string, error)
	FetchDocumentMetadata(ctx context.Context, network *HermesNetwork) ([]*models.Document, error)
}


func FetchSource(source Source, workers int, network *HermesNetwork, db *HermesDatabase) {
	pipeline := NewDownloadPipeline(network, workers)

	defer pipeline.Close()

	ctx := context.Background()
	pipeline.Start(ctx)

	go func() {
		urls, err := source.FetchDocumentsUrls(ctx, network)
		if err != nil {
			for _, url := range urls {
				pipeline.In() <- url
			}
			close(pipeline.In())
		}
	}()

	for res := range pipeline.Out() {
		if res.IsDownloaded {
			db.InsertTable(&res.Document)
		}
		// else : store into the invalid url database
	}
}


type DownloadDocument interface {
	GetPages() string
	GetUrls() []models.DistantDocument
}
