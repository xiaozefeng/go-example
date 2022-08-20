package insert

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var errInvalidEntity = errors.New("invalid entity")

func GenInsertStmt(target any) (string, []any, error) {

	typeOf := reflect.TypeOf(target)
	//valueOf := reflect.ValueOf(target)
	// 检测 target 是否符合我们的要求
	// 我们只支持有限的几种输入
	if typeOf.Kind() != reflect.Ptr || typeOf.Elem().Kind() != reflect.Struct {
		return "", nil, errInvalidEntity
	}
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	// 使用 strings.Builder 来拼接 字符串
	bd := strings.Builder{}

	// 构造 INSERT INTO XXX，XXX 是你的表名，这里我们直接用结构体名字
	bd.WriteString(fmt.Sprint("INSERT INTO %s", typeOf.Name()))

	// 遍历所有的字段，构造出来的是 INSERT INTO XXX(col1, col2, col3)
	// 在这个遍历的过程中，你就可以把参数构造出来
	// 如果你打算支持组合，那么这里你要深入解析每一个组合的结构体
	// 并且层层深入进去

	// 拼接 VALUES，达成 INSERT INTO XXX(col1, col2, col3) VALUES

	// 再一次遍历所有的字段，要拼接成 INSERT INTO XXX(col1, col2, col3) VALUES(?,?,?)
	// 注意，在第一次遍历的时候我们就已经拿到了参数的值，所以这里就是简单拼接 ?,?,?

	// return bd.String(), args, nil
	panic("implement me")
}
