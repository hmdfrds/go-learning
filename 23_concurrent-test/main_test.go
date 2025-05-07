package main

import (
	"sync"
	"testing"
)

func TestWorkerProcessesJobsAndExits(t *testing.T) {
	numTestJobs := 3
	jobs := make(chan int, numTestJobs)
	results := make(chan string, numTestJobs)

	var wg sync.WaitGroup

	wg.Add(1)
	go worker(1, jobs, results, &wg)

	jobs <- 10
	jobs <- 20
	jobs <- 30
	close(jobs)

	wg.Wait()

	close(results)

	receivedResults := make([]string, 0, numTestJobs)

	for result := range results {
		receivedResults = append(receivedResults, result)
	}

	if len(receivedResults) != numTestJobs {
		t.Fatalf("Expected %d results, but got %d", numTestJobs, len(receivedResults))
	}
}
