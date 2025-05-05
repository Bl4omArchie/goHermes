package core


import (
    _ "fmt"
    "sync"
)

type DownloadTask struct {
    url string
    filepath string
    taskId int
}

type DownloadResult struct {
    status int
    hash string
    taskId int
}

type DownloadPool struct {
    tasks chan DownloadTask
    results chan DownloadResult
    wg *sync.WaitGroup
}

// This worker accept a document url as a task and return the hash of the downloaded document 
func DownloadWorker(tasks <-chan DownloadTask, results chan<- DownloadResult, errChannel *ErrorChannel) {
    for task := range tasks {
		hashResult, _ := DownloadDocumentReturnHash(task.url, task.filepath, errChannel)
        if hashResult == "" {
            results <- DownloadResult{status: 0, hash: "", taskId: task.taskId}
        } else {
            results <- DownloadResult{status: 1, hash: hashResult, taskId: task.taskId}
        }
    }
}

func StartDownloadPool(numWorkers int, errChannel *ErrorChannel) *DownloadPool {
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
