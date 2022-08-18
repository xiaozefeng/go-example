package reflect

import (
	"errors"
	"reflect"
)

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

var (
	targetNilErr     = errors.New("target must not be nil")
	invalidTypeErr   = errors.New("target must be Ptr")
	notFoundFieldErr = errors.New("field not found")
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
		return notFoundFieldErr
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
