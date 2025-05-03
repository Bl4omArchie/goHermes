package engine


import (
    _ "fmt"
    "sync"
    "github.com/Bl4omArchie/eprint-DB/core/utility"
)

type DownloadTask struct {
    url string
    filepath string
}

type DownloadResult struct {
    status int
    hash string
}

type DownloadPool struct {
    tasks chan DownloadTask
    results chan DownloadResult
    wg *sync.WaitGroup
}

// This worker accept a document url as a task and return the hash of the downloaded document 
func DownloadWorker(tasks <-chan DownloadTask, results chan<- DownloadResult, errChannel *utility.ErrorChannel) {
    for task := range tasks {
		hashResult, _ := utility.DownloadDocumentReturnHash(task.url, task.filepath, errChannel)
        if hashResult == "" {
            results <- DownloadResult{status: 0, hash: ""}
        } else {
            results <- DownloadResult{status: 1, hash: hashResult}
        }
    }
}

func StartDownloadPool(numWorkers int, errChannel *utility.ErrorChannel) *DownloadPool {
    tasks := make(chan DownloadTask)
    results := make(chan DownloadResult)
    var wg sync.WaitGroup
   
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
		go func() {
			defer wg.Done()
			DownloadWorker(tasks, results, errChannel)
		}()
    }
	go func() {
		wg.Wait()
		close(results)
	}()

	return &DownloadPool{
		tasks:   tasks,
		results: results,
		wg:      &wg,
	}
}
