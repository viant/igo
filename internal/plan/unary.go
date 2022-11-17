package plan

import (
	"fmt"
	"github.com/viant/igo/internal"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"reflect"
	"unsafe"
)

func (s *Scope) compileUnaryExpr(unaryExpr *ast.UnaryExpr) (et.New, reflect.Type, error) {
	switch expr := unaryExpr.X.(type) {
	case *ast.CompositeLit: //pointer to a type
		return s.compileCompositeLiteral(unaryExpr.Op, expr)
	case *ast.ParenExpr:
		switch unaryExpr.Op.String() {
		case "!":
			newEst, _, err := s.compileExpr(unaryExpr.X)
			if err != nil {
				return nil, nil, err
			}
			return func(control *et.Control) (internal.Compute, error) {
				compute, err := newEst(control)
				if err != nil {
					return nil, err
				}
				return func(ptr unsafe.Pointer) unsafe.Pointer {
					vPtr := compute(ptr)
					out := *(*bool)(vPtr)
					*(*bool)(vPtr) = !out
					return vPtr
				}, nil
			}, boolType, nil
		case "":
			return s.compileExpr(unaryExpr.X)
		}
	}
	return nil, nil, fmt.Errorf("unsupported type: %T", unaryExpr.X)
}
