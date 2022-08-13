package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

func main() {
	pool := NewPool(Capacity(2), QueueSize(2), Policy(Discard))
	for i := 0; i < 1000; i++ {
		i := i
		pool.Submit(func() {
			fmt.Println("task: do something:", i)
		})
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)
	select {
	case <-done:
		fmt.Println("exit.")
	}
}
func Capacity(capacity int32) Option {
	return func(options *Options) {
		if capacity <= 1 {
			panic("invalid capacity, too small")
		}
		options.capacity = capacity
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
func NewPool(opts ...Option) *Pool {
	o := Options{
		capacity:        10,
		queueSize:       100,
		rejectionPolicy: UseCaller,
	}
	for _, opt := range opts {
		opt(&o)
	}

	queue := make(chan task, o.queueSize)
	go func() {
		for {
			select {
			case t := <-queue:
				fmt.Println("从队列捞到任务")
				t()
			}
		}
	}()
	return &Pool{
		capacity:        o.capacity - 1,
		numOfActive:     0,
		queue:           queue,
		queueSize:       o.queueSize,
		rejectionPolicy: o.rejectionPolicy,
	}

}

type Pool struct {
	capacity        int32     //最大 goroutine数量
	numOfActive     int32     // 当前活动数量
	queue           chan task // 任务队列
	rejectionPolicy RejectionPolicy
	queueSize       int
}

func (p *Pool) Submit(task task) {
	// 获取当前活跃goroutine数量
	activeNum := atomic.LoadInt32(&p.numOfActive)
	if activeNum >= p.capacity {
		select {
		case p.queue <- task:
			fmt.Println("任务成功加入队列")
		default:
			// 如果超过了 capacity 加入队列
			if len(p.queue) == p.queueSize {
				fmt.Println("队列已满")
				switch p.rejectionPolicy {
				case Discard:
					fmt.Println("任务被丢弃")
				case UseCaller:
					fmt.Println("使用当前goroutine执行任务")
					task()
				}

			}
		}
	} else {
		atomic.AddInt32(&p.numOfActive, 1)
		go func() {
			defer atomic.StoreInt32(&p.numOfActive, atomic.LoadInt32(&p.numOfActive)-1)
			task()
		}()
	}
}

// 简化的任务模型
type task func()

type Option func(*Options)

type Options struct {
	capacity        int32
	queueSize       int
	rejectionPolicy RejectionPolicy
}

type RejectionPolicy int8

const (
	Discard RejectionPolicy = iota + 1
	UseCaller
)
