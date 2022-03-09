package plan

import (
	"github.com/viant/igo/exec/est"
	"go/ast"
	"go/token"
)

func (s *Scope) compileAssignStmt(assignStmt *ast.AssignStmt) (est.New, error) {
	define := assignStmt.Tok == token.DEFINE
	callExpr, isCallExprAssigment := assignStmt.Rhs[0].(*ast.CallExpr)
	rOperands, err := s.assembleOperands(assignStmt.Rhs, false)
	if err != nil {
		return nil, err
	}
	lOperands, err := s.assembleOperands(assignStmt.Lhs, true)
	if err != nil {
		return nil, err
	}
	if isCallExprAssigment {
		return s.compileCallExprAssign(callExpr, lOperands)
	}

	var group = make([]est.New, len(lOperands))
	for i := 0; i < len(lOperands); i++ {
		if group[i], err = s.compileAssign(lOperands[i], rOperands[i], assignStmt.Tok, define); err != nil {
			return nil, err
		}
	}
	return est.NewGroupStmt(group, false), nil
}

func (s *Scope) compileAssign(left, right *est.Operand, assignToken token.Token, define bool) (est.New, error) {
	switch assignToken {

	case token.ADD_ASSIGN, token.SUB_ASSIGN, token.AND_ASSIGN, token.OR_ASSIGN, token.AND_NOT_ASSIGN, token.MUL_ASSIGN, token.REM_ASSIGN, token.QUO_ASSIGN:
		newFn, _ := est.NewBinaryExpr(assignToken, left, right)
		return newFn, nil
	default:
		if left.Selector.Type == nil {
			s.adjust(left.Selector, right.Type.Type())
			left.Type = right.Type
		}
		return est.NewAssignExpr(left, right), nil
	}
}

func (s *Scope) assembleOperands(expressions []ast.Expr, defined bool) ([]*est.Operand, error) {
	var result = make([]*est.Operand, len(expressions))
	var err error
	for i, expr := range expressions {
		if result[i], err = s.assembleOperand(expr, defined); err != nil {
		}
	}
	return result, nil
}
