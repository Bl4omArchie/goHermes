package utility

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type ErrorReport struct {
	Message  string
	Timestamp string
}

type ErrorChannel struct {
	logfile *os.File
	wec   chan ErrorReport
	count int
	mu    sync.Mutex
}

func CreateErrorReport(msg string, wec *ErrorChannel) {
	wec.wec <- ErrorReport{
		Message:  msg,
		Timestamp: time.Now().Format(time.RFC850),
	}
	wec.mu.Lock()
	wec.count++
	wec.mu.Unlock()
}

func CreateErrorChannel() *ErrorChannel {
	return &ErrorChannel{
		logfile: CreateLogFile(),
		wec:   make(chan ErrorReport),
		count: 0,
	}
}

func CreateLogFile() (*os.File) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", os.ModeDir)
		if err != nil {
			fmt.Printf("Error creating logs directory: %v\n", err)
			return nil
		}
	}

	filePath := fmt.Sprintf("logs/error_%d.log", time.Now().Unix())
	logfile, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return nil
	}
	return logfile
}

func ListenerLogFile(errChannel *ErrorChannel) {
	for errReport := range errChannel.wec {
		logEntry := fmt.Sprintf("%s : %s\n", errReport.Timestamp, errReport.Message)
		_, err := errChannel.logfile.WriteString(logEntry)
		if err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
			return
		}
	}
}
