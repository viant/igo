package plan

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope_Function_SliceShapeAdapter_DataDriven(t *testing.T) {
	type SourceRow struct {
		Name   string
		Score  int
		Active bool
	}
	type TargetRow struct {
		Label string
		Value int
	}

	testCases := []struct {
		name   string
		input  []SourceRow
		needle string
		min    int
		expect []TargetRow
	}{
		{
			name: "filter-and-adapt-two",
			input: []SourceRow{
				{Name: "Alpha", Score: 10, Active: true},
				{Name: "beta", Score: 3, Active: true},
				{Name: "galaxy", Score: 8, Active: true},
				{Name: "metal", Score: 6, Active: false},
			},
			needle: "al",
			min:    7,
			expect: []TargetRow{
				{Label: "alpha", Value: 20},
				{Label: "galaxy", Value: 16},
			},
		},
		{
			name: "no-matches",
			input: []SourceRow{
				{Name: "echo", Score: 2, Active: true},
				{Name: "foxtrot", Score: 9, Active: false},
			},
			needle: "zz",
			min:    1,
			expect: []TargetRow{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scope := NewScope()
			_ = scope.RegisterType(reflect.TypeOf(SourceRow{}))
			_ = scope.RegisterType(reflect.TypeOf(TargetRow{}))
			// ensure external stdlib function use in expected validation is aligned
			scope.RegisterFunc("toLower", strings.ToLower)

			fn, err := scope.Function(`func(rows []SourceRow, needle string, minScore int) []TargetRow {
				out := make([]TargetRow, 0)
				needleLower := strings.ToLower(needle)
				for _, row := range rows {
					if !row.Active {
						continue
					}
					nameLower := strings.ToLower(row.Name)
					if contains(nameLower, needleLower) && row.Score >= minScore {
						out = append(out, TargetRow{
							Label: nameLower,
							Value: row.Score * 2,
						})
					}
				}
				return out
			}`)
			if !assert.NoError(t, err) {
				return
			}

			transform, ok := fn.(func([]SourceRow, string, int) []TargetRow)
			if !assert.True(t, ok) {
				return
			}

			actual := transform(tc.input, tc.needle, tc.min)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestScope_Compile_RuntimeSliceShapeAdapter(t *testing.T) {
	type SourceRow struct {
		Name   string
		Score  int
		Active bool
	}
	type TargetRow struct {
		Label string
		Value int
	}

	scope := NewScope()
	_ = scope.RegisterType(reflect.TypeOf(SourceRow{}))
	_ = scope.RegisterType(reflect.TypeOf(TargetRow{}))
	_, _ = scope.DefineVariable("inRows", reflect.TypeOf([]SourceRow{}))
	_, _ = scope.DefineVariable("needle", reflect.TypeOf(""))
	_, _ = scope.DefineVariable("minScore", reflect.TypeOf(0))
	_, _ = scope.DefineVariable("outRows", reflect.TypeOf([]TargetRow{}))
	_, _ = scope.DefineVariable("outSum", reflect.TypeOf(0))

	execution, err := scope.Compile(`
		outRows = make([]TargetRow, 0)
		outSum = 0
		needleLower := strings.ToLower(needle)
		for _, row := range inRows {
			if !row.Active {
				continue
			}
			nameLower := strings.ToLower(row.Name)
			if contains(nameLower, needleLower) && row.Score >= minScore {
				item := TargetRow{Label: nameLower, Value: row.Score + 1}
				outRows = append(outRows, item)
				outSum += item.Value
			}
		}
	`)
	if !assert.NoError(t, err) {
		return
	}

	state := execution.NewState()
	err = state.SetValue("inRows", []SourceRow{
		{Name: "ALPHA", Score: 10, Active: true},
		{Name: "beta", Score: 3, Active: true},
		{Name: "galaxy", Score: 8, Active: true},
		{Name: "metal", Score: 7, Active: false},
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

	value, err := state.Value("outRows")
	if !assert.NoError(t, err) {
		return
	}
	rows, ok := value.([]TargetRow)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, []TargetRow{
		{Label: "alpha", Value: 11},
		{Label: "galaxy", Value: 9},
	}, rows)

	sum, err := state.Int("outSum")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 20, sum)
}
