package close

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func work(ch chan int) {
	for {
		select {
		case t, ok := <-ch:
			if !ok {
				return
			}
			fmt.Printf("task :%d, is done \n", t)
			time.Sleep(time.Millisecond * 5)
		}
	}
}

func sendTask() {
	ch := make(chan int, 10)
	go work(ch)

	for i := 0; i < 1000; i++ {
		ch <- i
	}
	close(ch)
}

func TestCloseChannel(t *testing.T) {
	t.Log(runtime.NumGoroutine())
	sendTask()
	time.Sleep(time.Second * 2)
	t.Log(runtime.NumGoroutine())
}
