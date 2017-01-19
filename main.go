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
	"strings"
)

var (
	test = flag.Bool("test", false, "Set to include test files")
)

func usage() {
	fmt.Fprintf(os.Stderr, "go-src <flags> <package>")

}

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	deps, err := getDepFiles(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
	fmt.Println("Total Dependencies: ")
	for _, dep := range deps {
		fmt.Println(dep)
	}
}

func getDepFiles(path string) ([]string, error) {
	var deps []string
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	files := getFiles(pkgs)
	for _, file := range files {
		fDeps, err := getFileDeps(file)
		if err != nil {
			return nil, fmt.Errorf("Error parsing file %s : %v", file, err)
		}
		deps = append(deps, fDeps...)
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
		deps = append(deps, impPath)
	}
	return deps, nil
}
