package plan

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope_Compile_FilterRowsByContains_DataDriven(t *testing.T) {
	type Row struct {
		Name  string
		Score int
	}

	testCases := []struct {
		name        string
		needle      string
		expectCount int
		expectSum   int
	}{
		{
			name:        "matches-two-rows",
			needle:      "al",
			expectCount: 2,
			expectSum:   17,
		},
		{
			name:        "matches-one-row",
			needle:      "gam",
			expectCount: 1,
			expectSum:   5,
		},
		{
			name:        "matches-none",
			needle:      "zzz",
			expectCount: 0,
			expectSum:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scope := NewScope()
			_ = scope.RegisterType(reflect.TypeOf(Row{}))
			_, err := scope.DefineVariable("needle", reflect.TypeOf(""))
			if !assert.NoError(t, err) {
				return
			}

			execution, err := scope.Compile(`
				rows := []*Row{
					&Row{Name: "alpha", Score: 10},
					&Row{Name: "beta", Score: 3},
					&Row{Name: "galaxy", Score: 7},
					&Row{Name: "gamma", Score: 5},
				}
				filtered := make([]*Row, 0)
				for _, row := range rows {
					if contains(row.Name, needle) {
						filtered = append(filtered, row)
					}
				}

				filterCount := len(filtered)
				filterSum := 0
				for _, row := range filtered {
					filterSum += row.Score
				}
			`)
			if !assert.NoError(t, err) {
				return
			}

			state := execution.NewState()
			err = state.SetString("needle", tc.needle)
			if !assert.NoError(t, err) {
				return
			}
			execution.Exec(state)

			count, err := state.Int("filterCount")
			if !assert.NoError(t, err) {
				return
			}
			sum, err := state.Int("filterSum")
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tc.expectCount, count)
			assert.Equal(t, tc.expectSum, sum)
		})
	}
}

func TestScope_Compile_FilterRuntimeRowsByContains_DataDriven(t *testing.T) {
	type Row struct {
		Name  string
		Score int
	}

	testCases := []struct {
		name        string
		rows        []*Row
		needle      string
		expectCount int
		expectSum   int
	}{
		{
			name: "runtime-rows-two-matches",
			rows: []*Row{
				{Name: "delta", Score: 4},
				{Name: "alphabet", Score: 8},
				{Name: "galaxy", Score: 6},
				{Name: "metal", Score: 5},
			},
			needle:      "al",
			expectCount: 3,
			expectSum:   19,
		},
		{
			name: "runtime-rows-one-match",
			rows: []*Row{
				{Name: "echo", Score: 3},
				{Name: "foxtrot", Score: 7},
				{Name: "golf", Score: 2},
			},
			needle:      "fox",
			expectCount: 1,
			expectSum:   7,
		},
		{
			name: "runtime-rows-no-match",
			rows: []*Row{
				{Name: "iris", Score: 1},
				{Name: "juno", Score: 2},
			},
			needle:      "zzz",
			expectCount: 0,
			expectSum:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scope := NewScope()
			_ = scope.RegisterType(reflect.TypeOf(Row{}))
			_, err := scope.DefineVariable("rows", reflect.TypeOf([]*Row{}))
			if !assert.NoError(t, err) {
				return
			}
			_, err = scope.DefineVariable("needle", reflect.TypeOf(""))
			if !assert.NoError(t, err) {
				return
			}

			execution, err := scope.Compile(`
				filtered := make([]*Row, 0)
				for _, row := range rows {
					if contains(row.Name, needle) {
						filtered = append(filtered, row)
					}
				}

				resultCount := len(filtered)
				resultSum := 0
				for _, row := range filtered {
					resultSum += row.Score
				}
			`)
			if !assert.NoError(t, err) {
				return
			}

			state := execution.NewState()
			err = state.SetValue("rows", tc.rows)
			if !assert.NoError(t, err) {
				return
			}
			err = state.SetString("needle", tc.needle)
			if !assert.NoError(t, err) {
				return
			}
			execution.Exec(state)

			count, err := state.Int("resultCount")
			if !assert.NoError(t, err) {
				return
			}
			sum, err := state.Int("resultSum")
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tc.expectCount, count)
			assert.Equal(t, tc.expectSum, sum)
		})
	}
}

func TestScope_Compile_FilterRuntimeRowsAdvanced_DataDriven(t *testing.T) {
	type Row struct {
		Name  string
		Score int
	}

	testCases := []struct {
		name        string
		rows        []*Row
		needle      string
		minScore    int
		expectCount int
		expectSum   int
	}{
		{
			name: "case-insensitive-multi-condition",
			rows: []*Row{
				{Name: "ALPHA", Score: 10},
				{Name: "beta", Score: 4},
				{Name: "Galaxy", Score: 9},
				{Name: "metal", Score: 5},
			},
			needle:      "AL",
			minScore:    6,
			expectCount: 2,
			expectSum:   19,
		},
		{
			name: "nil-and-empty-rows",
			rows: []*Row{
				nil,
				{Name: "", Score: 100},
				{Name: "omega", Score: 3},
				{Name: "Omicron", Score: 8},
			},
			needle:      "om",
			minScore:    5,
			expectCount: 1,
			expectSum:   8,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scope := NewScope()
			_ = scope.RegisterType(reflect.TypeOf(Row{}))
			_, err := scope.DefineVariable("rows", reflect.TypeOf([]*Row{}))
			if !assert.NoError(t, err) {
				return
			}
			_, err = scope.DefineVariable("needle", reflect.TypeOf(""))
			if !assert.NoError(t, err) {
				return
			}
			_, err = scope.DefineVariable("minScore", reflect.TypeOf(0))
			if !assert.NoError(t, err) {
				return
			}

			execution, err := scope.Compile(`
				filtered := make([]*Row, 0)
				needleLower := strings.ToLower(needle)
				for _, row := range rows {
					if row == nil {
						continue
					}
					if row.Name == "" {
						continue
					}
					if contains(strings.ToLower(row.Name), needleLower) && row.Score >= minScore {
						filtered = append(filtered, row)
					}
				}

				resultCount := len(filtered)
				resultSum := 0
				for _, row := range filtered {
					resultSum += row.Score
				}
			`)
			if !assert.NoError(t, err) {
				return
			}

			state := execution.NewState()
			err = state.SetValue("rows", tc.rows)
			if !assert.NoError(t, err) {
				return
			}
			err = state.SetString("needle", tc.needle)
			if !assert.NoError(t, err) {
				return
			}
			err = state.SetInt("minScore", tc.minScore)
			if !assert.NoError(t, err) {
				return
			}
			execution.Exec(state)

			count, err := state.Int("resultCount")
			if !assert.NoError(t, err) {
				return
			}
			sum, err := state.Int("resultSum")
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tc.expectCount, count)
			assert.Equal(t, tc.expectSum, sum)
		})
	}
}
