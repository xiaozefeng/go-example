package main

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func leakTimeout() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Microsecond)
	defer cancelFunc()
	// 正确的做法是这里改为 buffer 为 1
	done := make(chan struct{})
	go func() {
		time.Sleep(200 * time.Millisecond)
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
	case <-done:
	}
}

func TestLeak(t *testing.T) {
	for i := 0; i < 1000; i++ {
		leakTimeout()
	}
	time.Sleep(3 * time.Second)
	t.Log(runtime.NumGoroutine())
}
