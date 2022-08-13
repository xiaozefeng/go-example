package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server, cleanup, err := newHTTPServer(":8080")
	if err != nil {
		panic(err)
	}
	eg, ctx := errgroup.WithContext(context.Background())

	// start server
	eg.Go(func() error {
		fmt.Println("listen http server on ", server.Addr)
		return server.ListenAndServe()
	})

	// shutdown server
	eg.Go(func() error {
		<-ctx.Done()
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return server.Shutdown(c)
	})

	// monitor signal
	eg.Go(func() error {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT)
		select {
		case <-done:
			fmt.Println("优雅退出.")
			cleanup()
			return errors.New("exit")
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	log.Fatal(eg.Wait())
}

func newHTTPServer(addr string) (*http.Server, func(), error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return server, func() {
		fmt.Println("清理资源完成.")
	}, nil
}
