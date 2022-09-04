package main

import (
	"errors"
	"fmt"
	"github.com/xiaozefeng/go-example/advanced/ast/gen/annotation"
	"github.com/xiaozefeng/go-example/advanced/ast/gen/http"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// 实际上 main 函数这里要考虑接收参数
// src 源目标
// dst 目标目录
// type src 里面可能有很多类型，那么用户可能需要指定具体的类型
// 这里我们简化操作，只读取当前目录下的数据，并且扫描下面的所有源文件，然后生成代码
// 在当前目录下运行 go install 就将 main 安装成功了，
// 可以在命令行中运行 gen
// 在 testdata 里面运行 gen，则会生成能够通过所有测试的代码
func main() {
	err := gen("./testdata")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("success")
}

func gen(src string) error {
	// 第一步找出符合条件的文件
	srcFiles, err := scanFiles(src)
	if err != nil {
		return err
	}
	// 第二步，AST 解析源代码文件，拿到 service definition 定义
	defs, err := parseFiles(srcFiles)
	if err != nil {
		return err
	}
	// 生成代码
	return genFiles(src, defs)
}

// 根据 defs 来生成代码
// src 是源代码所在目录，在测试里面它是 ./testdata
func genFiles(src string, defs []http.ServiceDefinition) error {
	for _, def := range defs {
		name := def.GenName()
		filename := underscoreName(name) + ".go"
		f, err := os.Create(filepath.Join(src, filename))
		if err != nil {
			return err
		}
		err = http.Gen(f, def)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func parseFiles(srcFiles []string) ([]http.ServiceDefinition, error) {
	defs := make([]http.ServiceDefinition, 0, 20)
	for _, src := range srcFiles {
		// 你需要利用 annotation 里面的东西来扫描 src，然后生成 file
		fileSet := token.NewFileSet()
		f, err := parser.ParseFile(fileSet, src, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		visitor := annotation.SingleFileEntryVisitor{}
		ast.Walk(&visitor, f)
		var file = visitor.Get()

		for _, typ := range file.Types {
			_, ok := typ.Annotations.Get("HttpClient")
			if !ok {
				continue
			}
			def, err := parseServiceDefinition(file.Node.Name.Name, typ)
			if err != nil {
				return nil, err
			}
			defs = append(defs, def)
		}
	}
	return defs, nil
}

// 你需要利用 typ 来构造一个 http.ServiceDefinition
// 注意你可能需要检测用户的定义是否符合你的预期
func parseServiceDefinition(pkg string, typ annotation.Type) (http.ServiceDefinition, error) {

	var name = typ.Node.Name.Name
	ans := typ.Ans
	if len(ans) > 0 {
		for _, an := range ans {
			if an.Key == "ServiceName" && len(an.Value) > 0 {
				name = an.Value
			}
		}
	}

	methods := make([]http.ServiceMethod, 0, len(typ.Fields))
	for _, field := range typ.Fields {
		ft, ok := field.Node.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		if ft.Params == nil || len(ft.Params.List) != 2 {
			return http.ServiceDefinition{}, errors.New("gen: 方法必须接收两个参数，其中第一个参数是 context.Context，第二个参数请求")
		}
		if ft.Results == nil || len(ft.Results.List) != 2 {
			return http.ServiceDefinition{}, errors.New("gen: 方法必须返回两个参数，其中第一个返回值是响应，第二个返回值是error")
		}
		pType, ok := ft.Params.List[1].Type.(*ast.StarExpr)
		if !ok {
			return http.ServiceDefinition{}, errors.New("gen: 第二个参数必须是指针")
		}
		rType, ok := ft.Results.List[0].Type.(*ast.StarExpr)
		if !ok {
			return http.ServiceDefinition{}, errors.New("gen: 第一个返回值必须是指针")
		}
		annotations := field.Ans
		var path = "/" + field.Node.Names[0].Name
		for _, anno := range annotations {
			if anno.Key == "Path" {
				path = anno.Value
			}
		}
		pTypeIdent := pType.X.(*ast.Ident)
		rTypeIdent := rType.X.(*ast.Ident)
		var method = http.ServiceMethod{
			Name:         field.Node.Names[0].Name,
			Path:         path,
			ReqTypeName:  pTypeIdent.Name,
			RespTypeName: rTypeIdent.Name,
		}
		methods = append(methods, method)
	}
	var res = http.ServiceDefinition{
		Package: pkg,
		Name:    name,
		Methods: methods,
	}

	return res, nil
}

// 返回符合条件的 Go 源代码文件，也就是你要用 AST 来分析这些文件的代码
func scanFiles(src string) ([]string, error) {
	entries, err := os.ReadDir(src)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			continue
		}
		if !strings.HasSuffix(info.Name(), `.go`) {
			continue
		}
		abs, err := filepath.Abs(src + "/" + info.Name())
		if err != nil {
			return nil, err
		}
		res = append(res, abs)
	}
	return res, nil
}

// underscoreName 驼峰转字符串命名，在决定生成的文件名的时候需要这个方法
// 可以用正则表达式，然而我写不出来，我是正则渣
func underscoreName(name string) string {
	var buf []byte
	for i, v := range name {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
