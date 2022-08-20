package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println(unsafe.Sizeof(T{}))
	var c chan string
	fmt.Println(unsafe.Sizeof(c))
}

type T struct {
	a int8
	b int64
	c int16
	d int8
	e chan int
}
