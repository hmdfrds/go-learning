package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidArgument = errors.New("invalid argument")

type OperationError struct {
	Op  string
	Err error
}

func (e *OperationError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("operation %s failed", e.Op)
	}
	return fmt.Sprintf("operation %s failed: %v", e.Op, e.Err)
}

func (e *OperationError) Unwrap() error {
	return e.Err
}

func divide(a, b float64) (float64, error) {
	if b == 0.0 {
		return 0.0, fmt.Errorf("division by zero")
	}
	if a < 0 || b < 0 {
		return 0.0, &OperationError{Op: "divide", Err: fmt.Errorf("input validation: %w", ErrInvalidArgument)}
	}
	return a / b, nil
}

func main() {

	testCases := []struct{ a, b float64 }{
		{10.0, 2.0},
		{8.0, 0.0},
		{-4.0, 2.0},
		{10.0, -5.0},
	}

	for _, testCase := range testCases {
		result, err := divide(testCase.a, testCase.b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			var opErr *OperationError
			if errors.Is(err, ErrInvalidArgument) {
				fmt.Println("Cause: Invalid arguments provided.")
			}
			if errors.As(err, &opErr) {
				fmt.Printf("Operation Failed: %s\n", opErr.Op)
			}
			if err.Error() == "division by zero" {
				fmt.Print("Cause: Division by zero.\n")
			}

		} else {
			fmt.Printf("Result: %.2f\n", result)
		}
		fmt.Print("---\n\n")

	}

}
