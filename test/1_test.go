package test

import (
	"fmt"
	"hash/crc32"
	"sort"
	"testing"
	"time"
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

var set = make(map[int]bool, 0)

func printOnce(num int) {
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
}
func Test2(t *testing.T) {
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	time.Sleep(time.Second)
}

type aa struct {
	data []int
}

func getA() aa {
	return aa{
		data: []int{1, 2, 3},
	}
}

func Test3(t *testing.T) {
	//a := getA()
	//a.data[0] = 66
	//fmt.Println(a)
	var a interface{}
	a = 1
	//b,ok := a.(int)
	fmt.Println(a.(int))
}
func bbbb() (a int, b bool) {
	return aaaa()
}

func aaaa() (a int, b bool) {
	a++
	b = true
	return
}

func Test98(t *testing.T) {
	a := crc32.ChecksumIEEE([]byte("5lkdgklfjkg"))
	fmt.Println(a)
}

func Test99(t *testing.T) {
	a := []int{6, 3, 7, 34, 4, 66}
	sort.Ints(a)
	key := sort.Search(len(a), func(i int) bool {
		return a[i] == 34
	})
	fmt.Println(a)
	fmt.Println(key)
}
