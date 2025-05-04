package main

import "fmt"

type Rectangle struct {
	Width  float64
	Height float64
}

func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Height + r.Width)
}

func main() {

	rectangle := Rectangle{}

	fmt.Print("Enter rectangle width: ")
	if _, err := fmt.Scanln(&rectangle.Width); err != nil {
		fmt.Println("Error while reading width: ", err)
		return
	}

	fmt.Print("Enter rectangle height: ")
	if _, err := fmt.Scanln(&rectangle.Height); err != nil {
		fmt.Println("Error while reading height: ", err)
		return
	}

	fmt.Printf("Area: %10.2f\n", rectangle.Area())

	fmt.Printf("Perimeter: %.2f\n", rectangle.Perimeter())

}
