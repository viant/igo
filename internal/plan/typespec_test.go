package plan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScope_Compile_VarInterfaceDecl(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	var x interface{}
	x = 10
	y := x.(int) + 1
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("y")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 11, actual)
}

func TestScope_Compile_TypeSpecAliasAndInterface(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	type Any interface{}
	type Scores map[string]int
	type Label = string

	var x Any
	x = 7

	s := make(Scores, 1)
	s["a"] = x.(int)

	var l Label
	l = "ok"
	z := s["a"] + len(l)
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("z")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 9, actual)
}
