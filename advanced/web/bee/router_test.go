package bee

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_addRoute(t *testing.T) {
	var handler HandleFunc = func(ctx *Context) { fmt.Println("hello") }

	tests := []struct {
		method  string
		path    string
		handler HandleFunc
	}{
		{
			method:  http.MethodGet,
			path:    "/",
			handler: handler,
		},
		{
			method:  http.MethodGet,
			path:    "/user",
			handler: handler,
		},
		{
			method:  http.MethodGet,
			path:    "/order/detail",
			handler: handler,
		},
		{
			method:  http.MethodGet,
			path:    "/home",
			handler: handler,
		},
		{
			method:  http.MethodGet,
			path:    "/user/*",
			handler: handler,
		},
		{
			method:  http.MethodGet,
			path:    "/order/detail/:id",
			handler: handler,
		},
	}

	wantRouter := &router{trees: map[string]*node{
		http.MethodGet: {
			path:    "/",
			handler: handler,
			children: map[string]*node{
				"user": &node{
					path:     "user",
					handler:  handler,
					starNode: &node{path: "*", handler: handler},
				},
				"order": &node{
					path: "order",
					children: map[string]*node{
						"detail": &node{path: "detail", handler: handler, pathNode: &node{path: ":id", handler: handler}},
					},
				},
				"home": &node{
					path:    "home",
					handler: handler,
				},
			},
		},
	}}

	actualRouter := newRouter()
	for _, test := range tests {
		actualRouter.addRoute(test.method, test.path, handler)
	}

	_, ok := wantRouter.equal(actualRouter)
	assert.True(t, ok)

	findTests := []struct {
		method   string
		path     string
		found    bool
		wantPath string
	}{
		{
			method:   http.MethodGet,
			path:     "/",
			found:    true,
			wantPath: "/",
		},
		{
			method:   http.MethodGet,
			path:     "/order/detail",
			found:    true,
			wantPath: "detail",
		},
		{
			method:   http.MethodGet,
			found:    true,
			path:     "/order",
			wantPath: "order",
		},

		{
			method:   http.MethodGet,
			found:    true,
			path:     "/user/*",
			wantPath: "*",
		},
		{
			method:   http.MethodGet,
			found:    true,
			path:     "/order/detail/:id",
			wantPath: ":id",
		},
	}

	for _, tt := range findTests {
		matchInfo, found := actualRouter.findRoute(tt.method, tt.path)
		assert.Equal(t, tt.found, found)
		if found {
			assert.Equal(t, tt.wantPath, matchInfo.node.path)
		}
	}

}

func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}
