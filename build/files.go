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

	// .go files and their local dependencies necessary for building the package
	// as a library.
	Files *set.StringSet
	Deps  *set.StringSet

	// .go files and their local dependencies necessary for building the package
	// tests, in addition to the ones above.
	TestFiles *set.StringSet
	TestDeps  *set.StringSet

	// Names of test functions within the package.
	TestFuncs *set.StringSet
}

// GetDirectoryInfo scans the supplied directory, determining what package it
// represents (see below), what .go files it contains (test and non-test), and
// what their local package dependencies are.
//
// Sub-directories are not traversed. It is assumed that all of the .go files
// in the directory (not including its sub-directories) belong to the same
// package.
func GetDirectoryInfo(dir string) DirectoryInfo {
	var visitor directoryInfoVisitor
	visitor.originalDir = dir

	path.Walk(dir, &visitor, nil)
	return DirectoryInfo{
		visitor.packageName,
		&visitor.files,
		&visitor.deps,
		&visitor.testFiles,
		&visitor.testDeps,
		&visitor.testFuncs,
	}
}

type directoryInfoVisitor struct {
	originalDir string // The directory supplied by the user.

	packageName string
	files       set.StringSet
	deps        set.StringSet
	testFiles   set.StringSet
	testDeps    set.StringSet
	testFuncs   set.StringSet
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

	// Is this a normal source file or a test file?
	files := &v.files
	deps := &v.deps
	isTest := strings.HasSuffix(file, "_test.go")
	if isTest {
		files = &v.testFiles
		deps = &v.testDeps
	}

	files.Insert(file)

	contents, err := ioutil.ReadFile(file)
	if err == nil {
		for dep := range parse.GetImports(string(contents)).Iter() {
			if strings.HasPrefix(dep, "./") {
				deps.Insert(dep[2:])
			}
		}

		if v.packageName == "" {
			v.packageName = parse.GetPackageName(string(contents))
		}
	}

	if isTest {
		v.testFuncs.Union(parse.GetTestFunctions(string(contents)))
	}
}
