package core

import (
	"fmt"
	"os"
	"sync"
)

type WorkerErrorReport struct {
	Message  string
	WorkerID int
}

type WorkerErrorChannel struct {
	wec   chan WorkerErrorReport
	count int
	mu    sync.Mutex
}

func CreateWorkerErrorReport(msg string, wid int, wec *WorkerErrorChannel) {
	wec.wec <- WorkerErrorReport{
		Message:  msg,
		WorkerID: wid,
	}
	wec.mu.Lock()
	wec.count++
	wec.mu.Unlock()
}

func CreateWorkerErrorChannel() *WorkerErrorChannel {
	return &WorkerErrorChannel{
		wec:   make(chan WorkerErrorReport),
		count: 0,
	}
}

func GenerateLog(wec *WorkerErrorChannel, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return
	}
	defer file.Close()

	for errReport := range wec.wec {
		logEntry := fmt.Sprintf("WorkerID: %d, Message: %s\n", errReport.WorkerID, errReport.Message)
		_, err := file.WriteString(logEntry)
		if err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
			return
		}
	}
}