// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

// The build package implements the core building logic of igo.
package build

import (
	"igo/parse"
	"igo/set"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// GetPackageInfo computes a set of files that must be compiled to build the
// package contained in the supplied directory, and dependencies of those
// files. If includeTests is true, this will include test files.
func GetPackageInfo(dir string, includeTests bool) (files *set.StringSet, deps *set.StringSet) {
	var visitor packageInfoVisitor
	visitor.includeTests = includeTests
	visitor.originalDir = dir

	path.Walk(dir, &visitor, nil)
	return &visitor.files, &visitor.deps
}

type packageInfoVisitor struct {
	includeTests bool
	originalDir string

	files set.StringSet
	deps set.StringSet
}

func (v *packageInfoVisitor) VisitDir(path string, d *os.Dir) bool {
	// Ignore sub-directories, but do recurse into the original directory.
	return path == v.originalDir
}

func (v *packageInfoVisitor) VisitFile(file string, d *os.Dir) {
	// Ignore files that aren't Go source.
	if path.Ext(file) != ".go" {
		return
  }

	// Ignore test files if we haven't been told to include them.
	if !v.includeTests && strings.HasSuffix(file, "_test.go") {
		return
	}

	v.files.Insert(file)

	contents, err := ioutil.ReadFile(file)
	if err == nil {
		v.deps.Union(parse.ExtractImports(string(contents)))
	}
}
