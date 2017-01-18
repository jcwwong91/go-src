/*
package testing does absolutely nothing but import random librarys to test the
go-src import reading
*/
package testing

import (
	"go-src/testing/bar"
	"go-src/testing/foo"
)

func Test() {
	foo.Foo()
	bar.Bar()
}
