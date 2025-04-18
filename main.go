package main

import (
	"fmt"
	"github.com/Bl4omArchie/eprint-DB/core"
)

func get_header() {
	fmt.Println("\033[1;34m     ______     _       _    ______ \033[0m")
	fmt.Println("\033[1;34m     | ___ \\   (_)     | |   |  _  \\\033[0m")
	fmt.Println("\033[1;34m  ___| |_/ / __ _ _ __ | |_  | | | |\033[0m")
	fmt.Println("\033[1;34m / _ \\  __/ '__| | '_ \\| __| | | | |\033[0m")
	fmt.Println("\033[1;34m|  __/ |  | |  | | | | | |_  | |/ / \033[0m")
	fmt.Println("\033[1;34m \\___\\_|  |_|  |_|_| |_|\\__| |___(_)\033[0m")
	fmt.Println("\033[1;32mArchie - 2025 | ePrint-DB | Cryptographic Papers\033[0m")
	fmt.Println("\033[1;33m------------------------------------------------\033[0m\n")
}

func get_menu() {
	fmt.Println("\033[1;33mMenu Options:\033[0m")
	fmt.Println("\033[1;32m1.\033[0m \033[1;34mDownload papers\033[0m")
	fmt.Println("\033[1;32m2.\033[0m \033[1;34mCreate database\033[0m")
	fmt.Println("\033[1;32m3.\033[0m \033[1;34mQuit\033[0m")
	fmt.Println("\033[1;33m------------------------------------------------\033[0m")
}

func get_user_input_int() int {
	var choice int
	fmt.Print("Enter your choice: ")
	_, err := fmt.Scan(&choice)

	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return get_user_input_int()
	}
	fmt.Print("\n")
	return choice
}

func main() {
	get_header()
	get_menu()
	choice := get_user_input_int()

	switch choice {
		case 1:
			core.GetDocsPerYears([]string{"2014", "2015", "2016", "2017", "2018", "2019", "2020", "2021", "2022", "2023", "2024", "2025"}, "pdf")
		case 2:
			core.CreateDB()
		case 3:
			fmt.Println("Exiting program...")
			return
		default:
			fmt.Println("Invalid choice.")
	}
}
