package escape

import (
	"fmt"
	"sync"
)

type user struct {
	name string
}

// 返回指针
func newUser(name string) *user {
	return &user{name: name}
}

func printSomething() {

	var mux sync.Mutex
	mux.Unlock()
	str := "Hello"
	// 参数是指针
	fmt.Println(str)
}

// 返回接口
func returnAny() any {
	return "any thing"
}

// 闭包
func countFn() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

// 使用 channel 传递指针
func sendPointToChannel() {
	ch := make(chan *user, 1)
	go func() {

		u := &user{}
		ch <- u
	}()
	<-ch
}
