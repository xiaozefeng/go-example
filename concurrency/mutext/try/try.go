package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	go func() {
		mu.Lock()
		defer mu.Unlock()
		time.Sleep(2 * time.Second)
	}()
	time.Sleep(time.Millisecond * 100)
	ok := mu.TryLock()
	if !ok {
		fmt.Println("can not get lock")
		return
	}
	fmt.Println("get lock")
}
