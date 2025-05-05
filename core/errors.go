package core

import (
	"fmt"
	"os"
	"time"
	"sync/atomic"
)

type ErrorReport struct {
	Message  string
	Timestamp string
}

type ErrorChannel struct {
	logfile *os.File
	errChannel chan ErrorReport
	count atomic.Uint32
}

func CreateErrorReport(msg string, errChannel *ErrorChannel) {
	errChannel.errChannel <- ErrorReport{
		Message:  msg,
		Timestamp: time.Now().Format(time.RFC850),
	}
	errChannel.count.Add(1)
}

func CreateErrorChannel() *ErrorChannel {
	return &ErrorChannel{
		logfile: CreateLogFile(),
		errChannel:   make(chan ErrorReport),
		count: atomic.Uint32{},
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
	for errReport := range errChannel.errChannel {
		logEntry := fmt.Sprintf("%s : %s\n", errReport.Timestamp, errReport.Message)
		_, err := errChannel.logfile.WriteString(logEntry)
		if err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
			return
		}
	}
}
