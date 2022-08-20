package channel

import (
	"fmt"
	"testing"
	"time"
)

/*
*

	用 channel 实现不断输出 1234的程序
*/
func TestChannel(t *testing.T) {
	c1 := make(chan struct{})
	c2 := make(chan struct{})
	c3 := make(chan struct{})
	c4 := make(chan struct{}, 1)
	c4 <- struct{}{}
	go work(c4, c1, 1)
	go work(c1, c2, 2)
	go work(c2, c3, 3)
	go work(c3, c4, 4)
	time.Sleep(5 * time.Second)
}

func work(c chan struct{}, owner chan struct{}, id int) {
	for {
		<-c
		fmt.Println(id)
		time.Sleep(100 * time.Millisecond)
		owner <- struct{}{}
	}
}
