package plan

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope_Compile_ComplexStructAssignAndAccess_DataDriven(t *testing.T) {
	type Bar struct {
		ID1 int
	}
	type Foo struct {
		ID   int
		Name string
		Bar  *Bar
	}

	testCases := []struct {
		name       string
		code       string
		resultPath string
		expect     int
	}{
		{
			name: "nested-struct-assign-branch",
			code: `
				n := "abc"
				foo := &Foo{ID:123, Name:n, Bar:&Bar{ID1:456}}
				z := 0
				if foo.Name == "abc" {
					z = foo.ID + foo.Bar.ID1
				}
			`,
			resultPath: "z",
			expect:     579,
		},
		{
			name: "nested-struct-dummy-logic",
			code: `
				n := "xy"
				foo := &Foo{ID:120, Name:n, Bar:&Bar{ID1:30}}
				r := 0
				if foo.Name != "" {
					r = foo.ID + foo.Bar.ID1 + len(foo.Name) + 1
				}
			`,
			resultPath: "r",
			expect:     153,
		},
		{
			name: "slice-of-struct-pointers-with-field-access",
			code: `
				foos := []*Foo{
					&Foo{ID:5, Name:"a", Bar:&Bar{ID1:2}},
					&Foo{ID:7, Name:"bc", Bar:&Bar{ID1:4}},
				}
				sum := 0
				for i, foo := range foos {
					sum += foo.ID + foo.Bar.ID1 + i
					if foo.Name != "" {
						sum += 1
					}
				}
			`,
			resultPath: "sum",
			expect:     21,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scope := NewScope()
			_ = scope.RegisterType(reflect.TypeOf(Foo{}))
			_ = scope.RegisterType(reflect.TypeOf(Bar{}))

			execution, err := scope.Compile(tc.code)
			if !assert.NoError(t, err) {
				return
			}

			state := execution.NewState()
			execution.Exec(state)

			actual, err := state.Int(tc.resultPath)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tc.expect, actual)
		})
	}
}
