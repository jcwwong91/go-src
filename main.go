// Author: Jason Wong
/*
go-src will take any go package and list all the source files that the specified
package depends on
*/
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var (
	root            = string(filepath.Separator) + "usr" + string(filepath.Separator) + "lib" + string(filepath.Separator) + "go"
	test            = flag.Bool("test", false, "Set to include test files")
	printFiles      = flag.Bool("files", false, "Set to print output the source files this is dependent on, other prints the packages")
	includeMainLibs = flag.Bool("main", false, "Set to include looking through main go libraries")

	parsed = make(map[string]bool)
)

func usage() {
	fmt.Fprintf(os.Stderr, "go-src <flags> <package>")

}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	deps, err := getDepFiles(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	fmt.Println("Total Dependencies: ")
	for _, dep := range deps {
		fmt.Println(dep)
	}
}

func getDepFiles(path string) ([]string, error) {
	var deps []string
	var files []string
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	files = getFiles(pkgs)
	for _, file := range files {
		fDeps, err := getFileDeps(file)
		if err != nil {
			return nil, fmt.Errorf("Error parsing file %s : %v", file, err)
		}
		deps = append(deps, fDeps...)
	}

	if *printFiles {
		for _, dep := range deps {
			f, err := getDepFiles(dep)
			if err != nil {
				return nil, err
			}
			files = append(files, f...)
		}
		return files, nil
	}
	return deps, nil
}

func getFiles(pkgs map[string]*ast.Package) []string {
	var files []string
	for _, v := range pkgs {
		for file, _ := range v.Files {
			if !*test && strings.HasSuffix(file, "_test.go") {
				continue
			}
			if parsed[file] {
				continue
			}
			parsed[file] = true
			files = append(files, file)
		}
	}
	return files
}

func getFileDeps(filename string) ([]string, error) {
	var deps []string
	fs := token.NewFileSet()
	aFile, err := parser.ParseFile(fs, filename, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	for _, imp := range aFile.Imports {
		impPath := imp.Path.Value[1 : len(imp.Path.Value)-1] // strip quotes
		if impPath == "C" {
			continue
		}
		if parsed[impPath] {
			continue
		}
		parsed[impPath] = true
		impPath = findPackage(impPath)
		if impPath == "" {
			continue
		}
		deps = append(deps, impPath)
	}
	return deps, nil
}

func findPackage(impPath string) string {
	var path string
	var err error
	if *includeMainLibs {
		path = root + string(filepath.Separator) + "src" + string(filepath.Separator) + impPath
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			return path
		}
	}
	path = os.Getenv("GOPATH") + string(filepath.Separator) + "src" + string(filepath.Separator) + impPath
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		return path
	}
	return ""
}
