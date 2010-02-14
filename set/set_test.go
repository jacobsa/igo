// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package set

import (
	"container/vector"
	"reflect"
	"sort"
	"testing"
)

func getSorted(set *StringSet) []string {
	var sorted vector.StringVector
	for val := range set.Iter() {
		sorted.Push(val)
	}

	sort.SortStrings(sorted)
	return sorted.Data()
}

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

	sorted := getSorted(&set)
	expected := []string{}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, sorted)
	}
}

func TestMultipleElements(t *testing.T) {
	var set StringSet

	set.Insert("foo")

	expectDoesntContain(t, &set, "")
	expectContains(t, &set, "foo")
	expectDoesntContain(t, &set, "bar")

	sorted := getSorted(&set)
	expected := []string{"foo"}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, sorted)
	}

	set.Insert("")

	expectContains(t, &set, "")
	expectContains(t, &set, "foo")
	expectDoesntContain(t, &set, "bar")

	sorted = getSorted(&set)
	expected = []string{"", "foo"}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, sorted)
	}
}

func TestRepeatedInserts(t *testing.T) {
	var set StringSet

	set.Insert("")
	set.Insert("")

	expectContains(t, &set, "")

	sorted := getSorted(&set)
	expected := []string{""}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, sorted)
	}

	set.Insert("foo")
	set.Insert("foo")
	set.Insert("foo")

	expectContains(t, &set, "")
	expectContains(t, &set, "foo")

	sorted = getSorted(&set)
	expected = []string{"", "foo"}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, sorted)
	}
}
