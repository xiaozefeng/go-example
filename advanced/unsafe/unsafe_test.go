package unsafe

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccessor_GetIntField(t *testing.T) {
	tests := []struct {
		name        string
		target      any
		fieldName   string
		exceptedVal any
		exceptedErr error
	}{
		{
			"get int field",
			&User{
				Age: 18,
			},
			"Age",
			18,
			nil,
		},
		{
			"nil target",
			nil,
			"Age",
			18,
			TargetNilErr,
		},
		{
			"filed not foud",
			&User{Age: 18},
			"age",
			18,
			FieldNoFoundErr,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			accessor, err := NewAccessor(test.target)
			if err != nil {
				return
			}
			field, err := accessor.GetIntField(test.fieldName)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.exceptedVal, field)
		})
	}
}

func TestAccessor_GetStringField(t *testing.T) {
	tests := []struct {
		name        string
		target      any
		fieldName   string
		exceptedVal any
		exceptedErr error
	}{
		{
			"get string field",
			&User{
				Name: "mickey",
			},
			"Name",
			"mickey",
			nil,
		},
		{
			"nil target",
			nil,
			"Name",
			"mickey",
			TargetNilErr,
		},
		{
			"filed not found",
			&User{Name: "mickey"},
			"name",
			"mickey",
			FieldNoFoundErr,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			accessor, err := NewAccessor(test.target)
			if err != nil {
				return
			}
			field, err := accessor.GetStringField(test.fieldName)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.exceptedVal, field)
		})
	}
}

func TestAccessor_SetInt32Field(t *testing.T) {
	tests := []struct {
		name        string
		target      *User
		filedName   string
		val         int32
		exceptedErr error
	}{
		{
			"set age",
			&User{},
			"Age",
			18,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			accessor, err := NewAccessor(test.target)
			if err != nil {
				return
			}
			err = accessor.SetInt32Field(test.filedName, test.val)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.EqualValues(t, test.val, test.target.Age)
		})
	}
}

func BenchmarkUnsafeGetField(b *testing.B) {
	accessor, err := NewAccessor(&User{Name: "mickey"})
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := accessor.GetIntField("Name")
		if err != nil {
			b.Error(err)
		}

	}
}

func BenchmarkUnsafeSetFiled(b *testing.B) {
	accessor, err := NewAccessor(&User{})
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := accessor.SetStringField("Name", "mickey")
		if err != nil {
			b.Error(err)
		}
	}
}
