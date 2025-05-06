package main

import "fmt"

func producer(message chan<- string, count int) {
	for i := 0; i < count; i++ {
		formattedMessage := fmt.Sprintf("Message %d", i+1)
		message <- formattedMessage
		fmt.Printf("Send: %s\n", formattedMessage)
	}
	close(message)
}

func main() {
	messageChannel := make(chan string)

	go producer(messageChannel, 3)

	for msg := range messageChannel {
		fmt.Printf("Received: %s\n", msg)
	}

	fmt.Println("Main: All message received. Producer finished.")
}
