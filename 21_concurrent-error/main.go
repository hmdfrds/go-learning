package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func performTask(ctx context.Context, taskID int, shouldFail bool) error {
	fmt.Printf("Task %d: Starting\n", taskID)
	workDuration := 500*time.Millisecond + time.Duration(taskID%3)*300*time.Millisecond

	select {
	case <-ctx.Done():
		fmt.Printf("Task %d: Cancelled by context (err: %v)\n", taskID, ctx.Err())
		return ctx.Err()
	case <-time.After(workDuration):
		if shouldFail {
			fmt.Printf("Task %d: Work completed, but FAILING as instructed.\n", taskID)
			return fmt.Errorf("task %d failed intentionally", taskID)
		}
		fmt.Printf("Task %d: Work completed successfully.\n", taskID)
		return nil
	}
}

func main() {
	taskConfigs := []struct {
		id         int
		shouldFail bool
	}{
		{id: 1, shouldFail: false},
		{id: 2, shouldFail: true},
		{id: 3, shouldFail: false},
		{id: 4, shouldFail: false},
		{id: 5, shouldFail: false},
	}

	parentCtx := context.Background()
	g, gCtx := errgroup.WithContext(parentCtx)
	g.SetLimit(2)

	for _, taskConfig := range taskConfigs {
		cfg := taskConfig
		g.Go(func() error { return performTask(gCtx, cfg.id, cfg.shouldFail) })
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("\nMain: Operation failed with error: %v\n", err)
		fmt.Println("Main: Other tasks may have been cancelled due to this error.")
	} else {
		fmt.Println("\nMain: All tasks completed successfully (or were cancelled without an error from themselves).")
	}

	fmt.Println("Main: Program finished.")
}
