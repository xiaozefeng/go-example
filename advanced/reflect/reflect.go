package reflect

import (
	"errors"
	"reflect"
)

type User struct {
	Name string
	age  int
}

func (u *User) Hello(name string) string {
	return "Hello," + name
}
func (u User) GetName() string {
	return u.Name
}
func iterateField(v any) (map[string]any, error) {
	if v == nil {
		return nil, errors.New("input can not be nil")
	}
	val := reflect.ValueOf(v)

	t := reflect.TypeOf(v)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("not struct kind")
	}
	m := make(map[string]any)
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		fieldVal := val.Field(i)
		if field.IsExported() {
			m[field.Name] = fieldVal.Interface()
		} else {
			m[field.Name] = reflect.Zero(field.Type).Interface()
		}
	}
	return m, nil
}

type FuncInfo struct {
	Name   string
	In     []reflect.Type
	Out    []reflect.Type
	Result []any
}

func iterateMethods(target any) (map[string]*FuncInfo, error) {
	if target == nil {
		return nil, targetNilErr
	}
	typeOf := reflect.TypeOf(target)
	valueOf := reflect.ValueOf(target)
	if valueOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			return nil, targetNilErr
		}
	}
	//for typeOf.Kind() == reflect.Ptr {
	//	typeOf = typeOf.Elem()
	//	valueOf = valueOf.Elem()
	//}
	if typeOf.Kind() != reflect.Struct && typeOf.Kind() != reflect.Ptr {
		return nil, invalidTypeErr
	}
	numField := typeOf.NumField()
	result := make(map[string]*FuncInfo, numField)
	for i := 0; i < numField; i++ {
		m := typeOf.Method(i)
		numIn := m.Type.NumIn()
		params := make([]reflect.Value, 0, numIn)
		params = append(params, reflect.ValueOf(target))
		in := make([]reflect.Type, 0, numIn)
		for j := 0; j < numIn; j++ {
			p := m.Type.In(j)
			in = append(in, p)
			if j > 0 {
				params = append(params, reflect.Zero(p))
			}
		}

		ret := m.Func.Call(params)
		outNum := m.Type.NumOut()
		out := make([]reflect.Type, 0, outNum)
		res := make([]any, 0, outNum)
		for k := 0; k < outNum; k++ {
			out = append(out, m.Type.Out(k))
			res = append(res, ret[k].Interface())
		}

		result[m.Name] = &FuncInfo{
			Name:   m.Name,
			In:     in,
			Out:    out,
			Result: res,
		}
	}
	return result, nil
}

var (
	targetNilErr          = errors.New("target must not be nil")
	invalidTypeErr        = errors.New("target must be Ptr")
	fieldNotFoundErr      = errors.New("field not found")
	methodNotFoundErr     = errors.New("method not found")
	invalidMethodParamErr = errors.New("invalid method param")
)

func setField(target any, fieldName string, val any) error {
	v := reflect.ValueOf(target)
	t := reflect.TypeOf(target)

	if t.Kind() != reflect.Ptr {
		return invalidTypeErr
	}
	if v.Kind() == reflect.Ptr {
		isNil := v.IsNil()
		if isNil {
			return targetNilErr
		}
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	typeField, found := t.FieldByName(fieldName)
	if !found {
		return fieldNotFoundErr
	}
	filed := v.FieldByName(fieldName)
	if filed.CanSet() {
		if val == nil {
			filed.Set(reflect.Zero(typeField.Type))
		} else {
			filed.Set(reflect.ValueOf(val))
		}
	}
	return nil
}

func getFieldVal(target any, fieldName string) (any, error) {
	if target == nil {
		return nil, targetNilErr
	}
	typeOf := reflect.TypeOf(target)
	valueOf := reflect.ValueOf(target)
	if valueOf.Kind() == reflect.Ptr {
		isNil := valueOf.IsNil()
		if isNil {
			return nil, targetNilErr
		}
	}

	for typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	_, found := typeOf.FieldByName(fieldName)
	if !found {
		return nil, fieldNotFoundErr
	}
	return valueOf.FieldByName(fieldName).Interface(), nil
}

func New(target any) (any, error) {
	if target == nil {
		return nil, targetNilErr
	}
	typeOf := reflect.TypeOf(target)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	ptr := reflect.New(typeOf)
	return ptr.Elem().Interface(), nil
}

func iterateMethodName(target any) ([]string, error) {
	if target == nil {
		return nil, targetNilErr
	}
	typeOf := reflect.TypeOf(target)
	valueOf := reflect.ValueOf(target)
	if valueOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			return nil, targetNilErr
		}
	}
	if typeOf.Kind() != reflect.Struct && typeOf.Kind() != reflect.Ptr {
		return nil, invalidTypeErr
	}
	numMethod := typeOf.NumMethod()
	result := make([]string, 0, numMethod)
	for i := 0; i < numMethod; i++ {
		method := typeOf.Method(i)
		result = append(result, method.Name)
	}
	return result, nil
}

func CallMethod(target any, methodName string, params []any) ([]any, error) {
	if target == nil {
		return nil, targetNilErr
	}
	typeOf := reflect.TypeOf(target)
	valueOf := reflect.ValueOf(target)
	if valueOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			return nil, targetNilErr
		}
	}
	method, found := typeOf.MethodByName(methodName)
	if !found {
		return nil, methodNotFoundErr
	}
	numIn := method.Type.NumIn() - 1

	if numIn != len(params) {
		return nil, invalidMethodParamErr
	}
	in := make([]reflect.Value, 0)
	in = append(in, valueOf) //第一个参数是调用者自己
	for i := 0; i < numIn; i++ {
		in = append(in, reflect.ValueOf(params[i]))
	}
	ret := method.Func.Call(in)
	numOut := method.Type.NumOut()
	result := make([]any, 0, numOut)
	for i := 0; i < numOut; i++ {
		result = append(result, ret[i].Interface())
	}
	return result, nil
}

func IterateRangeAble(rangeAble any) ([]any, error) {
	if rangeAble == nil {
		return nil, targetNilErr
	}
	typeOf := reflect.TypeOf(rangeAble)
	valueOf := reflect.ValueOf(rangeAble)
	if valueOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			return nil, targetNilErr
		}
	}
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	kind := typeOf.Kind()
	switch {
	case kind == reflect.Array || kind == reflect.Slice:
		result := make([]any, 0, valueOf.Len())
		for i := 0; i < valueOf.Len(); i++ {
			result = append(result, valueOf.Index(i).Interface())
		}
		return result, nil
	case kind == reflect.Map:
		result := make([]any, 0, len(valueOf.MapKeys()))
		mapRange := valueOf.MapRange()
		for mapRange.Next() {
			result = append(result, mapRange.Value().Interface())
		}
		return result, nil
	default:
		return nil, invalidTypeErr
	}

}
