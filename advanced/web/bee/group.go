package bee

import "net/http"

type Group struct {
	prefix string
	s      Server
}

func (g *Group) Get(path string, handler HandleFunc) {
	g.s.addRoute(http.MethodGet, g.prefix+path, handler)
}
func (g *Group) Post(path string, handler HandleFunc) {
	g.s.addRoute(http.MethodPost, g.prefix+path, handler)
}
func (g *Group) Put(path string, handler HandleFunc) {
	g.s.addRoute(http.MethodPut, g.prefix+path, handler)
}
