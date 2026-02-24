package plan

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestScope_Compile_TypeAssertExpr(t *testing.T) {
	scope := NewScope()
	scope.DefineVariable("x", reflect.TypeOf((*interface{})(nil)).Elem())
	execution, err := scope.Compile(`
	v := x.(int)
	z := v + 1
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	state.SetValue("x", 3)
	execution.Exec(state)
	actual, err := state.Int("z")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 4, actual)
}

func TestScope_Compile_TypeAssertValueOk(t *testing.T) {
	scope := NewScope()
	scope.DefineVariable("x", reflect.TypeOf((*interface{})(nil)).Elem())
	execution, err := scope.Compile(`
	v, ok := x.(string)
	z := 0
	if ok {
		z = len(v)
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	state.SetValue("x", 3)
	execution.Exec(state)
	actual, err := state.Int("z")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 0, actual)
}

func TestScope_Compile_TypeSwitch(t *testing.T) {
	scope := NewScope()
	scope.DefineVariable("x", reflect.TypeOf((*interface{})(nil)).Elem())
	execution, err := scope.Compile(`
	z := 0
	switch x.(type) {
	case int:
		z = 1
	case string:
		z = 2
	default:
		z = 3
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	state.SetValue("x", "abc")
	execution.Exec(state)
	actual, err := state.Int("z")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 2, actual)
}

func TestScope_Compile_TypeSwitchAssign(t *testing.T) {
	scope := NewScope()
	scope.DefineVariable("x", reflect.TypeOf((*interface{})(nil)).Elem())
	execution, err := scope.Compile(`
	z := 0
	switch v := x.(type) {
	case int:
		z = v + 10
	case string:
		z = len(v)
	default:
		z = 7
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	state.SetValue("x", 5)
	execution.Exec(state)
	actual, err := state.Int("z")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 15, actual)
}

func TestScope_Compile_TypeSwitchCaseNilAndMultiType(t *testing.T) {
	scope := NewScope()
	scope.DefineVariable("x", reflect.TypeOf((*interface{})(nil)).Elem())
	execution, err := scope.Compile(`
	z := 0
	switch x.(type) {
	case nil:
		z = 11
	case int, string:
		z = 22
	default:
		z = 33
	}
	`)
	if !assert.Nil(t, err) {
		return
	}

	{
		state := execution.NewState()
		execution.Exec(state)
		actual, err := state.Int("z")
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, 11, actual)
	}
	{
		state := execution.NewState()
		state.SetValue("x", "abc")
		execution.Exec(state)
		actual, err := state.Int("z")
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, 22, actual)
	}
}
