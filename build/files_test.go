// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

package build

import (
	"container/vector"
	"fmt"
	"igo/set"
	"once"
	"os"
	"path"
	"rand"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func seedRand() { rand.Seed(time.Nanoseconds()) }

func createTempDir() string {
	once.Do(seedRand)
	result := fmt.Sprintf("/tmp/files_test.%d", rand.Uint32())
	err := os.Mkdir(result, 0700)
	if err != nil {
		panic(fmt.Sprintf("Can't create dir [%s]: %s", result, err))
	}

	return result
}

func createFile(dir string, name string) *os.File {
	fullPath := path.Join(dir, name)
	result, err := os.Open(fullPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(fmt.Sprintf("Can't open file [%s]: %s", fullPath, err))
	}

	return result
}

func writeFile(file *os.File, contents string) {
	_, err := file.Write(strings.Bytes(contents))
	if err != nil {
		panic(fmt.Sprintf("Couldn't write to file: %s", err))
	}
}

func expectSetContents(t *testing.T, expected []string, set *set.StringSet) {
	var contents vector.StringVector
	for val := range set.Iter() {
		contents.Push(val)
	}

	sort.Sort(&contents)
	sort.SortStrings(expected)

	if !reflect.DeepEqual(contents.Data(), expected) {
		t.Errorf("Expected: %v\nGot: %v", expected, contents.Data())
	}
}

func TestEmptyDir(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	files, deps := GetPackageInfo(dir, true)
	expectSetContents(t, []string{}, files)
	expectSetContents(t, []string{}, deps)
}

func TestOneFile(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	file := createFile(dir, "file.go")
	defer file.Close()

	writeFile(file, `
		// Some comments
		package blah

		import (
			"./foo"
			"fmt"
			"http"
		)

		func DoNothing() {}
	`)

	files, deps := GetPackageInfo(dir, true)
	expectSetContents(t, []string{path.Join(dir, "file.go")}, files)
	expectSetContents(t, []string{"./foo", "fmt", "http"}, deps)
}

func TestSeveralFiles(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	// foo.go
	prodFile1 := createFile(dir, "foo.go")
	defer prodFile1.Close()

	writeFile(prodFile1, `
		package blah
		import "fmt"
		func DoNothing() {}
	`)

	// bar.go
	prodFile2 := createFile(dir, "bar.go")
	defer prodFile2.Close()

	writeFile(prodFile2, `
		package blah
	`)

	// foo_test.go
	testFile1 := createFile(dir, "foo_test.go")
	defer testFile1.Close()

	writeFile(testFile1, `
		package blah
		import "./asdf"
	`)

	// bar_test.go
	testFile2 := createFile(dir, "bar_test.go")
	defer testFile2.Close()

	writeFile(testFile2, `
		package blah
		import "./qwerty"
	`)

	// Without tests
	files, deps := GetPackageInfo(dir, true)
	expectSetContents(t,
		[]string{
			path.Join(dir, "foo.go"),
			path.Join(dir, "bar.go"),
		},
		files)
	expectSetContents(t, []string{"fmt"}, deps)

	// With tests
	files, deps = GetPackageInfo(dir, true)
	expectSetContents(t,
		[]string{
			path.Join(dir, "foo.go"),
			path.Join(dir, "foo_test.go"),
			path.Join(dir, "bar.go"),
			path.Join(dir, "bar_test.go"),
		},
		files)
	expectSetContents(t, []string{"fmt", "./asdf", "./qwerty"}, deps)
}

func TestIgnoresSubdir(t *testing.T) {
	dir := createTempDir()
	defer os.RemoveAll(dir)

	// foo.go, directly in the directory.
	file := createFile(dir, "foo.go")
	defer file.Close()
	writeFile(file, `
		package blah
		import (
			"./foo"
			"http"
		)
	`)

	// bar.go, in a sub-directory.
	subdir := path.Join(dir, "subdir")
	err := os.Mkdir(path.Join(dir, "subdir"), 0700)
	if err != nil {
		panic(fmt.Sprintf("Can't create dir [%s]: %s", subdir, err))
	}

	irrelevantFile := createFile(subdir, "blah.go")
	defer irrelevantFile.Close()
	writeFile(irrelevantFile, `
		package other
		import (
			"./asd"
		)
	`)

	files, deps := GetPackageInfo(dir, true)
	expectSetContents(t, []string{path.Join(dir, "foo.go")}, files)
	expectSetContents(t, []string{"./foo", "http"}, deps)
}
