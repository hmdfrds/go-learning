package main

import (
	"fmt"
	"sync"
)

var counter = 0
var mutex sync.Mutex

func incrementWorker(id int, iterations int, wg *sync.WaitGroup, useMutex bool) {
	defer wg.Done()
	for i := 0; i < iterations; i++ {
		if useMutex {
			mutex.Lock()
			currentCounter := counter
			currentCounter++
			counter = currentCounter
			mutex.Unlock()
		} else {
			currentCounter := counter
			currentCounter++
			counter = currentCounter
		}
	}
}

func main() {
	const numGoroutines = 100
	const iterationsPerGoroutine = 1000

	expectedCount := numGoroutines * iterationsPerGoroutine

	fmt.Println("--- Scenario 1: Incrementing without Mutex (Potential Race Condition) ---")
	counter = 0
	var wg1 sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg1.Add(1)
		go incrementWorker(i, iterationsPerGoroutine, &wg1, false)
	}

	wg1.Wait()
	fmt.Printf("Without Mutex - Final Counter: %d (Expected: %d)\n", counter, expectedCount)

	fmt.Println("Try running this scenario with 'go run -race main.go' to detect race conditions.")

	fmt.Println("\n--- Scenario 2: Incrementing with Mutex (Race Condition Protected) ---")
	counter = 0
	var wg2 sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg2.Add(1)
		go incrementWorker(i, iterationsPerGoroutine, &wg2, true)
	}

	wg2.Wait()
	fmt.Printf("With Mutex - Final Counter: %d (Expected: %d)\n", counter, expectedCount)

	fmt.Println("\nMain: Program finished.")
}
