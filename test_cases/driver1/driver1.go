package main

import "./bar"
import "./foo"
import "fmt"

func main() {
	bar.RunDoStuff()
	fmt.Printf("Result: %d\n", foo.ComputeTwoPlusTwo())
}
