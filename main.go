package main


import (
	"fmt"
	"time"
	"sync"
	"github.com/Bl4omArchie/ePrint-DB/src/api"
)


func main() {
	var wg1 sync.WaitGroup

	input_list := []string{"1996"}
    wg1.Add(1)

	start := time.Now()
	fmt.Println("Goroutine #1 launched ...")
	go api.DownloadPapers(input_list, &wg1)

    wg1.Wait()
	fmt.Println("Final execution time:", time.Since(start))
}