package main

import (
	"fmt"
	"sync"
	"time"
)

func printMessage(label string, count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		fmt.Printf("%s: %d\n", label, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go printMessage("Go routine A", 5, &wg)

	wg.Add(1)
	go printMessage("Go routine B", 3, &wg)

	fmt.Println("Main: Goroutines launched. Waiting for them to finish...")
	wg.Wait()
	fmt.Println("Main: All goroutines complete. Program finished.")
}
