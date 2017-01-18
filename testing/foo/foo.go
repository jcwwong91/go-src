/*
Package foo does nothing but call random libraries to test the go-src library
*/
package foo

import (
	"fmt"
)

func Foo() {
	fmt.Sprintf("=)")
}
