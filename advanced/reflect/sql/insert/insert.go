package insert

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var errInvalidEntity = errors.New("invalid entity")

func GenInsertStmt(target any) (string, []any, error) {
	if target == nil {
		return "", nil, errInvalidEntity
	}

	valueOf := reflect.ValueOf(target)
	if valueOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			return "", nil, errInvalidEntity
		}
	}
	typeOf := reflect.TypeOf(target)
	//valueOf := reflect.ValueOf(target)
	// 检测 target 是否符合我们的要求
	// 我们只支持有限的几种输入
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return "", nil, errInvalidEntity
	}

	// 使用 strings.Builder 来拼接 字符串
	bd := strings.Builder{}

	// 构造 INSERT INTO XXX，XXX 是你的表名，这里我们直接用结构体名字
	bd.WriteString(fmt.Sprintf("INSERT INTO `%s`", typeOf.Name()))

	// 遍历所有的字段，构造出来的是 INSERT INTO XXX(col1, col2, col3)
	numField := typeOf.NumField()
	if numField == 0 {
		return "", nil, errInvalidEntity
	}
	// 在这个遍历的过程中，你就可以把参数构造出来
	fields := make([]string, 0, numField)
	args := make([]any, 0, numField)
	for i := 0; i < numField; i++ {
		f := typeOf.Field(i)
		valueField := valueOf.Field(i)
		if f.IsExported() {
			kind := f.Type.Kind()
			switch {
			case kind == reflect.Struct && !isImplValuer(f):
				innerNumField := f.Type.NumField()
				for j := 0; j < innerNumField; j++ {
					fields = append(fields, f.Type.Field(j).Name)
					args = append(args, valueField.Field(j).Interface())
				}
			case kind == reflect.Ptr && !isImplValuer(f):
				f.Type = f.Type.Elem()
				valueField = valueField.Elem()
				if !valueField.IsNil() {
					for k := 0; k < f.Type.NumField(); k++ {
						fields = append(fields, f.Type.Field(k).Name)
						args = append(args, valueField.Field(k).Interface())
					}
				}
			default:
				fields = append(fields, f.Name)
				args = append(args, valueField.Interface())
			}
		}
	}
	bd.WriteString("(")
	for i, f := range fields {
		bd.WriteString(fmt.Sprintf("`%s`", f))
		if i != len(fields)-1 {
			bd.WriteString(",")
		}
	}
	bd.WriteString(") VALUES(")
	for i := range fields {
		bd.WriteString("?")
		if i != len(fields)-1 {
			bd.WriteString(",")
		}
	}
	bd.WriteString(");")

	return bd.String(), args, nil

	// 如果你打算支持组合，那么这里你要深入解析每一个组合的结构体
	// 并且层层深入进去

	// 拼接 VALUES，达成 INSERT INTO XXX(col1, col2, col3) VALUES

	// 再一次遍历所有的字段，要拼接成 INSERT INTO XXX(col1, col2, col3) VALUES(?,?,?)
	// 注意，在第一次遍历的时候我们就已经拿到了参数的值，所以这里就是简单拼接 ?,?,?

	// return bd.String(), args, nil
}

func isImplValuer(f reflect.StructField) bool {
	return f.Type.Implements(reflect.TypeOf((*driver.Valuer)(nil)).Elem())
}

func isImpl(t reflect.Type, target reflect.Type) bool {
	return t.Implements(target)
}
