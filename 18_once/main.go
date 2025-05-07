package main

import (
	"fmt"
	"sync"
	"time"
)

type HeavyService struct {
	ID string
}

var (
	serviceInstance *HeavyService
	initOnce        sync.Once
)

func initializeService() {
	fmt.Println("ServiceInitializer: Starting heavy initialization process...")
	time.Sleep(2 * time.Second)
	serviceInstance = &HeavyService{ID: "SvcInst%001-Alpha"}
	fmt.Printf("ServiceInitializer: HeavyService initialized with ID: %s\n", serviceInstance.ID)
}

func getServiceInstance(goroutineID int) *HeavyService {
	fmt.Printf("Goroutine %d: Attempting to get service instance.\n", goroutineID)
	initOnce.Do(initializeService)
	fmt.Printf("Goroutine %d: Service instance obtained.\n", goroutineID)
	return serviceInstance
}

func main() {
	var wg sync.WaitGroup

	numConcurrentRequesters := 5

	for i := 0; i < numConcurrentRequesters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if instance := getServiceInstance(i); instance != nil {
				fmt.Printf("Goroutine %d: Using service. ID: %s\n", i, instance.ID)
			} else {
				fmt.Printf("Goroutine %d: Failed to get service instance.\n", i)
			}
		}()
	}
	wg.Wait()

	if finalInstance := getServiceInstance(99); finalInstance != nil {
		fmt.Printf("Main: Final check - Service ID: %s\n", finalInstance.ID)
	}

	fmt.Println("Main: All operation completed.")
}
