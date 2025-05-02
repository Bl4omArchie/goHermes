package utility

import "fmt"


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