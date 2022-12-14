package bee

import (
	"log"
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler

	Start(addr string) error

	addRoute(method, path string, handler HandleFunc, middlewares ...Middleware)
}

func NewServer() *HTTPServer {
	return &HTTPServer{router: newRouter()}
}

type HTTPServer struct {
	router      *router
	middlewares []Middleware
}

func (s *HTTPServer) Start(addr string) error {
	// pre
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// post
	return http.Serve(listen, s)
}

func (s *HTTPServer) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *HTTPServer) addRoute(method, path string, handler HandleFunc, middlewares ...Middleware) {
	s.router.addRoute(method, path, handler, middlewares...)
}

func (s *HTTPServer) Group(prefix string) *Group {
	return &Group{prefix: prefix, s: s}
}

func (s *HTTPServer) serve(ctx *Context) {
	request := ctx.Request
	matchInfo, found := s.router.findRoute(request.Method, request.URL.Path)
	if !found || matchInfo.node == nil || matchInfo.node.handler == nil {
		// 404
		ctx.RespStatusCode = http.StatusNotFound
		return
	}

	ctx.PathParams = matchInfo.pathParams
	ctx.MatchedRoute = matchInfo.node.route

	var root = matchInfo.node.handler
	if matchInfo.middlewares != nil && len(matchInfo.middlewares) > 0 {
		for i := len(matchInfo.middlewares) - 1; i >= 0; i-- {
			root = matchInfo.middlewares[i](root)
		}
	}
	root(ctx)
}

func (s *HTTPServer) Get(path string, handler HandleFunc, middlewares ...Middleware) {
	s.addRoute(http.MethodGet, path, handler, middlewares...)
}
func (s *HTTPServer) Post(path string, handler HandleFunc, middlewares ...Middleware) {
	s.addRoute(http.MethodPost, path, handler, middlewares...)
}
func (s *HTTPServer) Put(path string, handler HandleFunc, middlewares ...Middleware) {
	s.addRoute(http.MethodPut, path, handler, middlewares...)
}

func (s *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Fatalln("??????????????????", err)
	}
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Request: request,
		Resp:    writer,
	}

	root := s.serve
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		m := s.middlewares[i]
		root = m(root)
	}

	// ?????????????????????
	// ????????? next ??????????????????
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}
