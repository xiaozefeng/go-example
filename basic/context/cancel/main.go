package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	for i := 1; i <= 3; i++ {
		i := i
		go func(ctx context.Context) {
			backgroundTask(ctx, i)
		}(ctx)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	select {
	case <-done:
		log.Println("收到信号: 停止所有任务")
		cancel()
	}
}

func backgroundTask(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("worker:%d 取消任务\n", id)
			return
		default:
			log.Printf("worker:%d 执行任务中\n", id)
			// 模拟任务延时
			time.Sleep(time.Millisecond * 1000)
		}
	}
}
