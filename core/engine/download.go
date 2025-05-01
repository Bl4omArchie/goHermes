package engine


import (
    "fmt"
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

// This worker accept a document url as a task and return the hash of the downloaded document 
func DownloadWorker(tasks <-chan DownloadTask, results chan<- DownloadResult, errChannel *utility.ErrorChannel) {
    for task := range tasks {
		hashResult := utility.DownloadDocumentReturnHash(task.url, task.filepath, errChannel)
        if hashResult == "" {
            results <- DownloadResult{status: 0, hash: ""}
        }
        results <- DownloadResult{status: 1, hash: hashResult}
    }
}

func StartDownloadPool(numWorkers int, numTasks int, errChannel *utility.ErrorChannel) (chan DownloadTask) {
    tasks := make(chan DownloadTask, numTasks)
    results := make(chan DownloadResult, numTasks)

    for i := 1; i <= numWorkers; i++ {
        go DownloadWorker(tasks, results, errChannel)
    }
    return tasks
}

func ListenDownloadPool(numTasks int, results <- chan DownloadResult) {
    for k := 1; k <= numTasks; k++ {
        result := <-results
        fmt.Println(result.status)
    }
}
