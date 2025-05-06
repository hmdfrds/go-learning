package main

import (
	"fmt"
	"time"
)

func fastProducer(message chan<- string, numMessages int) {
	for i := 0; i < numMessages; i++ {
		msg := fmt.Sprintf("Message %d", i+1)
		fmt.Printf("Producer: Preparing to send '%s'.\n", msg)
		message <- msg
		fmt.Printf("Producer: Successfully send '%s'.\n", msg)
	}
	close(message)
}

func main() {
	const bufferCapacity = 3
	const messageToSend = 5

	messageChannel := make(chan string, bufferCapacity)

	go fastProducer(messageChannel, messageToSend)

	time.Sleep(500 * time.Millisecond)

	fmt.Println("Consumer: Waking up to receive messages.")

	for msg := range messageChannel {
		fmt.Printf("Consumer: Received '%s'.\n", msg)
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Consumer: All messages received. Channel closed.")

}
