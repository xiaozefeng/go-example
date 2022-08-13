package main

import (
	"fmt"
	"time"
)

func main() {
	const numberOfJobs = 5
	const numberOfWorkers = 3
	jobs := make(chan int, numberOfJobs)
	results := make(chan int, numberOfJobs)

	for i := 0; i < numberOfWorkers; i++ {
		go worker(i, jobs, results)
	}

	for i := 0; i < numberOfJobs; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 0; i < numberOfJobs; i++ {
		fmt.Println("result: ", <-results)
	}
}

func worker(id int, jobs chan int, results chan int) {
	for job := range jobs {
		fmt.Printf("worker: %d started job: %d \n", id, job)
		time.Sleep(time.Second)
		fmt.Printf("worker: %d finished job: %d\n", id, job)
		results <- job * 2
	}
}
