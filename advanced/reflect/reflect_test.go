package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	Name string
	age  int
}

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
