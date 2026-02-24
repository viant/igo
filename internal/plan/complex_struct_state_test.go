package plan

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope_Compile_StateStructCopyNestedPointerAccess(t *testing.T) {
	type Bar struct {
		ID1 int
	}
	type Foo struct {
		ID   int
		Name string
		Bar  *Bar
	}

	scope := NewScope()
	_ = scope.RegisterType(reflect.TypeOf(Foo{}))
	_ = scope.RegisterType(reflect.TypeOf(Bar{}))
	_, err := scope.DefineVariable("input", reflect.TypeOf(Foo{}))
	if !assert.NoError(t, err) {
		return
	}
	execution, err := scope.Compile(`
		work := input
		out := 0
		if work.Bar != nil {
			out += work.Bar.ID1
		}
		if work.ID > work.Bar.ID1 {
			out += work.ID - work.Bar.ID1
		}
	`)
	if !assert.NoError(t, err) {
		return
	}

	state := execution.NewState()
	err = state.SetValue("input", Foo{
		ID:   10,
		Name: "abc",
		Bar:  &Bar{ID1: 4},
	})
	if !assert.NoError(t, err) {
		return
	}
	execution.Exec(state)
	actual, err := state.Int("out")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 10, actual)
}

func TestScope_Compile_StateStructCopyNestedPointerNilGuard(t *testing.T) {
	type Bar struct {
		ID1 int
	}
	type Foo struct {
		ID  int
		Bar *Bar
	}

	scope := NewScope()
	_ = scope.RegisterType(reflect.TypeOf(Foo{}))
	_ = scope.RegisterType(reflect.TypeOf(Bar{}))
	_, err := scope.DefineVariable("input", reflect.TypeOf(Foo{}))
	if !assert.NoError(t, err) {
		return
	}
	execution, err := scope.Compile(`
		work := input
		out := 0
		if work.Bar != nil {
			out += work.Bar.ID1
		}
	`)
	if !assert.NoError(t, err) {
		return
	}

	state := execution.NewState()
	err = state.SetValue("input", Foo{ID: 3, Bar: nil})
	if !assert.NoError(t, err) {
		return
	}
	execution.Exec(state)
	actual, err := state.Int("out")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 0, actual)
}
