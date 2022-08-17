package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Order struct {
	orderNo     string
	userId      int
	orderStatus int
}

func main() {
	var order = Order{
		orderNo:     "123",
		userId:      1,
		orderStatus: 1,
	}
	sql, err := genInsertSQL(&order)
	if err != nil {
		panic(err)
	}
	fmt.Println(sql)

}

// default mysql
func genInsertSQL(entity interface{}) (string, error) {
	if entity == nil {
		return "", errors.New("must not nil input")
	}
	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", errors.New("not struct")
	}
	var table = t.Name()
	var columns []string
	var values []string

	for i := 0; i < v.NumField(); i++ {
		columns = append(columns, Camel2Case(t.Field(i).Name))
		switch v.Field(i).Kind() {
		case reflect.String:
			values = append(values, fmt.Sprintf("'%s'", v.Field(i).String()))
		case reflect.Int:
			values = append(values, fmt.Sprintf("%d", v.Field(i).Int()))
		}
	}
	var insert = Insert{
		table:   table,
		columns: columns,
		values:  values,
	}
	return insert.Build()
}

type Insert struct {
	table   string
	columns []string
	values  []string
}

func (i *Insert) Build() (string, error) {
	columns := lo.Map[string, string](i.columns, func(x string, _ int) string {
		return Camel2Case(x)
	})
	columnSQL := strings.Join(columns, `,`)
	valuesSQL := strings.Join(i.values, `,`)

	return fmt.Sprintf("insert %s(%s) values(%s)", Camel2Case(i.table), columnSQL, valuesSQL), nil
}

func Camel2Case(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

type Buffer struct {
	*bytes.Buffer
}

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = cases.Title(language.English).String(name)
	return strings.Replace(name, " ", "", -1)
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		_, _ = b.Write(val)
	case rune:
		_, _ = b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	_, _ = b.WriteString(s)
	return b
}
