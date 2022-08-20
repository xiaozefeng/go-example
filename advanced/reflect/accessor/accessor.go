package accessor

import (
	"errors"
	"reflect"
)

type FieldMeta struct {
	index int
}

type Accessor struct {
	fields map[string]FieldMeta
	target reflect.Value
}

var (
	targetNilErr          = errors.New("target must not be nil")
	invalidTypeErr        = errors.New("target must be Ptr")
	fieldNotFoundErr      = errors.New("field not found")
	methodNotFoundErr     = errors.New("method not found")
	invalidMethodParamErr = errors.New("invalid method param")
	canNotSetErr          = errors.New("can not set filed")
)

func (a *Accessor) GetField(filedName string) (any, error) {
	f, ok := a.fields[filedName]
	if !ok {
		return nil, fieldNotFoundErr
	}
	return a.target.Field(f.index).Interface(), nil
}

func (a *Accessor) SetField(filedName string, val any) error {
	f, ok := a.fields[filedName]
	if !ok {
		return fieldNotFoundErr
	}
	field := a.target.Field(f.index)
	if !field.CanSet() {
		return canNotSetErr
	}
	field.Set(reflect.ValueOf(val))
	return nil
}

func NewAccessor(v any) (*Accessor, error) {
	if v == nil {
		return nil, errors.New("input can not be nil")
	}
	valueOf := reflect.ValueOf(v)
	typeOf := reflect.TypeOf(v)

	for typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
		valueOf = valueOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return nil, errors.New("not struct kind")
	}
	fields := make(map[string]FieldMeta)
	numField := typeOf.NumField()
	for i := 0; i < numField; i++ {
		f := typeOf.Field(i)
		fields[f.Name] = FieldMeta{index: i}
	}
	return &Accessor{
		fields: fields,
		target: valueOf,
	}, nil
}
