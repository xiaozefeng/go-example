package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan int, 5)
	for i := 0; i < 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.Tick(200 * time.Millisecond)

	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}
	fmt.Println("----------")
	burstyLimiter := make(chan time.Time, 3)

	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	burstryRequests := make(chan int, 5)
	for i := 0; i < 5; i++ {
		burstryRequests <- i
	}
	close(burstryRequests)

	for req := range burstryRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}

}
