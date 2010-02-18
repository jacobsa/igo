// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package test

import (
	"igo/set"
	"testing"
)

func createSet(contents []string) *set.StringSet {
	var result set.StringSet
	for _, val := range contents {
		result.Insert(val)
	}

	return &result
}

func expectSourceEqual(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Errorf("Expected:\n---------\n%s\n\nActual:\n---------\n%s", expected, actual)
	}
}

func TestEmptySet(t *testing.T) {
	// If there are no test functions, the test runner shouldn't import the
	// package.
	expected :=
`package main

import "testing"
var tests = []testing.Test {
}

func main() {
	testing.Main(tests)
}
`
	actual := GenerateTestMain("blah", createSet([]string{}))
	expectSourceEqual(t, expected, actual)
}

func TestNonEmptySet(t *testing.T) {
	// If there are no test functions, the test runner shouldn't import the
	// package.
	expected :=
`package main

import "testing"
import "./blah"

var tests = []testing.Test {
	blah.TestFoo,
	blah.TestBar,
	blah.TestBaz,
}

func main() {
	testing.Main(tests)
}
`
	funcs := createSet([]string{"TestFoo", "TestBar", "TestBaz"})
	actual := GenerateTestMain("blah", funcs)
	expectSourceEqual(t, expected, actual)
}
