package plan

import (
	"fmt"
	"github.com/viant/igo/internal/et"
	"github.com/viant/xunsafe"
	"go/ast"
	"go/token"
	"reflect"
)

func (s *Scope) compileTypeAssertExpr(assertExpr *ast.TypeAssertExpr) (et.New, reflect.Type, error) {
	if assertExpr.Type == nil {
		return nil, nil, fmt.Errorf("type assertion type was empty")
	}
	x, err := s.assembleOperand(assertExpr.X, false)
	if err != nil {
		return nil, nil, err
	}
	targetType, err := s.discoverType(assertExpr.Type)
	if err != nil {
		return nil, nil, err
	}
	return et.NewTypeAssert(x, targetType)
}

func (s *Scope) compileTypeAssertAssignStmt(assignStmt *ast.AssignStmt) (et.New, bool, error) {
	if len(assignStmt.Lhs) != 2 || len(assignStmt.Rhs) != 1 {
		return nil, false, nil
	}
	assertExpr, ok := assignStmt.Rhs[0].(*ast.TypeAssertExpr)
	if !ok || assertExpr.Type == nil {
		return nil, false, nil
	}
	if assignStmt.Tok != token.ASSIGN && assignStmt.Tok != token.DEFINE {
		return nil, true, fmt.Errorf("unsupported type assertion assignment token: %v", assignStmt.Tok)
	}
	targetType, err := s.discoverType(assertExpr.Type)
	if err != nil {
		return nil, true, err
	}
	x, err := s.assembleOperand(assertExpr.X, false)
	if err != nil {
		return nil, true, err
	}
	lOperands, err := s.assembleOperands(assignStmt.Lhs, assignStmt.Tok == token.DEFINE)
	if err != nil {
		return nil, true, err
	}
	if lOperands[0].Selector.Type == nil {
		_ = s.adjust(lOperands[0].Selector, targetType)
		lOperands[0].Type = xunsafe.NewType(targetType)
	}
	if lOperands[1].Selector.Type == nil {
		_ = s.adjust(lOperands[1].Selector, boolType)
		lOperands[1].Type = xunsafe.NewType(boolType)
	}
	if lOperands[0].Type == nil {
		lOperands[0].Type = xunsafe.NewType(lOperands[0].Selector.Type)
	}
	if lOperands[1].Type == nil {
		lOperands[1].Type = xunsafe.NewType(lOperands[1].Selector.Type)
	}
	newFn, err := et.NewTypeAssert2Assign(x, targetType, lOperands[0], lOperands[1])
	return newFn, true, err
}
