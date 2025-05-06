package core


import (
    "sync"
)

type DownloadPool struct {
    tasks chan EprintDoc
    wg *sync.WaitGroup
}

// This worker accept a document url as a task and return the hash of the downloaded document 
func DownloadWorker(tasks <-chan EprintDoc, logChannel *LogChannel) {
    for task := range tasks {
        if err := GetMetadataEprint(&task, logChannel); err == nil {
            hashResult, _ := DownloadDocumentReturnHash(task.UrlDownload, task.Filepath, logChannel)
            if hashResult != "" {
                task.Hash = hashResult
            }
        }
    }
}

func StartDownloadPool(numWorkers int, logChannel *LogChannel) *DownloadPool {
    tasks := make(chan EprintDoc)
    wg := &sync.WaitGroup{}
    
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            DownloadWorker(tasks, logChannel)
        }()
    }

	return &DownloadPool {
		tasks:   tasks,
        wg: wg,
	}
}
