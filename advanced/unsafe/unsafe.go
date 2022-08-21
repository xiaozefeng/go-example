package unsafe

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type User struct {
	Sex  bool
	Age  int32
	Name string
}

type FieldMeta struct {
	offset uintptr
	typ    reflect.Type
}

type Accessor struct {
	fields     map[string]FieldMeta
	targetAddr unsafe.Pointer
}

func (a *Accessor) GetIntField(fieldName string) (int, error) {
	ptr, _, err := a.getFieldPoint(fieldName)
	if err != nil {
		return 0, err
	}
	return *(*int)(ptr), nil
}
func (a *Accessor) GetField(fieldName string) (any, error) {
	ptr, f, err := a.getFieldPoint(fieldName)
	if err != nil {
		return 0, err
	}
	result := reflect.NewAt(f.typ, ptr)
	return result.Interface(), nil
}

func (a *Accessor) getFieldPoint(fieldName string) (unsafe.Pointer, FieldMeta, error) {
	f, ok := a.fields[fieldName]
	if !ok {
		return nil, f, FieldNoFoundErr
	}
	ptr := unsafe.Pointer(uintptr(a.targetAddr) + f.offset)
	if ptr == nil {
		return nil, f, ErrInvalidAddress(fieldName)
	}
	return ptr, f, nil
}

func (a *Accessor) GetStringField(fieldName string) (string, error) {
	ptr, _, err := a.getFieldPoint(fieldName)
	if err != nil {
		return "", err
	}
	return *(*string)(ptr), nil
}

func (a *Accessor) SetInt32Field(fieldName string, val int32) error {
	ptr, _, err := a.getFieldPoint(fieldName)
	if err != nil {
		return err
	}
	*(*int32)(ptr) = val
	return nil
}

func (a *Accessor) SetStringField(fieldName string, val string) error {
	ptr, _, err := a.getFieldPoint(fieldName)
	if err != nil {
		return err
	}
	*(*string)(ptr) = val
	return nil
}

func (a *Accessor) SetField(fieldName string, val any) error {
	ptr, f, err := a.getFieldPoint(fieldName)
	if err != nil {
		return err
	}
	res := reflect.NewAt(f.typ, ptr)
	if res.CanSet() {
		res.Set(reflect.ValueOf(val))
	}
	return nil
}

var (
	TargetNilErr    = errors.New("target must not be nil")
	InvalidTypeErr  = errors.New("type must be pointer or struct")
	FieldNoFoundErr = errors.New("field not found")
)

func ErrInvalidAddress(fieldName string) error {
	return fmt.Errorf("invalid address, field:%s", fieldName)
}

func NewAccessor(target any) (*Accessor, error) {
	if target == nil {
		return nil, TargetNilErr
	}
	typeOf := reflect.TypeOf(target)
	valueOf := reflect.ValueOf(target)
	if valueOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			return nil, TargetNilErr
		}
	}
	if typeOf.Kind() != reflect.Ptr || typeOf.Elem().Kind() != reflect.Struct {
		return nil, InvalidTypeErr
	}
	typeOf = typeOf.Elem()

	numField := typeOf.NumField()
	fields := make(map[string]FieldMeta, numField)
	for i := 0; i < numField; i++ {
		f := typeOf.Field(i)
		fields[f.Name] = FieldMeta{offset: f.Offset, typ: f.Type}
	}
	return &Accessor{
		fields:     fields,
		targetAddr: valueOf.UnsafePointer(),
	}, nil
}
