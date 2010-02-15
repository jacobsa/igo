// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

// The build package implements the core building logic of igo.
package build

import (
	"igo/set"
)

// GetPackageInfo computes a set of files that must be compiled to build the
// package contained in the supplied directory, and dependencies of those
// files. If includeTests is true, this will include test files.
func GetPackageInfo(dir string, includeTests bool) (files *set.StringSet, deps *set.StringSet) {
	deps = &set.StringSet{}
	files = &set.StringSet{}
	return
}
