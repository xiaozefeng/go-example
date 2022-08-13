package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("start a group of workers...")
	groupSignal := make(chan signal)
	c := spawnGroup(worker, 5, groupSignal)
	time.Sleep(1 * time.Second)
	fmt.Println("the group of workers start to work...")
	close(groupSignal)
	<-c
	fmt.Println("the group of workers work done.")
}

func worker(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: is work done. \n", i)
}

type signal struct{}

func spawnGroup(f func(i int), num int, groupSignal chan signal) <-chan signal {
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		i := i + 1
		go func() {
			<-groupSignal
			fmt.Printf("worker %d : start to work...\n", i)
			f(i)
			wg.Done()
		}()
	}

	c := make(chan signal)
	go func() {
		wg.Wait()
		c <- signal{}
	}()

	return c
}
