package baz

import "./foo"

func DoStuff() {
	if foo.ComputeTwoPlusTwo() != 4 {
		panic("uh oh")
	}
}
