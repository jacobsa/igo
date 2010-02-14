// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

// The parse package offers utility functions for extracting dependency
// information from Go source files.
package parse

// ExtractImports parses the supplied source code for a .go file and
// returns an array of package names that the file depends upon.
//
// For example, if source looks like the following:
//
//     import (
//       "./bar/baz"
//       "fmt"
//       "os"
//     )
//
//     func DoSomething() {
//       ...
//     }
//
// then the result will be [ "./bar/baz", "fmt", "os" ].
//
// An attempt is made to return the imports for the file even if there is a
// syntax error elsewhere in the file.
func ExtractImports(source string) []string {
	return []string{}
}
