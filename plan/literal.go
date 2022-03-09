package plan

import (
	"fmt"
	"github.com/viant/igo/exec/est"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

func literalOperand(node *ast.BasicLit) (*est.Operand, error) {
	var t reflect.Type
	var sV string
	var iV int
	var fV float64

	switch node.Kind {
	case token.INT:
		t = reflect.TypeOf(0)
		iV, _ = strconv.Atoi(node.Value)
		return est.NewOperand(nil, t, nil, iV), nil
	case token.FLOAT:
		t = reflect.TypeOf(0.0)
		fV, _ = strconv.ParseFloat(node.Value, 64)
		return est.NewOperand(nil, t, nil, fV), nil
	case token.STRING:
		sV = strings.Trim(node.Value, "\"`")
		t = reflect.TypeOf("")
		return est.NewOperand(nil, t, nil, sV), nil
	default:
		return nil, fmt.Errorf("unsupported token: %v", node.Kind)
	}
}
