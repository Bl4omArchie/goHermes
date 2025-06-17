package core

import (
	"sync"
)


type DownloadPool struct {
    tasks chan EprintDoc
    results chan DownloadResult
    waitgroup *sync.WaitGroup
}

type DownloadResult struct {
    status int
    toIngest EprintDoc
}

func DownloadWorker(tasks <-chan EprintDoc, results chan <- DownloadResult, engine *Engine) {
    for task := range tasks {
        if err := FetchMetadata(&task, engine.Log); err == nil {
            // Add Hash
            if hashResult, err := DownloadDocumentReturnHash(task.Doc.Url, task.Doc.Filepath, engine.Log); err == nil {
                task.Doc.Hash = hashResult
                results <- DownloadResult{status: 1, toIngest: task}
            } else  {
                results <- DownloadResult{status: 0}
            }
        } else  {
            results <- DownloadResult{status: 0}
        }
    }
}

func StartDownloadPool(numWorkers int, engine *Engine) *DownloadPool {
    tasks := make(chan EprintDoc)
    results := make(chan DownloadResult)
    var wg sync.WaitGroup

    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go func() {
            DownloadWorker(tasks, results, engine)
            wg.Done()
        }()
    }

    return &DownloadPool{
        tasks:   tasks,
        results: results,
        waitgroup: &wg,
    }
}
