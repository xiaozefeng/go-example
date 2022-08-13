package main

import "fmt"

func main() {
	f := foo()
	f()
	f()
	f()
	fmt.Println(f())
}

func foo() func() int {
	a := 1
	return func() int {
		a++
		return a
	}
}
