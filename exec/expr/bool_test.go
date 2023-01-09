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

	var testCases = []struct {
		description string
		input       interface{}
		expr        string
		expect      bool
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
	}

	for _, testCase := range testCases {
		scope := igo.NewScope()
		input := reflect.ValueOf(testCase.input)
		inputType := input.Type()
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
