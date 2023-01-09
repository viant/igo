package plan

import (
	"fmt"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"reflect"
)

func (s *Scope) compileUnaryStarExpr(z *ast.StarExpr) (et.New, reflect.Type, error) {
	op, err := s.assembleOperand(z.X, true)
	if err != nil {
		return nil, nil, err
	}
	if op.Selector == nil {
		return nil, nil, fmt.Errorf("invalid *expr, %v not a pointer", op.Name)
	}

	exprType := op.Selector.Type
	if exprType.Kind() != reflect.Ptr {
		return nil, nil, fmt.Errorf("invalid *expr, %v not a pointer", op.Name)
	}
	newFn, err := et.NewStarExpr(op)
	return newFn, exprType.Elem(), nil
}
