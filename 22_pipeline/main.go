package main

import (
	"fmt"
	"sync"
	"time"
)

func generator(done <-chan struct{}, startNum, count int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			currentGeneratedNumber := startNum + i
			select {
			case out <- currentGeneratedNumber:
			case <-done:
				return
			}
		}
	}()
	return out
}

func transformer(done <-chan struct{}, factor, offset int, inputStream <-chan int) <-chan int {
	out := make(chan int)

	go func() {

		defer close(out)

		for value := range inputStream {
			transformedValue := (value * factor) + offset
			select {
			case out <- transformedValue:
			case <-done:
				return
			}
		}
	}()
	return out
}

func sink(done <-chan struct{}, inputStream <-chan int, wg *sync.WaitGroup) {

	defer wg.Done()

	for value := range inputStream {
		select {
		case <-done:
			fmt.Println("Sink: Cancellation signal received, shutting down.")
			return
		default:
		}
		fmt.Printf("Sink: Processed value %d\n", value)
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Sink: Finished processing all values.")
}

func main() {
	done := make(chan struct{})
	defer close(done)
	var wg sync.WaitGroup

	startNumber := 1
	numberOfItems := 5
	multiplicationFactor := 2
	additionOffset := 1

	generatedOutout := generator(done, startNumber, numberOfItems)
	transformedOutput := transformer(done, multiplicationFactor, additionOffset, generatedOutout)
	wg.Add(1)
	go sink(done, transformedOutput, &wg)

	wg.Wait()
	fmt.Println("Main: Pipeline processing fully completed.")

}
