package bar

import "./bar/baz"
import "./foo"

type AnotherType struct {
	someVar foo.SomeType
}

func RunDoStuff() { baz.DoStuff() }
