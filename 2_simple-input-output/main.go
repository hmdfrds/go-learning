package main

import "fmt"

func main() {

	var name string
	fmt.Print("Please enter you name: ")
	fmt.Scanln(&name) // will only take first word
	fmt.Printf("\nHello, %s", name)
}
