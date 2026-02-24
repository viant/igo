package plan

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope_Compile_InlineHelperInLoop_AdaptShape(t *testing.T) {
	type InRow struct {
		Name  string
		Score int
	}
	type OutRow struct {
		Label string
		Value int
	}

	scope := NewScope()
	_ = scope.RegisterType(reflect.TypeOf(InRow{}))
	_ = scope.RegisterType(reflect.TypeOf(OutRow{}))
	_, _ = scope.DefineVariable("rows", reflect.TypeOf([]*InRow{}))
	_, _ = scope.DefineVariable("needle", reflect.TypeOf(""))
	_, _ = scope.DefineVariable("minScore", reflect.TypeOf(0))
	_, _ = scope.DefineVariable("outRows", reflect.TypeOf([]OutRow{}))
	_, _ = scope.DefineVariable("outSum", reflect.TypeOf(0))

	execution, err := scope.Compile(`
		adapt := func(row *InRow, needle string, minScore int) (OutRow, bool) {
			nameLower := strings.ToLower(row.Name)
			needleLower := strings.ToLower(needle)
			if contains(nameLower, needleLower) && row.Score >= minScore {
				return OutRow{Label: nameLower, Value: row.Score + 1}, true
			}
			return OutRow{}, false
		}

		outRows = make([]OutRow, 0)
		outSum = 0
		for _, row := range rows {
			item, ok := adapt(row, needle, minScore)
			if !ok {
				continue
			}
			outRows = append(outRows, item)
			outSum += item.Value
		}
	`)
	if !assert.NoError(t, err) {
		return
	}

	state := execution.NewState()
	err = state.SetValue("rows", []*InRow{
		{Name: "ALPHA", Score: 10},
		{Name: "beta", Score: 3},
		{Name: "galaxy", Score: 8},
	})
	if !assert.NoError(t, err) {
		return
	}
	err = state.SetString("needle", "al")
	if !assert.NoError(t, err) {
		return
	}
	err = state.SetInt("minScore", 8)
	if !assert.NoError(t, err) {
		return
	}
	execution.Exec(state)

	sum, err := state.Int("outSum")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 20, sum)

	value, err := state.Value("outRows")
	if !assert.NoError(t, err) {
		return
	}
	rowsOut, ok := value.([]OutRow)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, []OutRow{
		{Label: "alpha", Value: 11},
		{Label: "galaxy", Value: 9},
	}, rowsOut)
}

func TestScope_Compile_InlineHelperDeclaredWithVar(t *testing.T) {
	scope := NewScope()
	execution, err := scope.Compile(`
		var plus = func(x int, y int) int {
			return x + y
		}
		z := plus(4, 7)
	`)
	if !assert.NoError(t, err) {
		return
	}
	state := execution.NewState()
	execution.Exec(state)
	z, err := state.Int("z")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 11, z)
}
