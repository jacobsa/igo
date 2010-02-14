// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package parse

import (
	"reflect"
	"testing"
)

func TestEmptyFile(t *testing.T) {
	code := ""
	expected := []string{}

	imports := ExtractImports(code)
	if !reflect.DeepEqual(imports, expected) {
		t.Errorf("Expected empty deps, got: %v", imports)
	}
}

func TestMultipleImportStatements(t *testing.T) {
	code := `
		package asdf

		import "./foo/bar"
		import "fmt"
		import "./baz"

		func DoSomething() {
		}
	`
	expected := []string{"./foo/bar", "fmt", "./baz"}
	imports := ExtractImports(code)
	if !reflect.DeepEqual(imports, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, imports)
	}
}

func TestImportList(t *testing.T) {
	code := `
		package asdf

		import (
			"./foo/bar"
			"fmt"
			"./baz"
		)

		func DoSomething() {
		}
	`
	expected := []string{"./foo/bar", "fmt", "./baz"}
	imports := ExtractImports(code)
	if !reflect.DeepEqual(imports, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, imports)
	}
}

func TestSyntaxErrorBeforeImports(t *testing.T) {
	code := `
		package asdf
		kjdfgkjdlshfjghdfkg

		import "./foo/bar"
		import "fmt"
		import "./baz"

		func DoSomething() {
		}
	`

	// Shouldn't crash.
	ExtractImports(code)
}

func TestSyntaxErrorInImports(t *testing.T) {
	code := `
		package asdf

		import (
			"fmt"
			dkfjlgjkf
			"os"
		)

		func DoSomething() {
		}
	`

	// Shouldn't crash.
	ExtractImports(code)
}

func TestSyntaxErrorAfterImports(t *testing.T) {
	code := `
		package asdf

		import (
			"fmt"
			"os"
		)

		func DoSomething() {
	`
	expected := []string{"fmt", "os"}
	imports := ExtractImports(code)
	if !reflect.DeepEqual(imports, expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, imports)
	}
}
