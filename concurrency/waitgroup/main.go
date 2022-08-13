package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const numberOfJobs = 5
	var wg sync.WaitGroup

	for i := 0; i < numberOfJobs; i++ {
		wg.Add(1)
		i := i

		go func() {
			defer wg.Done()
			worker(i)
		}()
	}

	wg.Wait()
}

func worker(id int) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}
