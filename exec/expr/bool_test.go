package expr_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/igo"
	"reflect"
	"testing"
)

func TestNewBool(t *testing.T) {

	type case001 struct {
		V1, V2 int
	}

	type Case002 struct {
		IPtr *int
		I    int
		F    float64
		FPtr *float64
	}

	type case003 struct {
		list   []int
		active *bool
	}

	type case004 struct {
		active *bool
		id     int
	}

	type case005 struct {
		active *bool
		id     string
	}

	var testCases = []struct {
		description  string
		functionName string
		function     interface{}
		input        interface{}
		expr         string
		expect       bool
	}{
		{
			description: "primitive test - when true",
			expr: `if v.V1 > v.V2 {
	return true
}
	return false
`,
			input:  &case001{V1: 5, V2: 3},
			expect: true,
		},
		{
			description: "primitive test - when false",
			expr: `if v.V1 > v.V2 {
	return true
}
	return false
`,
			input:  &case001{V1: 5, V2: 7},
			expect: false,
		},

		{
			description: "pointer test - when not true",
			expr: `if v.IPtr != nil && *v.IPtr > 4 {
	return true
}
	return false
`,
			input:  &Case002{},
			expect: false,
		},

		{
			description: "pointer test - when  true",
			expr: `if v.IPtr != nil && *v.IPtr > 4 {
				return true
		}
			return false
		`,
			input:  &Case002{IPtr: intPtr(5)},
			expect: true,
		},

		{
			description: "len - when  true",
			expr: `if len(v.list) == 2 &&  *v.active {
				return true
		}
			return false
		`,
			input:  &case003{list: []int{1, 2}, active: boolPtr(true)},
			expect: true,
		},

		{
			description: "with udf",
			expr: `if *v.active &&  in(v.id, 1,2,3) {
				return true
		}
			return false
		`,
			input:        &case004{id: 1, active: boolPtr(true)},
			functionName: "in",
			function: func(id int, values []int) bool {
				for _, v := range values {
					if v == id {
						return true
					}
				}
				return false
			},
			expect: true,
		},
		{
			description: "with udf",
			expr: `if *v.active &&  in(v.id, "1","2","3") {
				return true
		}
			return false
		`,
			input:        &case005{id: "1", active: boolPtr(true)},
			functionName: "in",
			function: func(id string, values []string) bool {
				for _, v := range values {
					if v == id {
						return true
					}
				}
				return false
			},
			expect: true,
		},
	}

	for _, testCase := range testCases {
		scope := igo.NewScope()
		input := reflect.ValueOf(testCase.input)
		inputType := input.Type()
		if testCase.functionName != "" {
			scope.RegisterFunc(testCase.functionName, testCase.function)
		}
		_, err := scope.DefineVariable("v", inputType)
		assert.Nil(t, err, testCase.description)
		predicate, err := scope.BoolExpression(testCase.expr)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		state := predicate.NewState()
		_ = state.SetValue("v", testCase.input)
		actual := predicate.ComputeWithState(state)
		assert.Equal(t, testCase.expect, actual, testCase.description)
	}

}

func intPtr(i int) *int {
	return &i
}

func boolPtr(i bool) *bool {
	return &i
}
