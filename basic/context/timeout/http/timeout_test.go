package http

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan struct{})
	go func() {
		doSomething(ctx)
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		log.Println("oh no, I'v exceeded the deadline")
	case <-done:
		log.Println("program exit.")
	}
}

func doSomething(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("timeout")
			return
		default:
			log.Println("doing something coll.")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
