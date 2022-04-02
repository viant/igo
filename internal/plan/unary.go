package plan

import (
	"fmt"
	"github.com/viant/igo/internal/exec/et"
	"go/ast"
	"reflect"
)

func (s *Scope) compileUnaryExpr(unaryExpr *ast.UnaryExpr) (et.New, reflect.Type, error) {
	switch expr := unaryExpr.X.(type) {
	case *ast.CompositeLit: //pointer to a type
		return s.compileCompositeLiteral(unaryExpr.Op, expr)
	}
	return nil, nil, fmt.Errorf("unsupported type: %T", unaryExpr.X)
}
