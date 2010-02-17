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

type DirectoryInfo struct {
	PackageName string
	Files       *set.StringSet // .go files to be compiled in the directory.
	Deps        *set.StringSet // Imports of the .go files to be compiled.
}

// GetDirectoryInfo scans the supplied directory, determining what package it
// represents (see below), what .go files it contains, and what their package
// dependencies are. If includeTests is false, files ending in "_test.go" are
// excluded from this determination.
//
// Sub-directories are not traversed. It is assumed that all of the .go files
// in the directory (not including its sub-directories) belong to the same
// package.
func GetDirectoryInfo(dir string, includeTests bool) DirectoryInfo {
	var visitor directoryInfoVisitor
	visitor.includeTests = includeTests
	visitor.originalDir = dir

	path.Walk(dir, &visitor, nil)
	return DirectoryInfo{visitor.packageName, &visitor.files, &visitor.deps}
}

type directoryInfoVisitor struct {
	includeTests bool
	originalDir  string // The directory supplied by the user.

	packageName string
	files       set.StringSet
	deps        set.StringSet
}

func (v *directoryInfoVisitor) VisitDir(dir string, d *os.Dir) bool {
	// Ignore sub-directories, but do recurse into the original directory.
	return dir == v.originalDir
}

func (v *directoryInfoVisitor) VisitFile(file string, d *os.Dir) {
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
		v.deps.Union(parse.GetImports(string(contents)))
		if v.packageName == "" {
			v.packageName = parse.GetPackageName(string(contents))
		}
	}
}
