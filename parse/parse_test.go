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

////////////////////////////////
// GetPackageName
////////////////////////////////

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


////////////////////////////////
// GetImports
////////////////////////////////

func TestGetImportsEmptyFile(t *testing.T) {
	code := ""
	expected := []string{}

	imports := GetImports(code)
	expectContentsEqual(t, imports, expected)
}

func TestGetImportsMultipleImportStatements(t *testing.T) {
	code := `
		package asdf

		import "./foo/bar"
		import "fmt"
		import "./baz"

		func DoSomething() {
		}
	`
	expected := []string{"./foo/bar", "fmt", "./baz"}
	imports := GetImports(code)
	expectContentsEqual(t, imports, expected)
}

func TestGetImportsImportList(t *testing.T) {
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
	imports := GetImports(code)
	expectContentsEqual(t, imports, expected)
}

func TestGetImportsSyntaxErrorBeforeImports(t *testing.T) {
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
	GetImports(code)
}

func TestGetImportsSyntaxErrorInImports(t *testing.T) {
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
	GetImports(code)
}

func TestGetImportsSyntaxErrorAfterImports(t *testing.T) {
	code := `
		package asdf

		import (
			"fmt"
			"os"
		)

		func DoSomething() {
	`
	expected := []string{"fmt", "os"}
	imports := GetImports(code)
	expectContentsEqual(t, imports, expected)
}


////////////////////////////////
// GetTestFunctions
////////////////////////////////

func TestGetTestFunctionsEmptyFile(t *testing.T) {
	code := ""
	expected := []string{}

	imports := GetTestFunctions(code)
	expectContentsEqual(t, imports, expected)
}

func TestGetTestFunctionsNoTestFunctions(t *testing.T) {
	code := `
		package asdf

		import (
			"fmt"
			"os"
			"testing"
		)

		func DoSomething() {}
	`
	expected := []string{}

	imports := GetTestFunctions(code)
	expectContentsEqual(t, imports, expected)
}

func TestGetTestFunctionsSomeResults(t *testing.T) {
	code := `
		package asdf

		import (
			"fmt"
			"os"
			"testing"
		)

		func TestBlah(t *testing.T) {}
		func DoSomething() {}
		func TestFooBar(t *testing.T) {}
	`
	expected := []string{"TestBlah", "TestFooBar"}

	imports := GetTestFunctions(code)
	expectContentsEqual(t, imports, expected)
}

func TestGetTestFunctionsSyntaxError(t *testing.T) {
	code := `
		package asdf

		import (
			"fmt"
			"os"
			"testing"
		)

		func TestBlah(t *testing.T) {}
		func DoSomething() {}
		func TestFooBar(t *testing.T) {}
		func ljlsdfkj {
	`
	expected := []string{}

	imports := GetTestFunctions(code)
	expectContentsEqual(t, imports, expected)
}
