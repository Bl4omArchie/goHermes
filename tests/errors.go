package main

import (
	"fmt"
	"sync"
)

type ErrorReport struct {
	message string
	wid int
}

type ErrorChannel struct {
	wec   chan ErrorReport
	count int
	mu    sync.Mutex
}

func CreateErrorChannel() (*ErrorChannel) {
	return &ErrorChannel {
		wec: make(chan ErrorReport),
		count: 0,
	}
}

func CreateErrorReport(message string, wid int, ch *ErrorChannel) {
	ch.wec <- ErrorReport {message: message, wid: wid}
}

func ListenerErrorChannel(ch *ErrorChannel) {
	for report := range(ch.wec) {
		fmt.Println("Error on worker : ", report)
	}
}

func main() {
	ch := CreateErrorChannel()
	go ListenerErrorChannel(ch)
	go CreateErrorReport("yooo", 3, ch)
	go CreateErrorReport("mila", 3, ch)
}