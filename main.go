package main

import (
	"fmt"
	"github.com/Bl4omArchie/eprint-DB/core"
)

func getHeader() {
	fmt.Println("\033[1;34m     ______     _       _    ______ \033[0m")
	fmt.Println("\033[1;34m     | ___ \\   (_)     | |   |  _  \\\033[0m")
	fmt.Println("\033[1;34m  ___| |_/ / __ _ _ __ | |_  | | | |\033[0m")
	fmt.Println("\033[1;34m / _ \\  __/ '__| | '_ \\| __| | | | |\033[0m")
	fmt.Println("\033[1;34m|  __/ |  | |  | | | | | |_  | |/ / \033[0m")
	fmt.Println("\033[1;34m \\___\\_|  |_|  |_|_| |_|\\__| |___(_)\033[0m")
	fmt.Println("\033[1;32mArchie - 2025 | ePrint-DB | Cryptographic Papers\033[0m")
	fmt.Println("\033[1;33m------------------------------------------------\033[0m\n")
}

func getMenu() {
	fmt.Println("\033[1;33mMenu Options:\033[0m")
	fmt.Println("\033[1;32m1.\033[0m \033[1;34mDownload papers\033[0m")
	fmt.Println("\033[1;32m2.\033[0m \033[1;34mCreate database\033[0m")
	fmt.Println("\033[1;32m3.\033[0m \033[1;34mRead a specific PDF\033[0m")
	fmt.Println("\033[1;32m4.\033[0m \033[1;34mQuit\033[0m")
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
			core.StartEngine()
		case 2:
			fmt.Println("Not implemented...")
		case 3:
			fmt.Println("Not implemented...")
		case 4:
			fmt.Println("Exiting program...")
		default:
			fmt.Println("Invalid option")
	}
}
