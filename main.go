package main

import (
	"fmt"
	"github.com/Bl4omArchie/eprint-DB/core"
)

func getHeader() {
	fmt.Println("\033[1;34m              _   _                               \033[0m")
	fmt.Println("\033[1;34m   __ _  ___ | | | | ___ _ __ _ __ ___   ___  ___ \033[0m")
	fmt.Println("\033[1;34m  / _` |/ _ \\| |_| |/ _ \\ '__| '_ ` _ \\ / _ \\/ __|\033[0m")
	fmt.Println("\033[1;34m | (_| | (_) |  _  |  __/ |  | | | | | |  __/\\__ \\\033[0m")
	fmt.Println("\033[1;34m  \\__, |\\___/|_| |_|\\___|_|  |_| |_| |_|\\___||___/\033[0m")
	fmt.Println("\033[1;34m  |___/                                           \033[0m")
	fmt.Println("\033[1;32marchie - 2025 | goHermes | Computer science papers\033[0m")
	fmt.Println("\033[1;33m---------------------------------------------------\033[0m\n")
}


func getMenu() {
	fmt.Println("\033[1;33mMenu Options:\033[0m")
	fmt.Println("\033[1;32m1.\033[0m \033[1;34mDownload papers\033[0m")
	fmt.Println("\033[1;32m2.\033[0m \033[1;34mQuit\033[0m")
	fmt.Println("\033[1;33m------------------------------------------------\033[0m")
}

func GetIntegerInput(message string) int {
	var choice int
	fmt.Print("Enter your choice: ")
	_, err := fmt.Scan(&choice)

	if err != nil {
		fmt.Println(message)
		return GetIntegerInput(message)
	}
	fmt.Print("\n")
	return choice
}

func main() {
	getHeader()
	getMenu()
	choice := GetIntegerInput("Input your option : ")

	switch choice {
		case 1:
			fmt.Println("Launching downloading engine...")
			core.StartEngine()
		case 2:
			fmt.Println("Exiting program...")
		default:
			fmt.Println("Invalid option")
	}
}
