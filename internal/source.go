package hermes

import (
	"context"

	"github.com/Bl4omArchie/goHermes/models"
)


type Source interface {
	nextPage() string
	ExtractUrls() []models.DistantDocument
}



func FetchSource(source Source, workers int, network *HermesNetwork, db *HermesDatabase) {
	pipeline := NewDownloadPipeline(network, workers)

	defer pipeline.Close()

	ctx := context.Background()
	pipeline.Start(ctx)

	go func() {
		urls, err := source.ExtractUrls(ctx, network)
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
