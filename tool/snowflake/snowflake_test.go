package snowflake

import (
	"sync"
	"testing"
)

func TestWorker_NextID(t *testing.T) {
	w := NewWorker(5, 5)
	var wg sync.WaitGroup
	count := 10000
	ch := make(chan uint64, count)

	wg.Add(count)
	defer close(ch)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id, _ := w.NextID()
			ch <- id
		}()
	}
	wg.Wait()

	m := make(map[uint64]struct{})
	for i := 0; i < count; i++ {
		id := <-ch
		_, ok := m[id]
		if ok {
			t.Fatalf("repeat id %d", id)
			return
		}
		m[id] = struct{}{}
	}

}
