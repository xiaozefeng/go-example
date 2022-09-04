package ast

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

type printVisitor struct {
}

func (p *printVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		fmt.Println(nil)
		return p
	}
	valueOf := reflect.ValueOf(node)
	typeOf := reflect.TypeOf(node)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	fmt.Printf("val: %+v, type:%s\n", valueOf.Interface(), typeOf.Name())
	return p
}

type annotationVisitor struct {
	Annotations map[string]string
}

func (a *annotationVisitor) Visit(node ast.Node) (w ast.Visitor) {
	fn, ok := node.(*ast.File)
	if !ok {
		return a
	}
	list := fn.Doc.List
	if fn.Doc == nil || len(list) == 0 {
		return a
	}
	annotations := make(map[string]string)
	for _, doc := range list {
		if strings.HasPrefix(doc.Text, `// @`) {
			text := strings.TrimPrefix(doc.Text, `// @`)
			sp := strings.SplitN(text, " ", 2)
			if len(sp) > 0 {
				if len(sp[0]) == 0 {
					continue
				}
				var val string
				if len(sp) > 1 {
					val = sp[1]
				}
				annotations[sp[0]] = val
			}
		}
	}
	a.Annotations = annotations
	return a
}
