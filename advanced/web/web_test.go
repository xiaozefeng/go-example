package web

import (
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	s := &MyServer{}
	err := http.ListenAndServe("localhost:8080", s)
	if err != nil {
		return
	}
}
