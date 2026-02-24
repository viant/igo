package plan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScope_Compile_SwitchFallthrough(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	x := 0
	switch 1 {
	case 1:
		x = 1
		fallthrough
	case 2:
		x = x + 2
	default:
		x = 99
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	value, err := state.Int("x")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 3, value)
}

func TestScope_Compile_GotoForwardLabel(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	x := 0
	goto L
	x = 1
L:
	x = x + 2
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	value, err := state.Int("x")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 2, value)
}

func TestScope_Compile_LabeledBreakAndContinue(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	sumBreak := 0
OuterBreak:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sumBreak++
			if j == 1 {
				break OuterBreak
			}
		}
	}

	sumContinue := 0
OuterContinue:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				continue OuterContinue
			}
			sumContinue++
		}
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	sumBreak, err := state.Int("sumBreak")
	if !assert.Nil(t, err) {
		return
	}
	sumContinue, err := state.Int("sumContinue")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 2, sumBreak)
	assert.Equal(t, 3, sumContinue)
}

func TestScope_Compile_SpecializedIntForLoopControl(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		sum := 0
		for i := 0; i < 10; i++ {
			if i == 3 {
				continue
			}
			if i == 8 {
				break
			}
			sum += i
		}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	sum, err := state.Int("sum")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 25, sum)
}

func TestScope_Compile_SpecializedIntForLoopStepAndSymmetricCond(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		sum := 0
		for i := 0; 10 > i; i += 2 {
			sum += i
		}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	sum, err := state.Int("sum")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 20, sum)
}

func TestScope_Compile_SpecializedIntForLoopAssignPost(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		sum := 0
		for i := 1; i <= 9; i = i + 2 {
			sum += i
		}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	sum, err := state.Int("sum")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 25, sum)
}

func TestScope_Compile_SpecializedIntForLoopIntCall(t *testing.T) {
	scope := NewScope()
	total := 0
	scope.RegisterFunc("sinkInt", func(v int) {
		total += v
	})
	execution, err := scope.Compile(`
		for i := 0; i < 10; i++ {
			sinkInt(i)
		}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	assert.Equal(t, 45, total)
}

func TestScope_Compile_SpecializedIntForLoopGuardedIntCall(t *testing.T) {
	scope := NewScope()
	total := 0
	scope.RegisterFunc("sinkInt", func(v int) {
		total += v
	})
	execution, err := scope.Compile(`
		for i := 0; i < 10; i++ {
			if i % 3 != 0 {
				sinkInt(i)
			}
		}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	assert.Equal(t, 27, total)
}

func TestScope_Compile_SpecializedIntForLoopAccPatterns(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		a := 0
		b := 0
		c := 1
		d := 1000
		for i := 1; i <= 5; i++ {
			a -= i
		}
		for i := 0; i < 5; i++ {
			b += 3
		}
		for i := 1; i <= 4; i++ {
			c *= i
		}
		for i := 2; i <= 10; i += 2 {
			d /= 2
		}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	a, err := state.Int("a")
	if !assert.Nil(t, err) {
		return
	}
	b, err := state.Int("b")
	if !assert.Nil(t, err) {
		return
	}
	c, err := state.Int("c")
	if !assert.Nil(t, err) {
		return
	}
	d, err := state.Int("d")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, -15, a)
	assert.Equal(t, 15, b)
	assert.Equal(t, 24, c)
	assert.Equal(t, 31, d)
}

func TestScope_Compile_SpecializedIntForLoopGuardedAcc(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		sum := 0
		for i := 0; i < 10; i++ {
			if i < 5 {
				sum += i
			}
		}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	sum, err := state.Int("sum")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 10, sum)
}
