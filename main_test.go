package main

import (
	"testing"
)

type testCase struct {
	arg      string
	flags    []bool
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
	}
}
