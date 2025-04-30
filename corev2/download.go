package corev2


import (
    "fmt"
    "time"
)

type DownloadTask struct {
    url string
    filepath string
}

type DownloadResult struct {
    status int
    hash string
}

// This worker accept a document url as a task and return the hash of the downloaded document 
func DownloadWorker(tasks <-chan DownloadTask, results chan<- DownloadResult, errChannel ErrorChannel) {
    for task := range tasks {
		hashResult := DownloadDocumentReturnHash(task, errChannel)
        results <- hashResult
    }
}

func StartDownloadPool(numWorkers int, numTasks int, errChannel ErrorChannel) {
    tasks := make(chan DownloadTask, numTasks)
    results := make(chan DownloadResult, numTasks)

    for i := 1; i <= numWorkers; i++ {
        go DownloadWorker(tasks, results, errChannel)
    }
    close(tasks)
}

func ListenDownloadPool(numTasks int, results <- chan DownloadResult) {
    for k := 1; k <= numTasks; k++ {
        result := <-results
        fmt.Println(result.status)
    }
}
