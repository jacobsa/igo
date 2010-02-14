// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package set

import (
	"testing"
)

func expectContains(t *testing.T, set *StringSet, s string) {
	if !set.Contains(s) {
		t.Errorf("Expected set %v to contain string %s", set, s)
	}
}

func expectDoesntContain(t *testing.T, set *StringSet, s string) {
	if set.Contains(s) {
		t.Errorf("Expected set %v to not contain string %s", set, s)
	}
}

func TestEmptySet(t *testing.T) {
	var set StringSet

	expectDoesntContain(t, &set, "")
	expectDoesntContain(t, &set, "foo")
	expectDoesntContain(t, &set, "bar")
}

func TestMultipleElements(t *testing.T) {
	var set StringSet

	set.Insert("")

	expectContains(t, &set, "")
	expectDoesntContain(t, &set, "foo")
	expectDoesntContain(t, &set, "bar")

	set.Insert("foo")

	expectContains(t, &set, "")
	expectContains(t, &set, "foo")
	expectDoesntContain(t, &set, "bar")
}

func TestRepeatedInserts(t *testing.T) {
	var set StringSet

	set.Insert("")
	set.Insert("")

	expectContains(t, &set, "")

	set.Insert("foo")
	set.Insert("foo")
	set.Insert("foo")

	expectContains(t, &set, "")
	expectContains(t, &set, "foo")
}
