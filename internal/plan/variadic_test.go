package plan

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestScope_Compile_VariadicInterfaceReturns(t *testing.T) {
	scope := NewScope()
	scope.RegisterFunc("count", func(v ...interface{}) (int, error) {
		return len(v), nil
	})
	scope.RegisterFunc("countS", func(_ string, v ...interface{}) (int, error) {
		return len(v), nil
	})

	execution, err := scope.Compile(`
		c1, _ := count(1)
		c2, _ := count(1, 2, 3, 4, 5)
		c3, _ := countS("x", 1, 2, 3, 4, 5)
	`)
	if !assert.Nil(t, err) {
		return
	}

	state := execution.NewState()
	execution.Exec(state)

	c1, err := state.Int("c1")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 1, c1)

	c2, err := state.Int("c2")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 5, c2)

	c3, err := state.Int("c3")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 5, c3)
}

func TestScope_Compile_VariadicStringAndInterfaceSignatures(t *testing.T) {
	scope := NewScope()
	scope.DefineVariable("out", reflect.TypeOf((*interface{})(nil)).Elem())
	scope.RegisterFunc("join", func(prefix string, v ...interface{}) string {
		result := prefix
		for _, item := range v {
			result += "-" + asString(item)
		}
		return result
	})
	var withBaseSeen string
	scope.RegisterFunc("withBaseSink", func(base interface{}, v ...interface{}) interface{} {
		if len(v) == 0 {
			withBaseSeen = asString(base)
			return base
		}
		withBaseSeen = asString(base) + ":" + asString(v[0])
		return withBaseSeen
	})
	count := 0
	scope.RegisterFunc("sink", func(prefix string, v ...interface{}) {
		count += len(prefix) + len(v)
	})

	execution, err := scope.Compile(`
		a := join("p", 1, "x", true)
		out = withBaseSink("root", 10)
		sink("z", 1, 2, 3)
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)

	a, err := state.String("a")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "p-1-x-true", a)

	assert.Equal(t, "root:10", withBaseSeen)
	assert.Equal(t, 4, count)
}
