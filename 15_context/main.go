package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, workerID int, stepDuration time.Duration, totalSteps int, resultChan chan<- string) {
	fmt.Printf("Worker %d: Starting (will attempt %d steps, %v per step).\n", workerID, totalSteps, stepDuration)
	for i := 0; i < totalSteps; i++ {
		select {
		case <-ctx.Done():
			errReason := ctx.Err()
			statusMsg := fmt.Sprintf("Worker %d: Stopped at step %d/%d. Reason: %v", workerID, i+1, totalSteps, errReason)
			fmt.Println(statusMsg)
			resultChan <- statusMsg
			return
		case <-time.After(stepDuration):
			fmt.Printf("Worker %d: Completed step %d/%d.\n", workerID, i+1, totalSteps)
		}
	}

	successMsg := fmt.Sprintf("Worker %d: Finished all %d. steps successfully.", workerID, totalSteps)
	fmt.Println(successMsg)
	resultChan <- successMsg
}

func main() {
	results := make(chan string, 1)

	fmt.Println("--- Scenario 1: Worker completes before timeout ---")
	parentCtx1 := context.Background()
	ctx1, cancel1 := context.WithTimeout(parentCtx1, 4*time.Second)
	defer cancel1()
	go worker(ctx1, 1, 1*time.Second, 3, results)
	fmt.Printf("Main (Scenario 1): %s\n", <-results)

	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n--- Scenario 2: Worker times out before completion ---")
	parentCtx2 := context.Background()
	ctx2, cancel2 := context.WithTimeout(parentCtx2, 1500*time.Millisecond)
	defer cancel2()
	go worker(ctx2, 2, 1*time.Second, 3, results)
	fmt.Printf("Main (Scenario 2): %s\n", <-results)

	fmt.Println("\nMain: Program finished.")

}
