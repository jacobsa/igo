// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package parse

import (
	"reflect"
	"testing"
)

func TestGetPackageNameEmptyFile(t *testing.T) {
	code := ""
	expected := ""

	packageName := GetPackageName(code)
	if packageName != expected {
		t.Errorf("Expected %s, got: %s", expected, packageName)
	}
}

func TestGetPackageNameGoodFile(t *testing.T) {
	code := `
		package asdf

		import "./foo/bar"
		import "fmt"
		import "./baz"

		func DoSomething() {
		}
	`
	expected := "asdf"

	packageName := GetPackageName(code)
	if packageName != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, packageName)
	}
}

func TestGetPackageNameSyntaxError(t *testing.T) {
	code := `
		package asdf ljksdhk ksd

		import "./foo/bar"
		import "fmt"
		import "./baz"

		func DoSomething() {
		}
	`
	expected := ""

	packageName := GetPackageName(code)
	if packageName != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, packageName)
	}
}

func TestExtractImportsEmptyFile(t *testing.T) {
	code := ""
	expected := []string{}

	imports := ExtractImports(code)
	if !reflect.DeepEqual(imports, expected) {
		t.Errorf("Expected empty deps, got: %v", imports)
	}
}

func TestExtractImportsMultipleImportStatements(t *testing.T) {
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

func TestExtractImportsImportList(t *testing.T) {
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

func TestExtractImportsSyntaxErrorBeforeImports(t *testing.T) {
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

func TestExtractImportsSyntaxErrorInImports(t *testing.T) {
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

func TestExtractImportsSyntaxErrorAfterImports(t *testing.T) {
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
