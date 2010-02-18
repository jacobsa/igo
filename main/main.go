package main

import (
	"container/vector"
	"flag"
	"fmt"
	"igo/build"
	"igo/deps"
	"igo/set"
	"igo/test"
	"io/ioutil"
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
	if command != "build" && command != "test" {
		printUsageAndExit()
	}

	// Grab dependency and file information for every local package, starting
	// with the specified one. We consider a package local if it starts with "./".
	requiredFiles := make(map[string]*set.StringSet)
	packageDeps := make(map[string]*set.StringSet)

	specifiedPackage := flag.Arg(1)
	var remainingPackages vector.StringVector
	remainingPackages.Push(specifiedPackage)

	for remainingPackages.Len() > 0 {
		packageName := remainingPackages.Pop()

		// Have we already processed this directory?
		_, alreadyDone := packageDeps[packageName]
		if alreadyDone {
			continue
		}

		dir := "./" + packageName
		dirInfo := build.GetDirectoryInfo(dir)
		if dirInfo.PackageName == "" {
			fmt.Printf("Couldn't find .go files to build in directory: %s\n", dir)
			os.Exit(1)
		}

		// Stash information about this package, and add its local dependencies to
		// the queue.
		requiredFiles[packageName] = dirInfo.Files
		packageDeps[packageName] = dirInfo.Deps

		// If we're testing and this is the package under test, also add its test
		// files and dependencies.
		if packageName == specifiedPackage && command == "test" {
			requiredFiles[packageName].Union(dirInfo.TestFiles)
			packageDeps[packageName].Union(dirInfo.TestDeps)
		}

		for dep := range packageDeps[packageName].Iter() {
			remainingPackages.Push(dep)
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
	if build.GetDirectoryInfo(specifiedPackage).PackageName == "main" {
		linkBinary(specifiedPackage)
	}

	// If we're testing, create a test runner, build it, and run it.
	if command == "test" {
		const outputFile = "igo-out/test_runner.go"
		testFuncs := build.GetDirectoryInfo(specifiedPackage).TestFuncs
		code := test.GenerateTestMain(specifiedPackage, testFuncs)
		err := ioutil.WriteFile("igo-out/test_runner.go", strings.Bytes(code), 0600)
		if err != nil {
			panic(err)
		}

		var files set.StringSet
		files.Insert("igo-out/test_runner.go")
		compileFiles(&files, "test_runner")
		linkBinary("test_runner")

		if !executeCommand("igo-out/test_runner", []string{}, "") {
			os.Exit(1)
		}
	}
}
