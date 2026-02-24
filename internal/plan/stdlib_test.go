package plan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope_Compile_LazyStdlibFunc_String(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		z := strings.ToUpper("abc")
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	z, err := state.String("z")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "ABC", z)
}

func TestScope_Compile_LazyStdlibVariadic(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		z := fmt.Sprintf("%d-%s", 7, "x")
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	z, err := state.String("z")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "7-x", z)
}

func TestScope_Compile_LazyStdlibType(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		var t time.Time
		z := t.IsZero()
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	z, err := state.Bool("z")
	if !assert.Nil(t, err) {
		return
	}
	assert.True(t, z)
}
