package plan

import (
	"fmt"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

func literalOperand(node *ast.BasicLit) (*et.Operand, error) {
	var t reflect.Type
	var sV string
	var iV int
	var fV float64
	var cV complex128

	switch node.Kind {
	case token.INT:
		t = reflect.TypeOf(0)
		iV, _ = strconv.Atoi(node.Value)
		return et.NewOperand(nil, t, nil, iV), nil
	case token.FLOAT:
		t = reflect.TypeOf(0.0)
		fV, _ = strconv.ParseFloat(node.Value, 64)
		return et.NewOperand(nil, t, nil, fV), nil
	case token.STRING:
		sV = strings.Trim(node.Value, "\"`")
		t = reflect.TypeOf("")
		return et.NewOperand(nil, t, nil, sV), nil
	case token.CHAR:
		unquoted, err := strconv.Unquote(node.Value)
		if err != nil {
			return nil, err
		}
		runes := []rune(unquoted)
		if len(runes) != 1 {
			return nil, fmt.Errorf("invalid rune literal: %s", node.Value)
		}
		t = reflect.TypeOf(rune(0))
		return et.NewOperand(nil, t, nil, runes[0]), nil
	case token.IMAG:
		parsed, err := strconv.ParseComplex(node.Value, 128)
		if err != nil {
			return nil, err
		}
		cV = parsed
		t = reflect.TypeOf(complex128(0))
		return et.NewOperand(nil, t, nil, cV), nil
	default:
		return nil, fmt.Errorf("unsupported token: %v", node.Kind)
	}
}
