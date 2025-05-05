package main

import (
	"fmt"
	"time"
)

func printMessage(label string, count int) {

	for i := 0; i < count; i++ {
		fmt.Printf("%s: %d\n", label, i)
	}
}

func main() {
	go printMessage("Goroutine A", 5)
	go printMessage("Goroutine B", 5)
	fmt.Println("Main function started goroutines.")
	time.Sleep(1 * time.Second) // educational purpose
	fmt.Println("Main function finished.")
}
