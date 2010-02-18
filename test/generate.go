// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

// The test package generates test runner programs for packages.
package test

import (
	"container/vector"
	"fmt"
	"igo/set"
)

// GenerateTestMain returns the source code for a test program that will run
// the specified test functions from the specified package.
func GenerateTestMain(packageName string, funcs *set.StringSet) string {
	var funcVec vector.StringVector
	for val := range funcs.Iter() {
		funcVec.Push(val)
	}

	result := ""
	result += "package main\n\n"
	result += "import \"testing\"\n"

	if funcVec.Len() > 0 {
		result += fmt.Sprintf("import \"./%s\"\n\n", packageName)
	}

	result += "var tests = []testing.Test {\n"
	for _, val := range funcVec.Data() {
		result += fmt.Sprintf("\ttesting.Test{\"%s\", %s.%s},\n", val, packageName, val)
	}

	result += "}\n\n"
	result += "func main() {\n"
	result += "\ttesting.Main(tests)\n"
	result += "}\n"

	return result
}
