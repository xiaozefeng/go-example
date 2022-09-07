package web

import (
	"fmt"
	"strings"
)

type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{trees: make(map[string]*node)}
}

func (r *router) addRoute(method, path string, handler HandleFunc) {
	if len(path) == 0 {
		panic("path 不能为空")
	}

	if path[0] != '/' {
		panic("path 必须以 / 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("path 不能以 / 结尾")
	}

	root, ok := r.trees[method]
	if !ok {
		root = &node{path: "/"}
		r.trees[method] = root
	}

	if path == "/" {
		if root.handler != nil {
			panic("路由冲突 [/]")
		}
		root.handler = handler
		return
	}
	segs := strings.Split(path[1:], `/`)

	cur := root
	for _, seg := range segs {
		if len(seg) == 0 {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		cur = cur.createChild(seg)
	}
	if cur.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}
	cur.handler = handler
}

func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	if len(path) == 0 {
		return nil, false
	}
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{node: root}, true
	}
	segs := strings.Split(strings.Trim(path, `/`), `/`)
	matchInfo := &matchInfo{}
	for _, seg := range segs {
		var matchPathParam bool
		root, matchPathParam, ok = root.childOf(seg)
		if !ok {
			return nil, false
		}
		if matchPathParam {
			matchInfo.putValue(root.path[1:], seg)
		}
	}
	matchInfo.node = root
	return matchInfo, true
}

type node struct {
	path    string
	handler HandleFunc

	children map[string]*node

	// 通配符匹配
	starNode *node

	// 路径参数
	pathNode *node
}

type matchInfo struct {
	node       *node
	pathParams map[string]string
}

func (m *matchInfo) putValue(key, val string) {
	if m.pathParams == nil {
		m.pathParams = map[string]string{key: val}
	}
	m.pathParams[key] = val
}

func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.pathNode != nil {
			return n.pathNode, true, true
		}
		return n.starNode, false, n.starNode != nil
	}
	res, ok := n.children[path]
	if !ok {
		if n.pathNode != nil {
			return n.pathNode, true, true
		}
		return n.starNode, false, n.starNode != nil
	}
	return res, false, ok
}

func (n *node) createChild(path string) *node {
	if path == "*" {
		if n.pathNode != nil {
			panic(fmt.Sprintf("web: 非法路由，已经有路径参数路由，不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.starNode == nil {
			n.starNode = &node{path: path}
		}
		return n.starNode
	}
	if path[0] == ':' {
		if n.starNode != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.pathNode != nil {
			if n.pathNode.path != path {
				panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.pathNode.path, path))
			}
		} else {
			n.pathNode = &node{path: path}
		}
		return n.pathNode
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}