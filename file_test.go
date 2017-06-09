package hamgo

import "testing"

func Test_currentPath(t *testing.T) {
	p1 := "/1/test.log"
	p2 := "/1/2/test.log"
	p3 := "test.log"
	p4 := "./test.log"
	println(p1, currentPath(p1))
	println(p2, currentPath(p2))
	println(p3, currentPath(p3))
	println(p4, currentPath(p4))
}
