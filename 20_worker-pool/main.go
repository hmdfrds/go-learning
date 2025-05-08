package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	numWorkers = 3
	numJobs    = 10
)

func worker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d: Started\n", id)

	for job := range jobs {
		fmt.Printf("Worker %d: Processing job %d\n", id, job)
		time.Sleep(500 * time.Millisecond)
		resultValue := job * job
		resultStr := fmt.Sprintf("Job %d (from Worker %d) -> Result: %d", job, id, resultValue)
		results <- resultStr
	}
	fmt.Printf("Worker %d: Finished and shutting down\n", id)
}

func main() {
	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}
	for j := 0; j < numJobs; j++ {
		jobs <- j
		fmt.Printf("Main: Dispatched job %d to jobs channel\n", j)
	}

	close(jobs)
	fmt.Println("Main: All jobs dispatched. Waiting for workers to complete all tasks...")

	wg.Wait()
	close(results)
	fmt.Println("Main: All workers completed. Collecting results...")

	for result := range results {
		fmt.Printf("Main: Collected -> %s\n", result)
	}

	fmt.Println("Main: All results collected. Program finished.")

}
