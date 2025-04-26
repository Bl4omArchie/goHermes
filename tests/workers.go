package main
 
import (
    "fmt"
    "time"
)

type Task struct {
	name string
	lastName string
}

type Result struct {
	fullname string
}
 
func worker(id int, tasks chan Task, results chan Result) {
    for task := range tasks {
        fmt.Printf("Worker %d processing task %d\n", id, task)
        time.Sleep(time.Second)
		res := &Result{fullname: task.name + task.lastName}
        results <- *res
    }
}
 
func main() {
    numWorkers := 3
    numTasks := 5
 
    tasks := make(chan Task, numTasks)
    results := make(chan Result, numTasks)
 
    for i := 1; i <= numWorkers; i++ {
        go worker(i, tasks, results)
    }
 
    for j := 1; j <= numTasks; j++ {
        tasks <- Task{name: "Pierre", lastName: "Justin"}
    }
    close(tasks)
 
    for k := 1; k <= numTasks; k++ {
        result := <-results
        fmt.Printf("Result: %s\n", result.fullname)
    }
}