package web

import (
	"context"
	"net"
	"net/http"
)

type Context struct {
	context.Context
	request *http.Request
	writer  http.ResponseWriter
}

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler

	Start(addr string) error

	AddRoute(method, path string, handler HandleFunc)
}

type MyServer struct {
}

func (m *MyServer) Start(addr string) error {
	// pre
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// post
	return http.Serve(listen, m)
}

func (m *MyServer) AddRoute(method, path string, handler HandleFunc) {
	// gin 的路由树 每个method 都会有颗路由树
	panic("implement me")
}

func (m *MyServer) Get(path string, handler HandleFunc) {
	m.AddRoute(http.MethodGet, path, handler)
}

func (m *MyServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		request: request,
		writer:  writer,
	}

	// 查找路由
	// 执行方法
	m.serve(ctx)
}

func (m *MyServer) serve(ctx *Context) {

}
