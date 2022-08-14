package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	pool := NewPool(10, QueueSize(1000), Policy(Discard))
	for i := 0; i < 1000; i++ {
		i := i
		pool.Submit(func() {
			fmt.Println("task: do something:", i)
			time.Sleep(500 * time.Millisecond)
		})
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)
	select {
	case <-done:
		fmt.Println("exit.")
		ctx, cancel := context.WithTimeout(context.Background(), pool.shutdownTimeout)
		defer cancel()
		pool.Shutdown(ctx)
	}
}

func QueueSize(queueSize int) Option {
	return func(options *Options) {
		options.queueSize = queueSize
	}
}
func Policy(policy RejectionPolicy) Option {
	return func(options *Options) {
		options.rejectionPolicy = policy
	}
}
func PoolShutdownTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.shutdownTimeout = timeout
	}
}
func NewPool(capacity int, opts ...Option) *Pool {
	if capacity <= 0 {
		panic("invalid capacity, too small")
	}
	o := Options{
		queueSize:       100,
		rejectionPolicy: UseCaller,
		shutdownTimeout: time.Second * 5,
	}
	for _, opt := range opts {
		opt(&o)
	}

	queue := make(chan task, o.queueSize)
	stop := make(chan struct{})
	p := &Pool{
		queue:           queue,
		rejectionPolicy: o.rejectionPolicy,
		stop:            stop,
		shutdownTimeout: o.shutdownTimeout,
	}
	for i := 0; i < capacity; i++ {
		i := i
		go work(queue, i, stop)
	}
	return p
}

func work(queue chan task, i int, stop chan struct{}) {
	func() {
		for {
			select {
			case t := <-queue:
				fmt.Printf("worker:%d,从队列捞到任务\n", i)
				t()
			case <-stop:
				fmt.Printf("workder:%d 后台任务退出\n", i)
				return
			}
		}
	}()
}

type Pool struct {
	capacity        int
	queue           chan task // 任务队列
	rejectionPolicy RejectionPolicy
	stop            chan struct{}
	shutdownTimeout time.Duration
}

func (p *Pool) Shutdown(ctx context.Context) {
	fmt.Println("start shutdown")
	done := make(chan struct{}, 1)
	go func() {
		for i := 0; i < 10; i++ {
			p.stop <- struct{}{}
			fmt.Println("send stop signal successful")
		}
		done <- struct{}{}
		fmt.Println("send done signal")
	}()
	select {
	case <-ctx.Done():
		fmt.Println("shutdown timeout")
	case <-done:
		fmt.Println("关闭所有后台worker")
	}
}
func (p *Pool) Submit(task task) {
	select {
	case p.queue <- task:
		fmt.Println("提交任务成功")
	default:
		switch p.rejectionPolicy {
		case Discard:
			fmt.Println("任务丢弃")
		case UseCaller:
			task()
		}
	}
}

// 简化的任务模型
type task func()

type Option func(*Options)

type Options struct {
	queueSize       int
	rejectionPolicy RejectionPolicy
	shutdownTimeout time.Duration
}

type RejectionPolicy int8

const (
	Discard RejectionPolicy = iota + 1
	UseCaller
)
