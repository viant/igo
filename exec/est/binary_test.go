package est

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/igo/exec"
	"go/token"
	"reflect"
	"testing"
)

type binaryTestCase struct {
	description string
	left        interface{}
	right       interface{}
	expect      interface{}
	unsupported bool
	token       token.Token
}

func (c binaryTestCase) init() (*exec.Variables, exec.Compute, error) {
	x := reflect.StructField{
		Name: "X",
		Type: reflect.TypeOf(c.left),
	}
	y := reflect.StructField{
		Name: "Y",
		Type: reflect.TypeOf(c.right),
	}
	if y.Type.ConvertibleTo(errType) {
		y.Type = errType
	}
	if x.Type.ConvertibleTo(errType) {
		x.Type = errType
	}
	destType := reflect.StructOf([]reflect.StructField{
		x, y,
	})
	varsNew := exec.VariablesNew(destType, nil)
	vars := varsNew()
	xSel, _ := vars.Selector("X")
	ySel, _ := vars.Selector("Y")
	nwBinary, _ := NewBinaryExpr(c.token, NewOperand(xSel, nil, nil, nil), NewOperand(ySel, nil, nil, nil))
	binaryExpr, err := nwBinary(&Control{})
	if err != nil {
		return vars, nil, err
	}
	ptr := vars.Pointer()
	xSel.SetValue(ptr, c.left)
	ySel.SetValue(ptr, c.right)
	return vars, binaryExpr, nil
}

func TestNewBinaryExpr(t *testing.T) {

	var testCases = []binaryTestCase{
		{
			description: "1 + 10",
			left:        1,
			token:       token.ADD,
			right:       10,
			expect:      11,
		},
		{
			description: "1 - 11",
			left:        1,
			token:       token.SUB,
			right:       11,
			expect:      -10,
		},
		{
			description: "2 * 10",
			left:        2,
			token:       token.MUL,
			right:       10,
			expect:      20,
		},
		{
			description: "2 > 10",
			left:        2,
			token:       token.GTR,
			right:       10,
			expect:      false,
		},
		{
			description: "10 > 2",
			left:        10,
			token:       token.GTR,
			right:       2,
			expect:      true,
		},
		{
			description: "2 > 10",
			left:        2,
			token:       token.LSS,
			right:       10,
			expect:      true,
		},
		{
			description: "10 > 2",
			left:        10,
			token:       token.LSS,
			right:       2,
			expect:      false,
		},
		{
			description: "2 >= 10",
			left:        2,
			token:       token.GEQ,
			right:       10,
			expect:      false,
		},
		{
			description: "10 >= 2",
			left:        10,
			token:       token.GEQ,
			right:       2,
			expect:      true,
		},

		{
			description: "2 <= 10",
			left:        2,
			token:       token.LEQ,
			right:       10,
			expect:      true,
		},
		{
			description: "10 <= 2",
			left:        10,
			token:       token.LEQ,
			right:       2,
			expect:      false,
		},
		{
			description: "1 << 2",
			left:        1,
			token:       token.SHL,
			right:       2,
			expect:      4,
		},
		{
			description: "4 >> 2",
			left:        4,
			token:       token.SHR,
			right:       2,
			expect:      1,
		},
		{
			description: "true && false",
			left:        true,
			token:       token.AND,
			right:       false,
			expect:      false,
		},
		{
			description: "true || false",
			left:        true,
			token:       token.OR,
			right:       false,
			expect:      true,
		},
		{
			description: "2 == 2",
			left:        2,
			token:       token.EQL,
			right:       2,
			expect:      true,
		},
		{
			description: "2 == 3",
			left:        2,
			token:       token.EQL,
			right:       3,
			expect:      false,
		},
		{
			description: `"2" == "2"`,
			left:        "2",
			token:       token.NEQ,
			right:       "2",
			expect:      false,
		},
		{
			description: "2.0 == 3.0",
			left:        2.0,
			token:       token.NEQ,
			right:       3.0,
			expect:      true,
		},
	}

	for _, testCase := range testCases {
		scope, expr, err := testCase.init()
		if testCase.unsupported {
			assert.NotNil(t, err, testCase.description)
			continue
		}
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		result := expr(scope.Pointer())
		var actual interface{}
		switch testCase.expect.(type) {
		case int:
			actual = *(*int)(result)
		case string:
			actual = *(*string)(result)
		case bool:
			actual = *(*bool)(result)
		case uint8:
			actual = *(*uint8)(result)
		case float64:
			actual = *(*float64)(result)
		case float32:
			actual = *(*float32)(result)
		}
		assert.Equal(t, testCase.expect, actual, testCase.description)
	}

}
