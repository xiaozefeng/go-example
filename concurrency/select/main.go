package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		c1 <- "done"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "done"
	}()
	for i := 0; i < 2; i++ {
		select {
		case v1 := <-c1:
			fmt.Println("received channel 1 value:", v1)
		case v2 := <-c2:
			fmt.Println("received channel 2 value:", v2)
		}
	}

}
