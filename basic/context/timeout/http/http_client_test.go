package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func get(url string, timeout time.Duration) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	req = req.WithContext(ctx)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func TestHTTPClientTimeout(t *testing.T) {
	bytes, err := get("http://httpbin.org/get", time.Millisecond*200)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("resp: %s", bytes)
}
