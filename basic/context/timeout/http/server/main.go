package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", serve)
	log.Fatal(http.ListenAndServe(":80", nil))

}

func serve(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	fmt.Println("server: handler started")
	defer fmt.Println("server: handler ended")

	select {
	case <-time.After(3 * time.Second):
		_, _ = fmt.Fprintf(w, "hello\n")
	case <-ctx.Done():
		err := ctx.Err()
		fmt.Println("server err:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
