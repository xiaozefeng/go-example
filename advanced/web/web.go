package web

import (
	"context"
	"net"
	"net/http"
)

type Context struct {
	context.Context
	Request    *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string
}

func (c *Context) WriteString(content string) error {
	_, err := c.Resp.Write([]byte(content))
	return err
}

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler

	Start(addr string) error

	AddRoute(method, path string, handler HandleFunc)
}

type Group struct {
	prefix string
	s      Server
}

func (g *Group) AddRoute(method, path string, handler HandleFunc) {
	g.s.AddRoute(method, g.prefix+path, handler)
}

func NewServer() *MyServer {
	return &MyServer{router: newRouter()}
}

type MyServer struct {
	router *router
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
	m.router.addRoute(method, path, handler)
}

func (m *MyServer) Group(prefix string) *Group {
	return &Group{prefix: prefix, s: m}
}

func (m *MyServer) serve(ctx *Context) {
	request := ctx.Request
	matchInfo, found := m.router.findRoute(request.Method, request.URL.Path)
	writer := ctx.Resp
	if !found || matchInfo.node.handler == nil {
		// 404
		http.NotFound(writer, request)
		return
	}

	ctx.PathParams = matchInfo.pathParams
	matchInfo.node.handler(ctx)
}

func (m *MyServer) Get(path string, handler HandleFunc) {
	m.AddRoute(http.MethodGet, path, handler)
}

func (m *MyServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Request: request,
		Resp:    writer,
	}

	// 查找路由
	// 执行方法
	m.serve(ctx)
}
