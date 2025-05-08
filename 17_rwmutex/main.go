package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	sharedData = make(map[string]int)
	rwMutex    sync.RWMutex
)

func dataReader(readerID int, keyToRead string, wg *sync.WaitGroup) {
	defer wg.Done()

	rwMutex.RLock()
	defer rwMutex.RUnlock()

	value, found := sharedData[keyToRead]

	fmt.Printf("Reader %d: Reading key '%s'... Value: %d (Found: %t)\n", readerID, keyToRead, value, found)
}

func dataWriter(writerID int, keyToWrite string, valueToWrite int, wg *sync.WaitGroup) {
	defer wg.Done()

	rwMutex.Lock()
	defer rwMutex.Unlock()

	sharedData[keyToWrite] = valueToWrite

	fmt.Printf("Writer %d: Wrote key '%s' with value %d.\n", writerID, keyToWrite, valueToWrite)
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go dataWriter(1, "configItem1", 100, &wg)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go dataReader(i, "configItem1", &wg)
		time.Sleep(20 * time.Millisecond)
	}

	wg.Add(1)
	go dataWriter(2, "configItem2", 200, &wg)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go dataReader(6, "configItem1", &wg)
		wg.Add(1)
		go dataReader(7, "configItem2", &wg)
		wg.Add(1)
		go dataReader(8, "nonExistentItem", &wg)
	}

	wg.Add(1)
	go dataWriter(3, "configItem1", 101, &wg)

	wg.Wait()

	rwMutex.RLock()
	fmt.Println("\nMain: Final content of sharedData:")
	for key, value := range sharedData {
		fmt.Printf("	'%s': %d\n", key, value)
	}
	rwMutex.RUnlock()

	fmt.Println("Main: All operations completed.")
}
