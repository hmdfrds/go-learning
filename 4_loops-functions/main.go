package main

import "fmt"

func add(n1, n2 int) int {
	return n1 + n2
}

func main() {

	totalSum := 0
	inputNumber := 0
	for {
		fmt.Print("Enter a number (or 0 to stop): ")

		// If there's error the \n will still stuck at the input buffer. So it will accept \n as next input. Not good. Use bufio instead.
		if _, err := fmt.Scanln(&inputNumber); err != nil {
			fmt.Println("Error while reading input: ", err)
			continue
		}
		if inputNumber == 0 {
			break
		}
		totalSum = add(totalSum, inputNumber)
	}

	fmt.Println("The final sum is: ", totalSum)

}
