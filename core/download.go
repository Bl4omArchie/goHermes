package core


import (
    "sync"
    "sync/atomic"
)

type DownloadResult struct {
    status int
    hash string
    taskId atomic.Uint64
}

type DownloadPool struct {
    tasks chan EprintDoc
    wg *sync.WaitGroup
}

// This worker accept a document url as a task and return the hash of the downloaded document 
func DownloadWorker(tasks <-chan EprintDoc, errChannel *ErrorChannel) {
    for task := range tasks {
        if err := GetMetadataEprint(&task, errChannel); err == nil {
            task.DocId.Add(1)
        }

		hashResult, _ := DownloadDocumentReturnHash(task.UrlDownload, task.Filepath, errChannel)
        if hashResult != "" {
            task.Hash = hashResult
        }
    }
}

func StartDownloadPool(numWorkers int, errChannel *ErrorChannel) *DownloadPool {
    tasks := make(chan EprintDoc)
    wg := &sync.WaitGroup{}
    
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            DownloadWorker(tasks, errChannel)
        }()
    }

	return &DownloadPool {
		tasks:   tasks,
        wg: wg,
	}
}
