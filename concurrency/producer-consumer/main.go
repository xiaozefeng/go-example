package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	c := make(chan int, 10)

	go producer(ctx, c)

	go consumer(ctx, c)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	select {
	case <-done:
		cancelFunc()
		fmt.Println("优雅退出")
		time.Sleep(1 * time.Second)
	}
}

var s1 = rand.NewSource(time.Now().UnixNano())

func producer(ctx context.Context, c chan int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("exit producer.")
			return
		case <-time.After(500 * time.Millisecond):
			i := rand.New(s1).Int()
			fmt.Println("send element:", i)
			c <- i
			//case c <- rand.New(s1).Int():
		}
	}
}

func consumer(ctx context.Context, c chan int) {
	for {
		select {
		case e := <-c:
			fmt.Println("receive element:", e)
		case <-ctx.Done():
			fmt.Println("exit consumer.")
			return
		}
	}
}
