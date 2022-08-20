package accessor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Name string
}

func Test(t *testing.T) {
	tests := []struct {
		name        string
		target      any
		filedName   string
		exceptedVal any
		exceptedErr error
	}{
		{
			"get name",
			&User{Name: "mickey"},
			"Name",
			"mickey",
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			accessor, err := NewAccessor(test.target)
			if err != nil {
				return
			}
			val, err := accessor.GetField(test.filedName)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.exceptedVal, val)
		})
	}
}

func TestAccessor_SetField(t *testing.T) {
	tests := []struct {
		name        string
		target      *User
		fieldName   string
		val         any
		exceptedErr error
	}{
		{
			"set name",
			&User{},
			"Name",
			"mickey",
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			accessor, err := NewAccessor(test.target)
			if err != nil {
				return
			}
			err = accessor.SetField(test.fieldName, test.val)
			assert.Equal(t, test.exceptedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, test.val, test.target.Name)
		})

	}
}

func BenchmarkAccessor_GetField(b *testing.B) {
	accessor, err := NewAccessor(&User{Name: "mickey"})
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := accessor.GetField("Name")
		if err != nil {
			b.Error(err)
		}

	}
}

func BenchmarkAccessor_SetField(b *testing.B) {
	accessor, err := NewAccessor(&User{})
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := accessor.SetField("Name", "mickey")
		if err != nil {
			b.Error(err)
		}
	}
}
