package main

import (
	"flag"
	"go/parser"
	"go/token"
	"testing"
)

type testCase struct {
	arg      string
	flags    map[string]string
	expected []string
}

func TestGetFileDeps(t *testing.T) {
	testCases := []testCase{
		{
			arg: "testing",
		},
		{
			arg:      "testing/foo/foo.go",
			expected: []string{"fmt"},
		},
		{
			arg:      "testing/bar/bar.go",
			expected: []string{"fmt"},
		},
		{
			arg:      "testing/testing.go",
			expected: []string{"go-src/testing/foo", "go-src/testing/bar"},
		},
	}

	for i, v := range testCases {
		res, err := getFileDeps(v.arg)
		if v.expected != nil && err != nil {
			t.Errorf("Test[%d]:Expected success but got error: %v", i, err)
			continue
		}
		if err != nil {
			continue
		}
		for _, v := range v.expected {
			var match bool
			for _, v2 := range res {
				if v == v2 {
					match = true
					break
				}
			}
			if !match {
				t.Errorf("Test[%d]:Failed to find match for %s", i, v)
			}
		}
	}

}

func testGetDepfiles(t *testing.T) {

	testCases := []testCase{
		{
			arg: "testing",
		},
		{
			arg:      "testing/foo/foo.go",
			expected: []string{"fmt"},
		},
		{
			arg:      "testing/bar/bar.go",
			expected: []string{"fmt"},
		},
		{
			arg:      "testing/testing.go",
			expected: []string{"go-src/testing/foo", "go-src/testing/bar", "fmt", "fmt"},
		},
	}

	for i, v := range testCases {
		res, err := getDepFiles(v.arg)
		if v.expected != nil && err != nil {
			t.Errorf("Test[%d]:Expected success but got error: %v", i, err)
			continue
		}
		if err != nil {
			continue
		}
		for _, v := range v.expected {
			var match bool
			for i2, v2 := range res {
				if v == v2 {
					match = true
					res = append(res[:i2], res[i2+1:]...)
					break
				}
			}
			if !match {
				t.Errorf("Test[%d]:Failed to find match for %s", i, v)
			}
		}
		if len(res) != 0 {
			t.Errorf("Test[%d]: Extra results %v", i, res)
		}
	}
}

func TestGetFiles(t *testing.T) {
	testCases := []testCase{
		{
			arg:      "testing",
			expected: []string{"testing/testing.go"},
		},
		{
			arg:      "testing",
			expected: []string{"testing/testing.go", "testing/testing_test.go"},
			flags:    map[string]string{"test": "true"},
		},
	}
	for i, v := range testCases {
		for fk, fv := range v.flags {
			flag.Set(fk, fv)
		}
		flag.Parse()
		fs := token.NewFileSet()
		pkgs, err := parser.ParseDir(fs, v.arg, nil, parser.ImportsOnly)
		if err != nil {
			t.Errorf("Failed to parse testing directory")
		}
		files := getFiles(pkgs)
		for _, f := range v.expected {
			var match bool
			for i2, f2 := range files {
				if f == f2 {
					match = true
					files = append(files[:i2], files[i2+1:]...)
					break
				}
			}
			if !match {
				t.Errorf("Test[%d]:Failed to get match for %s", i, f)
			}
		}
		if len(files) != 0 {
			t.Errorf("Test[%d]: Extra results %v", i, files)
		}
	}
}
