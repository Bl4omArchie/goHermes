package core

import (
	"sync"
)


type DownloadPool struct {
    tasks chan *Document
    results chan DownloadResult
    waitgroup *sync.WaitGroup
}

type DownloadResult struct {
    status int
    toIngest*Document
}

func DownloadWorker(tasks <-chan *Document, results chan <- DownloadResult, engine *Engine) {
    for task := range tasks {
        // Add Hash
        if hashResult, err := DownloadDocumentReturnHash(task.Url, task.Filepath, engine.Log); err == nil {
            task.Hash = hashResult
            results <- DownloadResult{status: 1, toIngest: task}
        } else  {
            results <- DownloadResult{status: 0}
        }
    }
}

func StartDownloadPool(numWorkers int, engine *Engine) *DownloadPool {
    tasks := make(chan *Document)
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
