package main

import (
	"fmt"
	"github.com/Bl4omArchie/eprint-DB/core"
)

func get_header() {
	fmt.Println("     ______     _       _    ______ ")
	fmt.Println("     | ___ \\   (_)     | |   |  _  \\")
	fmt.Println("  ___| |_/ / __ _ _ __ | |_  | | | |")
	fmt.Println(" / _ \\  __/ '__| | '_ \\| __| | | | |")
	fmt.Println("|  __/ |  | |  | | | | | |_  | |/ / ")
	fmt.Println(" \\___\\_|  |_|  |_|_| |_|\\__| |___(_)")
	fmt.Println("Archie - 2025")
	fmt.Println("-----------------------------------------\n")
}

func get_menu() {
	fmt.Println("1. Download papers")
	fmt.Println("2. Readme")
	fmt.Println("3. Quit")
}

func get_user_input() int {
	var choice int
	fmt.Print("Enter your choice: ")
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return get_user_input()
	}
	return choice
}

func main() {
	get_header()
	get_menu()
	choice := get_user_input()

	switch choice {
		case 1:
			core.GetDocsPerYears([]int{2024, 2025}, "pdf")
		case 2:
			fmt.Println("Readme selected.")
		case 3:
			fmt.Println("Quitting...")
		default:
			fmt.Println("Invalid choice.")
	}
}
