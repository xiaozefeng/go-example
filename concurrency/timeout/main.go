package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "done"
	}()
	select {
	case res := <-c1:
		fmt.Println("received from c1: ", res)
	case <-time.After(time.Second * 1):
		fmt.Println("received from c1 leakTimeout")

	}
	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "done"
	}()
	select {
	case res := <-c2:
		fmt.Println("received from c2: ", res)
	case <-time.After(time.Second * 3):
		fmt.Println("received from c2 leakTimeout")

	}

}
