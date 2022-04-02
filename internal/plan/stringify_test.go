package plan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringifyExpr(t *testing.T) {
	var testCases = []struct {
		description string
		expr        string
		expect      string
		depth       int
	}{
		{
			description: "array composite expr",
			expr:        "x[y].ID",
		},
		{
			description: "arithmetic expr",
			expr:        "x+y",
		},
		{
			description: "array expr",
			expr:        "x[y]",
		},

		{
			description: "selector array expr",
			expr:        "foo.x[y]",
		},
		{
			description: "func  expr",
			expr:        "foo(x,y)",
		},
		{
			description: "map expr",
			expr:        `m["a"].id`,
		},
		{
			description: "map expr with dep",
			expr:        `m["a"].id`,
			depth:       1,
			expect:      "m",
		},
		{
			description: "map expr with dep",
			expr:        `m["a"].id`,
			depth:       2,
			expect:      `m["a"]`,
		},
		{
			description: "map expr with dep",
			expr:        `m["a"].id`,
			depth:       3,
			expect:      `m["a"].id`,
		},
		{
			description: "anonymous idx expr",
			expr:        `_[1]`,
		},
		{
			description: "anonymous idx expr",
			expr:        `x+y`,
		},
	}

	for _, testCase := range testCases {
		if testCase.expect == "" {
			testCase.expect = testCase.expr
		}
		expr, err := expression(testCase.expr)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		actual := stringifyExpr(expr, testCase.depth)
		assert.Equal(t, testCase.expect, actual)
	}

}
