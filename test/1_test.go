package test

import (
	"fmt"
	"testing"
)

type b struct {
	cc int
}
type a struct {
	bb interface{}
}

func Test1(t *testing.T) {
	a1 := a{&b{3}}
	bb, ok := a1.bb.(*b)
	fmt.Println(bb, ok)
}
