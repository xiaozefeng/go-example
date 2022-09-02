package ast

import (
	"fmt"
	"go/ast"
	"reflect"
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
