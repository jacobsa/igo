// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package deps

import (
	"igo/set"
	"reflect"
	"testing"
)

// Return the index of the first occurence of needle in haystack, or
// len(haystack) if none is found.
func indexOf(haystack []string, needle string) int {
	for i, val := range haystack {
		if val == needle {
			return i
		}
	}

	return len(haystack)
}

func expectContains(t *testing.T, haystack []string, needle string) {
	if indexOf(haystack, needle) == len(haystack) {
		t.Errorf("Expected %v contains %s", haystack, needle)
	}
}

func expectComesBefore(t *testing.T, array []string, a string, b string) {
	expectContains(t, array, a)
	expectContains(t, array, b)

	if indexOf(array, a) > indexOf(array, b) {
		t.Errorf("Expected %s comes before %s, got: %v", a, b, array)
	}
}

func addDeps(depsMap map[string]*set.StringSet, name string, deps []string) {
	var set set.StringSet
	for _, val := range deps {
		set.Insert(val)
	}

	depsMap[name] = &set
}

func TestEmptyMap(t *testing.T) {
	input := make(map[string]*set.StringSet)
	result := BuildTotalOrder(input)

	if !reflect.DeepEqual(result, []string{}) {
		t.Errorf("Expected empty, got: %v", result)
	}
}

func TestSinglePackageNoDeps(t *testing.T) {
	input := make(map[string]*set.StringSet)
	addDeps(input, "foo", []string{})

	expected := []string{"foo"}
	result := BuildTotalOrder(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, result)
	}
}

func TestSinglePackageOnlyForeignDeps(t *testing.T) {
	input := make(map[string]*set.StringSet)
	addDeps(input, "foo", []string{"http", "os", "fmt"})

	expected := []string{"foo"}
	result := BuildTotalOrder(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, result)
	}
}

func TestTwoPackagesNoOrderingRestraints(t *testing.T) {
	input := make(map[string]*set.StringSet)
	addDeps(input, "foo", []string{"http", "os", "fmt"})
	addDeps(input, "bar", []string{})

	result := BuildTotalOrder(input)
	if len(result) != 2 {
		t.Errorf("Expected two elements, got %v", result)
	}

	expectContains(t, result, "foo")
	expectContains(t, result, "bar")
}

func TestTwoPackagesWithOrderingRestraints(t *testing.T) {
	input := make(map[string]*set.StringSet)
	addDeps(input, "bar", []string{"http", "foo"})
	addDeps(input, "foo", []string{})

	result := BuildTotalOrder(input)
	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot:%v", expected, result)
	}
}

func TestFourPackages(t *testing.T) {
	input := make(map[string]*set.StringSet)
	addDeps(input, "foo", []string{"baz"})
	addDeps(input, "bar", []string{"foo"})
	addDeps(input, "baz", []string{})
	addDeps(input, "tony", []string{"baz"})

	result := BuildTotalOrder(input)
	if len(result) != 4 {
		t.Errorf("Expected four elements, got %v", result)
	}

	expectComesBefore(t, result, "baz", "foo")
	expectComesBefore(t, result, "foo", "bar")
	expectComesBefore(t, result, "baz", "tony")
}

func TestCircularDependency(t *testing.T) {
	input := make(map[string]*set.StringSet)
	addDeps(input, "foo", []string{"bar"})
	addDeps(input, "bar", []string{"foo"})

	// Should give some output.
	result := BuildTotalOrder(input)
	if len(result) != 2 {
		t.Errorf("Expected two elements, got %v", result)
	}

	expectContains(t, result, "foo")
	expectContains(t, result, "bar")
}
