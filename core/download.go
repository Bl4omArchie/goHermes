package core


type DownloadPool struct {
    tasks chan EprintDoc
    results chan DownloadResult
}

type DownloadResult struct {
    status int
    toIngest EprintDoc
}

func DownloadWorker(tasks <-chan EprintDoc, results chan <- DownloadResult, logChannel *Log) {
    for task := range tasks {
        if err := GetMetadataEprint(&task, logChannel); err == nil {
            // Add Hash
            if hashResult, err := DownloadDocumentReturnHash(task.Doc.Url, task.Doc.Filepath, logChannel); err == nil {
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

func StartDownloadPool(numWorkers int, logChannel *Log) *DownloadPool {
    tasks := make(chan EprintDoc)
    results := make(chan DownloadResult)
    
    for i := 1; i <= numWorkers; i++ {
        go DownloadWorker(tasks, results, logChannel)
    }

	return &DownloadPool {
		tasks:   tasks,
        results: results,
	}
}
