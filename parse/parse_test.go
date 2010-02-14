// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package parse

import (
	"container/vector"
	"igo/set"
	"reflect"
	"sort"
	"testing"
)

func expectContentsEqual(t *testing.T, set *set.StringSet, expected []string) {
	var contents vector.StringVector
	for val := range set.Iter() {
		contents.Push(val)
	}

	sort.SortStrings(contents)
	sort.SortStrings(expected)

	if !reflect.DeepEqual(contents.Data(), expected) {
		t.Errorf("Expected:%v\nGot: %v", expected, contents)
	}
}

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
	expectContentsEqual(t, imports, expected)
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
	expectContentsEqual(t, imports, expected)
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
	expectContentsEqual(t, imports, expected)
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
	expectContentsEqual(t, imports, expected)
}
