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

func TestUnion(t *testing.T) {
	var set1 StringSet
	set1.Insert("foo")
	set1.Insert("baz")

	var set2 StringSet
	set2.Insert("foo")
	set2.Insert("bar")

	set1.Union(&set2)

	sorted1 := getSorted(&set1)
	expected1 := []string{"bar", "baz", "foo"}
	if !reflect.DeepEqual(sorted1, expected1) {
		t.Errorf("Expected: %v\nGot: %v", expected1, sorted1)
	}

	sorted2 := getSorted(&set2)
	expected2 := []string{"bar", "foo"}
	if !reflect.DeepEqual(sorted2, expected2) {
		t.Errorf("Expected: %v\nGot: %v", expected2, sorted2)
	}
}
