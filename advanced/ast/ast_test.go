package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"testing"
)

func TestInspect(t *testing.T) {
	src := `
package p
const c = 1.0
var X = f(3.14)*2 + c
`
	fs := token.NewFileSet()
	astFile, err := parser.ParseFile(fs, "src.go", src, 0)
	if err != nil {
		t.Error(err)
		return
	}
	ast.Inspect(astFile, func(node ast.Node) bool {
		var s string
		switch x := node.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident:
			s = x.Name
		}
		if len(s) > 0 {
			fmt.Printf("%s:\t%s\n", fs.Position(node.Pos()), s)
		}
		return true
	})
}

func Test(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("########### Manual Iteration ###########")

	fmt.Println("Imports:")
	for _, i := range node.Imports {
		fmt.Println(i.Path.Value)
	}

	fmt.Println("Comments:")
	for _, c := range node.Comments {
		fmt.Print(c.Text())
	}

	fmt.Println("Functions:")
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		fmt.Println(fn.Name.Name)
	}

	fmt.Println("########### Inspect ###########")
	ast.Inspect(node, func(n ast.Node) bool {
		// Find Return Statements
		ret, ok := n.(*ast.ReturnStmt)
		if ok {
			fmt.Printf("return statement found on line %d:\n\t", fset.Position(ret.Pos()).Line)
			printer.Fprint(os.Stdout, fset, ret)
			return true
		}
		// Find Functions
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			var exported string
			if fn.Name.IsExported() {
				exported = "exported "
			}
			fmt.Printf("%sfunction declaration found on line %d: \n\t%s\n", exported, fset.Position(fn.Pos()).Line, fn.Name.Name)
			return true
		}
		return true
	})
	fmt.Println()
}

func TestAst(t *testing.T) {
	fileSet := token.NewFileSet()
	var src = `package ast

import (
	"fmt"
	"go/ast"
	"reflect"
)

type printVisitor struct {
}

func (p *printVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		fmt.Println("node is nil")
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
`
	f, err := parser.ParseFile(fileSet, "src.go", src, parser.ParseComments)
	if err != nil {
		t.Error(err)
		return
	}
	ast.Walk(&printVisitor{}, f)
}
