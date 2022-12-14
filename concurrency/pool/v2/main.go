package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	fmt.Println("numOfGoroutine:", runtime.NumGoroutine())
	pool := NewPool(10,
		QueueSize(1000),
		Policy(Discard),
		PoolShutdownWaitTime(time.Second*2),
	)
	pool.Start()
	for i := 0; i < 1000; i++ {
		i := i
		pool.Submit(func() {
			fmt.Println("task: do something:", i)
			time.Sleep(500 * time.Millisecond)
		})
	}
	fmt.Println("numOfGoroutine:", runtime.NumGoroutine())
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)
	select {
	case <-done:
		fmt.Println("exit.")
		ctx, cancel := context.WithTimeout(context.Background(), pool.shutdownTimeout)
		defer cancel()
		pool.Shutdown(ctx)
	}
	fmt.Println("numOfGoroutine:", runtime.NumGoroutine())
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
func PoolShutdownWaitTime(timeout time.Duration) Option {
	return func(options *Options) {
		options.shutdownWaitTime = timeout
	}
}
func NewPool(capacity int, opts ...Option) *Pool {
	if capacity <= 0 {
		panic("invalid capacity, too small")
	}
	o := Options{
		queueSize:        100,
		rejectionPolicy:  UseCaller,
		shutdownTimeout:  time.Second * 5,
		shutdownWaitTime: 2 * time.Second,
	}
	for _, opt := range opts {
		opt(&o)
	}

	queue := make(chan task, o.queueSize)
	stop := make(chan struct{})
	p := &Pool{
		queue:            queue,
		rejectionPolicy:  o.rejectionPolicy,
		stop:             stop,
		shutdownTimeout:  o.shutdownTimeout,
		shutdownWaitTime: o.shutdownWaitTime,
		capacity:         capacity,
	}
	return p
}

func (p *Pool) Start() {
	for i := 0; i < p.capacity; i++ {
		i := i
		go work(p.queue, i, p.stop)
	}
}

func work(queue chan task, i int, stop chan struct{}) {
	func() {
		for {
			select {
			case t, ok := <-queue:
				if !ok {
					fmt.Printf("workder:%d ??????????????????\n", i)
					return
				}
				fmt.Printf("worker:%d,?????????????????????\n", i)
				t()
			case <-stop:
				fmt.Printf("workder:%d ??????????????????\n", i)
				return
			}
		}
	}()
}

type Pool struct {
	capacity         int
	queue            chan task // ????????????
	rejectionPolicy  RejectionPolicy
	stop             chan struct{}
	shutdownTimeout  time.Duration
	shutdownWaitTime time.Duration
}

func (p *Pool) Shutdown(ctx context.Context) {
	fmt.Println("start shutdown")
	done := make(chan struct{}, 1)
	go func() {
		for i := 0; i < p.capacity; i++ {
			p.stop <- struct{}{}
			fmt.Println("send stop signal successful")
		}
		close(p.queue)
		done <- struct{}{}
		fmt.Println("send done signal")
	}()
	select {
	case <-ctx.Done():
		fmt.Println("shutdown timeout")
		return
	case <-done:
		fmt.Println("??????????????????worker")
		time.Sleep(p.shutdownWaitTime)
	}
}
func (p *Pool) Submit(task task) {
	select {
	case p.queue <- task:
		fmt.Println("??????????????????")
	default:
		switch p.rejectionPolicy {
		case Discard:
			fmt.Println("????????????")
		case UseCaller:
			task()
		}
	}
}

// ?????????????????????
type task func()

type Option func(*Options)

type Options struct {
	queueSize        int
	rejectionPolicy  RejectionPolicy
	shutdownTimeout  time.Duration
	shutdownWaitTime time.Duration
}

type RejectionPolicy int8

const (
	Discard RejectionPolicy = iota + 1
	UseCaller
)
