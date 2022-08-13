package main

import (
	"fmt"
	"time"
)

type signal struct {
}

func main() {
	fmt.Println("start a worker...")
	c := spawn(worker)
	<-c
	fmt.Println("worker work done.")
}

func worker() {
	fmt.Println("worker is working...")
	time.Sleep(1 * time.Second)
}

func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		fmt.Println("worker start to work...")
		f()
		c <- signal{}
	}()
	return c
}
