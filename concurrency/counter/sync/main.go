package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			increase()
		}()
	}

	wg.Wait()
	fmt.Printf("current value: %d\n", c.i)

}

type counter struct {
	sync.Mutex
	i int
}

var c counter

func increase() int {
	c.Lock()
	defer c.Unlock()
	c.i++
	return c.i
}
