package core

import (
	"fmt"
	"os"
	"time"
	"sync/atomic"
)

type LogReport struct {
	Message  string
	Timestamp string
}

type Log struct {
	logfile *os.File
	logChannel chan LogReport
	count atomic.Uint64
}

func CreateLogReport(msg string, logChannel *Log) {
	logChannel.logChannel <- LogReport{
		Message:  msg,
		Timestamp: time.Now().Format(time.RFC850),
	}
	logChannel.count.Add(1)
}

func CreateLogChannel(engine *Engine) {
	engine.Log = &Log{
		logfile: CreateLogFile(),
		logChannel: make(chan LogReport),
		count: atomic.Uint64{},
	}
}

func CreateLogFile() (*os.File) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			fmt.Printf("failed to create logs directory: %w", err)
			return nil
		}
	}

	filePath := fmt.Sprintf("logs/log_%d.log", time.Now().Unix())
	logfile, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return nil
	}
	return logfile
}

func ListenerLogFile(logChannel *Log) {
	for errReport := range logChannel.logChannel {
		logEntry := fmt.Sprintf("%s : %s\n", errReport.Timestamp, errReport.Message)
		_, err := logChannel.logfile.WriteString(logEntry)
		if err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
			return
		}
	}
}
