// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package deps

import (
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

func TestEmptyMap(t *testing.T) {
	input := make(map[string][]string)
	result := BuildTotalOrder(input)

	if !reflect.DeepEqual(result, []string{}) {
		t.Errorf("Expected empty, got: %v", result)
	}
}

func TestSinglePackageNoDeps(t *testing.T) {
	input := make(map[string][]string)
	input["foo"] = []string{}

	expected := []string{"foo"}
	result := BuildTotalOrder(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, result)
	}
}

func TestSinglePackageOnlyForeignDeps(t *testing.T) {
	input := make(map[string][]string)
	input["foo"] = []string{"http", "os", "fmt"}

	expected := []string{"foo"}
	result := BuildTotalOrder(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, result)
	}
}
