package main

import (
	"fmt"
	"os"
)

func getHeader() {
	fmt.Println("\033[1;34m              _   _                               \033[0m")
	fmt.Println("\033[1;34m   __ _  ___ | | | | ___ _ __ _ __ ___   ___  ___ \033[0m")
	fmt.Println("\033[1;34m  / _` |/ _ \\| |_| |/ _ \\ '__| '_ ` _ \\ / _ \\/ __|\033[0m")
	fmt.Println("\033[1;34m | (_| | (_) |  _  |  __/ |  | | | | | |  __/\\__ \\\033[0m")
	fmt.Println("\033[1;34m  \\__, |\\___/|_| |_|\\___|_|  |_| |_| |_|\\___||___/\033[0m")
	fmt.Println("\033[1;34m  |___/                                           \033[0m")
	fmt.Println("\033[1;32marchie - 2026 | goHermes | Computer science papers\033[0m")
	fmt.Println("\033[1;33m---------------------------------------------------\033[0m\n")
}


type HermesEngineConfig struct {
	DatabaseName	string
	NumWorkers		int
	Sources 		[]string
}

func NewHermesEngineConfig(dbName string, workers int, sources []string) *HermesEngineConfig {
	return &HermesEngineConfig{
		DatabaseName: dbName,
		NumWorkers: workers,
		Sources: sources,
	}
}

func CreateHermesEngineConfig(args []string) *HermesEngineConfig {
	var sources []string = make([]string, 0)

	for _, source := range args[3:] {
		sources = append(sources, source)
	}
	return NewHermesEngineConfig(args[1], args[2], sources)
}


func main() {
	getHeader()

	args := os.Args[1:]
	_ = CreateHermesEngineConfig(args)
}
