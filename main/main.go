package main

import (
	"container/vector"
	"flag"
	"fmt"
	"igo/build"
	"igo/deps"
	"igo/set"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	if flag.NArg() != 2 || flag.Arg(0) != "build" {
		fmt.Println("Usage: igo build <directory name>")
		os.Exit(1)
	}

	// Grab dependency and file information for every local package, starting
	// with the specified one. We consider a package local if it starts with "./".
	requiredFiles := make(map[string]*set.StringSet)
	packageDeps := make(map[string]*set.StringSet)

	var remainingDirs vector.StringVector
	remainingDirs.Push(flag.Arg(1))

	for remainingDirs.Len() > 0 {
		dir := remainingDirs.Pop()

		dirInfo := build.GetDirectoryInfo(dir, false)
		if dirInfo.PackageName == "" {
			fmt.Printf("Couldn't find .go files to build in directory: %s\n", dir)
			os.Exit(1)
		}

		// Have we already processed this package?
		//
		// TODO(jacobsa): It would be more efficient to do this chack before we hit
		// the filesystem, based on directory name rather than package name.
		_, alreadyDone := packageDeps[dirInfo.PackageName]
		if alreadyDone { continue }

		// Stash information about this package, and add its local dependencies to
		// the queue.
		requiredFiles[dirInfo.PackageName] = dirInfo.Files
		packageDeps[dirInfo.PackageName] = dirInfo.Deps

		for dep := range dirInfo.Deps.Iter() {
			if strings.HasPrefix(dep, "./") { remainingDirs.Push(dep) }
		}
	}

	// Order the packages by their dependencies.
	totalOrder := deps.BuildTotalOrder(packageDeps)
	fmt.Println("Found these packages to compile:")
	for _, packageName := range totalOrder {
		fmt.Printf("  %s\n", packageName)
	}

	// Create a directory to hold outputs, deleting the old one first.
	os.RemoveAll("igo-out")
	os.Mkdir("igo-out", 0700)

	// Compile each of the packages in turn.
	for _, currentPackage := range totalOrder {
		fmt.Printf("\nCompiling package: %s\n", currentPackage)
	}
}
