// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

// The deps package offers utility functions for dealing with dependency
// graphs.
package deps

import (
	"igo/set"
)

// BuildTotalOrder accepts a map from package names to the dependencies of
// those packages, and returns a safe order in which to compile them (assuming
// there are no circular dependencies).
func BuildTotalOrder(dependencies map[string]*set.StringSet) []string {
	return []string{}
}
