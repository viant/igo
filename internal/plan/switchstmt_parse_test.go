package plan

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustParseExprList(t *testing.T, src []string) []ast.Expr {
	t.Helper()
	result := make([]ast.Expr, 0, len(src))
	for _, item := range src {
		expr, err := parser.ParseExpr(item)
		if !assert.NoError(t, err) {
			return nil
		}
		result = append(result, expr)
	}
	return result
}

func TestParseIntCaseValues_DataDriven(t *testing.T) {
	testCases := []struct {
		name   string
		exprs  []string
		values []int
		ok     bool
	}{
		{name: "single", exprs: []string{"1"}, values: []int{1}, ok: true},
		{name: "multiple", exprs: []string{"0", "2", "10"}, values: []int{0, 2, 10}, ok: true},
		{name: "non-int-literal", exprs: []string{`"x"`}, ok: false},
		{name: "ident", exprs: []string{"v"}, ok: false},
		{name: "negative-unary", exprs: []string{"-1"}, ok: false},
		{name: "atoi-overflow", exprs: []string{"99999999999999999999999999999999999"}, ok: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := mustParseExprList(t, tc.exprs)
			if list == nil {
				return
			}
			values, ok := parseIntCaseValues(list)
			assert.Equal(t, tc.ok, ok)
			if tc.ok {
				assert.Equal(t, tc.values, values)
			}
		})
	}
}
