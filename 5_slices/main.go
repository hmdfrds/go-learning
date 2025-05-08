package main

import "fmt"

func main() {
	var numSlice []int

	for {
		var inputNum int
		fmt.Print("Enter a positive number (or 0/negative to stop): ")
		if _, err := fmt.Scanln(&inputNum); err != nil {
			fmt.Println("Error while reading input: ", err)
			continue
		}

		if inputNum <= 0 {
			break
		}

		numSlice = append(numSlice, inputNum)
	}

	sliceLength := len(numSlice)
	if sliceLength == 0 {
		fmt.Println("No positive numbers entered.")
		return
	}

	var total int
	for _, num := range numSlice {
		total += num
	}

	average := float64(total) / float64(sliceLength)

	fmt.Printf("Average: %.2f", average)
}
