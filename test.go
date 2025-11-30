package main

import (
	"fmt"
)

type CarInterface interface {
	Drive()
}

type Car struct {
	Brand string
}

func (c *Car) Drive() {
	fmt.Println("Driving", c.Brand)
}

func test(car CarInterface) {
	car.Drive()
}

func main() {
	fmt.Println("Hello, World!")
	car := Car{Brand: "Toyota"}
	test(&car)
}
