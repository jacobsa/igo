package main

import (
	"container/vector"
	"flag"
	"fmt"
	"igo/build"
	"igo/deps"
	"igo/set"
	"os"
	"path"
	"strings"
)

// executeCommand runs the specified tool with the supplied arguments (not
// including the path to the tool itself), chdir'ing to the specified directory
// first. It returns true if and only if the child process returns zero.
func executeCommand(tool string, args []string, dir string) bool {
	fmt.Printf("%s %s\n", tool, strings.Join(args, " "))

	var fullArgs vector.StringVector
	fullArgs.Push(tool)
	fullArgs.AppendVector(&args)

	pid, err := os.ForkExec(
		tool,
		fullArgs.Data(),
		os.Environ(),
		dir,
		[]*os.File{os.Stdin, os.Stdout, os.Stderr})
	if err != nil {
		panic(err)
	}

	waitMsg, err := os.Wait(pid, 0)
	if err != nil {
		panic(err)
	}

	return waitMsg.ExitStatus() == 0
}

// compileFiles invokes 6g with the appropriate arguments for compiling the
// supplied set of .go files, and exits the program if the subprocess fails.
func compileFiles(files *set.StringSet, targetBaseName string) {
	compilerPath := path.Join(os.Getenv("GOBIN"), "6g")
	gopackPath := path.Join(os.Getenv("GOBIN"), "gopack")

	targetDir, _ := path.Split(targetBaseName)
	if targetDir != "" {
		os.MkdirAll(path.Join("igo-out", targetDir), 0700)
	}

	// Compile
	var compilerArgs vector.StringVector
	compilerArgs.Push("-o")
	compilerArgs.Push(targetBaseName + ".6")

	for file := range files.Iter() {
		compilerArgs.Push(path.Join("../", file))
	}

	if !executeCommand(compilerPath, compilerArgs.Data(), "igo-out/") {
		os.Exit(1)
	}

	// Pack
	var gopackArgs vector.StringVector
	gopackArgs.Push("grc")
	gopackArgs.Push(targetBaseName + ".a")
	gopackArgs.Push(targetBaseName + ".6")

	if !executeCommand(gopackPath, gopackArgs.Data(), "igo-out/") {
		os.Exit(1)
	}
}

// linkBinary calls 6l to link the binary of the given name, which must have
// already been compiled with compileFiles.
func linkBinary(name string) {
	linkerPath := path.Join(os.Getenv("GOBIN"), "6l")

	var linkerArgs vector.StringVector
	linkerArgs.Push("-o")
	linkerArgs.Push(name)
	linkerArgs.Push(name + ".6")

	if !executeCommand(linkerPath, linkerArgs.Data(), "igo-out/") {
		os.Exit(1)
	}
}

func printUsageAndExit() {
	fmt.Println("Usage:")
	fmt.Println("  igo build <directory name>")
	os.Exit(1)
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		printUsageAndExit()
	}

	command := flag.Arg(0)
	if command != "build" {
		printUsageAndExit()
	}

	// Grab dependency and file information for every local package, starting
	// with the specified one. We consider a package local if it starts with "./".
	requiredFiles := make(map[string]*set.StringSet)
	packageDeps := make(map[string]*set.StringSet)

	specifiedDir := flag.Arg(1)
	var remainingDirs vector.StringVector
	remainingDirs.Push("./" + specifiedDir)

	for remainingDirs.Len() > 0 {
		dir := remainingDirs.Pop()

		// Have we already processed this directory?
		_, alreadyDone := packageDeps[dir]
		if alreadyDone {
			continue
		}

		dirInfo := build.GetDirectoryInfo(dir, false)
		if dirInfo.PackageName == "" {
			fmt.Printf("Couldn't find .go files to build in directory: %s\n", dir)
			os.Exit(1)
		}

		// Stash information about this package, and add its local dependencies to
		// the queue.
		requiredFiles[dir] = dirInfo.Files
		packageDeps[dir] = dirInfo.Deps

		for dep := range dirInfo.Deps.Iter() {
			if strings.HasPrefix(dep, "./") {
				remainingDirs.Push(dep)
			}
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
		compileFiles(requiredFiles[currentPackage], currentPackage)
	}

	// If this is a binary, also link it.
	if build.GetDirectoryInfo(specifiedDir, false).PackageName == "main" {
		linkBinary(specifiedDir)
	}
}
