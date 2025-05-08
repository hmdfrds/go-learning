package main

import (
	"fmt"
)

func main() {

	var n1 int
	var n2 int
	fmt.Print("Enter the first number: ")

	if _, err := fmt.Scanln(&n1); err != nil {
		fmt.Println("Error while reading input: ", err)
		return
	}

	fmt.Print("Enter the second number: ")
	if _, err := fmt.Scanln(&n2); err != nil {
		fmt.Println("Error while reading input: ", err)
		return
	}

	sum := n1 + n2
	fmt.Println("The sum is: ", sum)
	if sum > 10 {
		fmt.Println("The sum is greater than 10.")
	}

}
