// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

// The deps package offers utility functions for extracting dependency
// information from Go source files, and working with that information.
package deps

import (
	"os"
)

// ExtractDependencies parses the supplied source code for a .go file and
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
func ExtractDependencies(source string) (deps []string, err os.Error) {
	return nil, os.NewError("Not implemented.")
}
