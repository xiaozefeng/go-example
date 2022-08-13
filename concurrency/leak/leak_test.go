package leak

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

// 没有任何逻辑往ch中发数据
func leak() {
	ch := make(chan int)

	go func() {
		val := <-ch
		fmt.Println("received a value", val)
	}()
}

func TestLeak(t *testing.T) {
	err := process("hello")
	if err != nil {
		t.Error(err)
	}

	time.Sleep(1 * time.Second)
}

func search(term string) (string, error) {
	time.Sleep(500 * time.Millisecond)
	return "some value", nil
}

type result struct {
	val string
	err error
}

func process(term string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	// 这里应该用 buffer channel ，避免泄露
	// 	results := make(chan result, 1)
	results := make(chan result)
	go func() {
		record, err := search(term)
		results <- result{record, err}
		log.Printf("发送数据成功")
	}()
	select {
	case <-ctx.Done():
		log.Println("timeout")
	case result := <-results:
		if result.err != nil {
			log.Println(result.err)
		}
		log.Println("received data", result.val)
	}
	return nil
}
