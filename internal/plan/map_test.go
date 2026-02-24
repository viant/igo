package plan

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScope_Compile_MapStringKey(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := map[string]int{"a": 10, "b": 20}
	x := m["a"]
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("x")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 10, actual)
}

func TestScope_Compile_MapSetAndGet(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := make(map[string]int, 2)
	m["a"] = 7
	m["b"] = 5
	x := m["a"] + m["b"]
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("x")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 12, actual)
}

func TestScope_Compile_MapValueOk(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := map[string]int{"a": 3}
	v1, ok1 := m["a"]
	v2, ok2 := m["b"]
	x := 0
	if ok1 {
		x += v1
	}
	if !ok2 {
		x += 100
	}
	x += v2
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("x")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 103, actual)
}

func TestScope_Compile_MapIntKey(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := map[int]string{10: "x"}
	s, ok := m[10]
	r := ""
	if ok {
		r = s
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.String("r")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "x", actual)
}

func TestScope_Compile_TypedVarMap(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	var m map[string]int
	m = make(map[string]int)
	m["k"] = 9
	x := m["k"]
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("x")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 9, actual)
}

func TestScope_Compile_MapStringStringAndBool(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m1 := map[string]string{"a": "x"}
	m2 := map[string]bool{"on": true}
	m1["b"] = "y"
	s := m1["a"] + m1["b"]
	_, ok := m2["off"]
	if !ok {
		m2["off"] = false
	}
	r := m2["on"]
	_ = r
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.String("s")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "xy", actual)
}

func TestScope_Compile_MapIntFloat64(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := map[int]float64{1: 2.5}
	m[2] = 1.5
	v1, ok1 := m[1]
	v2, ok2 := m[3]
	sum := 0.0
	if ok1 {
		sum += v1
	}
	if !ok2 {
		sum += 10.0
	}
	sum += v2
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Float64("sum")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 12.5, actual)
}

func TestScope_Compile_MapInt64String(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := make(map[int64]string, 2)
	m[int64(1)] = "a"
	m[int64(2)] = "b"
	s1, ok1 := m[int64(1)]
	s2, ok2 := m[int64(2)]
	out := ""
	if ok1 {
		out += s1
	}
	if ok2 {
		out += s2
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.String("out")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "ab", actual)
}

func TestScope_Compile_MapStringTime(t *testing.T) {
	scope := NewScope()
	scope.RegisterFunc("refTime", func() time.Time {
		return time.Unix(123, 0).UTC()
	})
	execution, err := scope.Compile(`
	m := make(map[string]time.Time, 1)
	m["t"] = refTime()
	_, ok := m["t"]
	x := 0
	if ok {
		x = 1
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("x")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 1, actual)
}

func TestScope_Compile_MapRangeKeyValue(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := map[string]int{"ab": 3, "c": 5}
	sum := 0
	for k, v := range m {
		sum += len(k) + v
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("sum")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 11, actual)
}

func TestScope_Compile_MapRangeKeyOnly(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := map[int]bool{1: true, 3: false}
	sum := 0
	for k := range m {
		sum += k
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Int("sum")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 4, actual)
}

func TestScope_Compile_MapRangeValueOnly(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := map[int64]float64{int64(1): 2.0, int64(2): 4.0}
	sum := 0.0
	for _, v := range m {
		sum += v
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Float64("sum")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 6.0, actual)
}

func TestScope_Compile_MapDelete(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := map[string]int{"a": 1, "b": 2}
	delete(m, "a")
	_, ok := m["a"]
	z := 0
	if !ok {
		z = 1
	}
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
	assert.Equal(t, 1, actual)
}

func TestScope_Compile_MapDeleteStringTime(t *testing.T) {
	scope := NewScope()
	scope.RegisterFunc("refTime", func() time.Time {
		return time.Unix(123, 0).UTC()
	})
	execution, err := scope.Compile(`
	m := make(map[string]time.Time, 1)
	m["t"] = refTime()
	delete(m, "t")
	_, ok := m["t"]
	z := 0
	if !ok {
		z = 1
	}
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
	assert.Equal(t, 1, actual)
}

func TestScope_Compile_MapUint64Float64FastPath(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
	m := make(map[uint64]float64, 2)
	var k1 uint64 = 1
	var k2 uint64 = 2
	m[k1] = 1.5
	m[k2] = 3.0
	v, ok := m[k1]
	sum := 0.0
	if ok {
		sum += v
	}
	for _, item := range m {
		sum += item
	}
	delete(m, k2)
	v2, ok2 := m[k2]
	if ok2 {
		sum = v2
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	actual, err := state.Float64("sum")
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 6.0, actual)
}
