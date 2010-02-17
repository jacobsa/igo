package main

import "./foo"
import "fmt"

func main() {
	fmt.Printf("Result: %d\n", foo.ComputeThreeTimesThree())
}
