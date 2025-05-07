package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var atomicCounter int64 = 0

func atomicIncrementWorker(workerID, iterations int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < iterations; i++ {
		atomic.AddInt64(&atomicCounter, 1)
	}
}

func main() {
	numGoroutines := 100
	iterationsPerGoroutine := 1000

	expectedCount := int64(numGoroutines * iterationsPerGoroutine)

	fmt.Println("--- Incrementing Counter with Atomic Operations ---")

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go atomicIncrementWorker(i, iterationsPerGoroutine, &wg)
	}
	wg.Wait()
	finalValue := atomic.LoadInt64(&atomicCounter)

	fmt.Printf("Atomic Operations - Final Counter: %d (Expected: %d)\n", finalValue, expectedCount)

	fmt.Println("Try running this scenario with 'go run -race main.go' to verify no race conditions.")
	fmt.Println("Main: Program finished.")
}
