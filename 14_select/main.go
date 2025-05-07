package main

import (
	"fmt"
	"time"
)

func messageProducer(label string, messages chan<- string, delay time.Duration, numMessages int) {
	for i := 0; i < numMessages; i++ {
		time.Sleep(delay)
		msg := fmt.Sprintf("%s: Message %d", label, i+1)
		messages <- msg
		fmt.Printf("%s send '%s'\n", label, msg)
	}
}

func main() {
	channelA := make(chan string)
	channelB := make(chan string)

	go messageProducer("Producer A", channelA, 800*time.Millisecond, 3)
	go messageProducer("Producer B", channelB, 1200*time.Millisecond, 2)

	for i := 0; i < 5; i++ {
		select {
		case msgA := <-channelA:
			{
				fmt.Printf("Main: Received from A -> '%s'\n", msgA)
			}
		case msgB := <-channelB:
			{
				fmt.Printf("Main: Received from B -> '%s'\n", msgB)
			}
		case <-time.After(1 * time.Second):
			{
				fmt.Println("Main: Timeout - No message received within 1 second.")
			}
		}
	}

	fmt.Println("Main: Finished receiving messages.")

}
