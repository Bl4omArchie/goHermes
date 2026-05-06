package hermes

import (
	"context"
	"github.com/Bl4omArchie/goHermes/models"
)


type NetworkPipeline struct {
	input    chan *DownloadInput
	output   chan *DownloadOutput
	workers  int
	network  *HermesNetwork
}


type DownloadInput struct {
	Url string
	Path string
}

type DownloadOutput struct {
	Document		*models.DownloadedDocument
	IsDownloaded	bool
	Error 			error
}


func NewNetworkPipeline(network *HermesNetwork, workers int) *NetworkPipeline {
	return &NetworkPipeline{
		input:   make(chan *DownloadInput),
		output:  make(chan *DownloadOutput),
		workers: workers,
		network: network,
	}
}


func NewDownloadOutput(doc *models.DownloadedDocument, status bool, err error) *DownloadOutput {
	return &DownloadOutput{
		Document: doc,
		IsDownloaded: status,
		Error: err,
	}
}


func (p *NetworkPipeline) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		go p.worker(ctx)
	}
}


func (p *NetworkPipeline) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case input, ok := <-p.input:
			if !ok {
				return
			}

			hash, err := p.network.DownloadDocumentReturnHash(ctx, input.Url, input.Path)
			p.output <- NewDownloadOutput(models.NewDownloadedDocument(input.Url, hash), err == nil, err)
		}
	}
}

func (p *NetworkPipeline) In() chan<- *DownloadInput {
	return p.input
}

func (p *NetworkPipeline) Out() <-chan *DownloadOutput {
	return p.output
}

func (p *NetworkPipeline) Close() {
	close(p.input)
	close(p.output)
}
