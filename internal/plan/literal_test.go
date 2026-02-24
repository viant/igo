package plan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScope_Compile_RuneLiteral(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	r := 'A'
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	value, err := state.Value("r")
	if !assert.Nil(t, err) {
		return
	}
	actual, ok := value.(rune)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, rune('A'), actual)
}

func TestScope_Compile_ImagLiteral(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	z := 2i
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	value, err := state.Value("z")
	if !assert.Nil(t, err) {
		return
	}
	actual, ok := value.(complex128)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, complex(0, 2), actual)
}

func TestScope_Compile_ImagLiteralFromCall(t *testing.T) {
	scope := NewScope()
	scope.RegisterFunc("imag2", func() complex128 {
		return 2i
	})
	execution, err := scope.Compile(`
	z := imag2()
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	value, err := state.Value("z")
	if !assert.Nil(t, err) {
		return
	}
	actual, ok := value.(complex128)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, complex(0, 2), actual)
}
