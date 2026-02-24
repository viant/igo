package plan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScope_Compile_Switch(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
			r := 0
			x := 2
			switch {
			case true:
				r += 3
				break
				r = 100
			default:
				r = 99
			}
			r += 1
		`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	x, err := state.Int("x")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 2, x)
	actual, err := state.Int("r")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 4, actual)
}

func TestScope_Compile_SwitchNoTag(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		a := 12
		r := 0
		switch {
		case a < 5:
			r = 1
		case a < 10:
			r = 2
		default:
			r = 3
		}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("r")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 3, actual)
}

func TestScope_Compile_Defer(t *testing.T) {
	scope := NewScope()
	var order []int
	scope.RegisterFunc("push", func(v int) {
		order = append(order, v)
	})
	execution, err := scope.Compile(`
		if true {
			defer push(1)
		}
		defer push(2)
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	assert.Equal(t, []int{2, 1}, order)
}

func TestScope_Compile_DeferWithReturn(t *testing.T) {
	scope := NewScope()
	called := 0
	scope.RegisterFunc("mark", func() {
		called++
	})
	execution, err := scope.Compile(`
		defer mark()
		return
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	assert.Equal(t, 1, called)
}
