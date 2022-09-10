package bee

import (
	"fmt"
	"regexp"
	"strings"
)

type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{trees: make(map[string]*node)}
}

type nodeType int

const (
	// 静态路由
	nodeTypeStatic = iota
	// 正则路由
	nodeTypeReg
	// 路径参数路由
	nodeTypeParam
	// 通配符路由
	nodeTypeAny
)

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
		root = &node{path: "/", typ: nodeTypeStatic}
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
	cur.route = path
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
			if root.typ == nodeTypeParam {
				matchInfo.putValue(root.paramName, seg)
			} else if root.typ == nodeTypeReg {
				subMatch := root.regExpr.FindStringSubmatch(seg)
				if len(subMatch) > 1 {
					matchInfo.putValue(root.paramName, subMatch[1])
				}
			}
		}
	}
	matchInfo.node = root
	return matchInfo, true
}

type node struct {
	typ     nodeType
	path    string
	handler HandleFunc
	route   string

	children map[string]*node

	// 通配符匹配
	starNode *node

	// 路径参数
	pathNode *node
	// 正则路由和参数路由都会使用这个字段
	paramName string

	regNode *node
	regExpr *regexp.Regexp
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
		if n.regNode != nil {
			matched := n.regNode.regExpr.MatchString(path)
			return n.regNode, true, matched
		}
		if n.pathNode != nil {
			return n.pathNode, true, true
		}
		if n.starNode == nil && n.path == "*" {
			return n, false, true
		}
		return n.starNode, false, n.starNode != nil
	}
	res, ok := n.children[path]
	if !ok {
		if n.regNode != nil {
			matched := n.regNode.regExpr.MatchString(path)
			return n.regNode, true, matched
		}
		if n.pathNode != nil {
			return n.pathNode, true, true
		}
		return n.starNode, false, n.starNode != nil
	}
	return res, false, ok
}

func (n *node) createChild(path string) *node {
	if strings.ContainsAny(path, `()`) {
		seg := strings.Split(path, `(`)
		if len(seg) != 2 {
			panic(fmt.Sprintf("web: 非法路由，不符合正则规范, 必须是 :name(你的正则)的格式 [%s]", path))
		}
		var paramName = seg[0][1:]
		reg := regexp.MustCompile(`(` + seg[1])
		if n.pathNode != nil {
			panic(fmt.Sprintf("web: 非法路由，已经有路径参数路由，不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.starNode != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.regNode != nil {
			panic(fmt.Sprintf("web: 非法路由, 重复注册正则路由 [%s]", path))
		}
		n.regNode = &node{path: path, typ: nodeTypeReg, regExpr: reg, paramName: paramName}
		//n.regExpr = reg
		//n.paramName = paramName
		return n.regNode
	}
	if path == "*" {
		if n.pathNode != nil {
			panic(fmt.Sprintf("web: 非法路由，已经有路径参数路由，不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.starNode == nil {
			n.starNode = &node{path: path, typ: nodeTypeAny}
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
			n.pathNode = &node{path: path, typ: nodeTypeParam, paramName: path[1:]}
		}
		return n.pathNode
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path, typ: nodeTypeStatic}
		n.children[path] = child
	}
	return child
}
