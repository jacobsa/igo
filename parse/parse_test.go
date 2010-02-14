// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package parse

import (
	"reflect"
	"testing"
)

func TestEmptyFile(t *testing.T) {
	deps, err := ExtractDependencies("")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expected := make([]string, 0)
	if !reflect.DeepEqual(deps, expected) {
		t.Errorf("Expected empty deps, got: %v", deps)
	}
}
