package main

import (
	"fmt"
	"sync"
)

func main() {
	c := newCounter()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.increase()
		}()
	}

	wg.Wait()
	fmt.Println("current value:", <-c.c)
}

type counter struct {
	c chan int
	i int
}

func newCounter() *counter {
	c := &counter{
		c: make(chan int),
	}

	go func() {
		for {
			c.c <- c.i
			c.i++
		}
	}()

	return c
}

func (c *counter) increase() int {
	return <-c.c
}
