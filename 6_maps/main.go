package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	var line string

	fmt.Print("Enter a line of tet: ")

	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error while reading input: ", err)
		os.Exit(1)
	}

	wordMap := make(map[string]int)
	for _, word := range strings.Fields(line) {
		wordMap[strings.ToLower(word)] += 1
	}

	fmt.Println("Word Counts:")

	for k, v := range wordMap {
		fmt.Printf("\"%s\":   %v\n", k, v)
	}

}
