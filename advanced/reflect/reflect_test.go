package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateField(t *testing.T) {
	tests := []struct {
		name        string
		in          any
		exceptedRes map[string]any
		exceptedErr error
	}{

		{
			"nil",
			nil,
			nil,
			errors.New("input can not be nil"),
		},
		{
			"string",
			"jackie",
			nil,
			errors.New("not struct kind"),
		},
		{
			"user",
			User{Name: "jackie"},
			map[string]any{"Name": "jackie"},
			nil,
		},
		{
			"user pointer",
			&User{Name: "jackie"},
			map[string]any{"Name": "jackie"},
			nil,
		},
		{
			"unexported field",
			&User{
				Name: "jackie",
				age:  18,
			},
			map[string]any{"Name": "jackie", "age": 0},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := iterateField(test.in)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.exceptedRes, res)
		})
	}
}

func TestSetField(t *testing.T) {
	tests := []struct {
		name        string
		target      *User
		fieldName   string
		val         any
		want        *User
		exceptedErr error
	}{
		{

			"set name",
			&User{
				Name: "jackie",
				age:  18,
			},
			"Name",
			"Mickey",
			&User{
				Name: "Mickey",
				age:  18,
			},
			nil,
		},
		{

			"nil target",
			nil,
			"Name",
			"Mickey",
			nil,
			targetNilErr,
		},
		{

			"set age",
			&User{
				Name: "mickey",
				age:  18,
			},
			"age",
			1,
			&User{
				Name: "mickey",
				age:  18,
			},
			nil,
		},

		{
			"set name to nil",
			&User{
				Name: "mickey",
				age:  18,
			},
			"Name",
			nil,
			&User{
				Name: "",
				age:  18,
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := setField(test.target, test.fieldName, test.val)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.want.Name, test.target.Name)
			assert.Equal(t, test.want.age, test.target.age)
		})
	}
}

func TestGetFieldVal(t *testing.T) {
	tests := []struct {
		name        string
		target      any
		fieldName   string
		exceptedVal any
		exceptedErr error
	}{
		{
			"target is nil",
			nil,
			"Name",
			nil,
			targetNilErr,
		},
		{
			"filed not found",
			&User{
				Name: "mickey",
			},
			"name",
			nil,
			fieldNotFoundErr,
		},
		{
			"normal",
			&User{
				Name: "mickey",
			},
			"Name",
			"mickey",
			nil,
		},
		{
			"target is struct",
			User{
				Name: "mickey",
			},
			"Name",
			"mickey",
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := getFieldVal(test.target, test.fieldName)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.exceptedVal, val)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		in          any
		exceptedVal any
		exceptedErr error
	}{
		{
			"nil",
			nil,
			nil,
			targetNilErr,
		},
		{
			"int",
			int(0),
			0,
			nil,
		},
		{
			"string",
			"",
			"",
			nil,
		},
		{
			"user",
			User{},
			User{},
			nil,
		},
		{
			"bool",
			true,
			false,
			nil,
		},
		{
			"pointer",
			&User{},
			User{},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := New(test.in)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.exceptedVal, val)
		})
	}
}

func TestIterateMethodName(t *testing.T) {
	tests := []struct {
		name        string
		target      any
		exceptedVal []string
		exceptedErr error
	}{
		{
			"pointer",
			&User{},
			[]string{"GetName", "Hello"},
			nil,
		},
		{
			"struct",
			User{}, //结构体只能拿到传递值的方法
			[]string{"GetName"},
			nil,
		},
		{
			"nil target",
			nil,
			nil,
			targetNilErr,
		},
		{
			"no struct",
			"hello",
			nil,
			invalidTypeErr,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			names, err := iterateMethodName(test.target)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.exceptedVal, names)
		})
	}
}

func TestCallMethod(t *testing.T) {
	tests := []struct {
		name        string
		target      any
		methodName  string
		params      []any
		exceptedVal []any
		exceptedErr error
	}{
		{
			"call Hello",
			&User{Name: "mickey"},
			"Hello",
			[]any{"mickey"},
			[]any{"Hello,mickey"},
			nil,
		},
		{
			"call GetName",
			&User{Name: "mickey"},
			"GetName",
			[]any{"mickey"},
			[]any{"mickey"},
			invalidMethodParamErr,
		},
		{
			"invalid method params",
			&User{Name: "mickey"},
			"GetName",
			nil,
			[]any{"mickey"},
			nil,
		},

		{
			"struct can not call pointer method",
			User{Name: "mickey"},
			"Hello",
			nil,
			nil,
			methodNotFoundErr,
		},
		{
			"nil target",
			nil,
			"",
			nil,
			nil,
			targetNilErr,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := CallMethod(test.target, test.methodName, test.params)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.exceptedVal, result)
		})
	}
}

func TestIterateRangeAble(t *testing.T) {
	tests := []struct {
		name        string
		target      any
		exceptedVal []any
		exceptedErr error
	}{
		{
			"iterate string slice",
			[]string{"mickey", "jackie"},
			[]any{"mickey", "jackie"},
			nil,
		},
		{
			"iterate int slice",
			[]int{101, 201},
			[]any{101, 201},
			nil,
		},
		{
			"iterate int array",
			[2]int{101, 102},
			[]any{101, 102},
			nil,
		},
		{
			"iterate map value",
			map[string]int{"a": 1, "b": 2},
			[]any{1, 2},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := IterateRangeAble(test.target)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, len(test.exceptedVal), len(result))
			for i := 0; i < len(result); i++ {
				assert.Equal(t, test.exceptedVal[i], result[i])
			}
		})
	}
}
