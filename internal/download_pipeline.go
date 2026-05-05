package hermes

import (
	"context"
	"github.com/Bl4omArchie/goHermes/models"
)


type NetworkAction int
const (
	DonwloadAndHash NetworkAction = iota
	Download
	Request
)


type NetworkPipeline struct {
	input    chan *NetworkInput
	output   chan *DownloadOutput
	workers  int
	network  *HermesNetwork
}


type NetworkInput struct {
	Url 	string
	Action  NetworkAction
}


type DownloadOutput struct {
	Document		*models.DownloadedDocument
	IsDownloaded	bool
	Error 			error
}


func NewNetworkPipeline(network *HermesNetwork, workers int) *NetworkPipeline {
	return &NetworkPipeline{
		input:   make(chan *NetworkInput),
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

			switch input.Action {
			case DonwloadAndHash:
				hash, err := p.network.DownloadDocumentReturnHash(ctx, url, doc.Path)
				p.output <- NewDownloadOutput(models.NewDownloadedDocument(doc.Url, hash), err == nil, err)
			
			case Request:
				hash, err := p.network.DownloadDocumentReturnHash(ctx, url, doc.Path)
				p.output <- NewDownloadOutput(models.NewDownloadedDocument(doc.Url, hash), err == nil, err)
			}
		}
	}
}

func (p *NetworkPipeline) In() chan<- string {
	return p.input
}

func (p *NetworkPipeline) Out() <-chan *DownloadOutput {
	return p.output
}

func (p *NetworkPipeline) Close() {
	close(p.input)
	close(p.output)
}
