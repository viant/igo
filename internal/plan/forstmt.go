package plan

import (
	"github.com/viant/igo/internal/et"
	"go/ast"
)

func (s *Scope) compileForStmt(forStmt *ast.ForStmt) (et.New, error) {
	var err error
	scope := s.SubScope()
	var cond *et.Operand
	var init, post, body et.New
	if forStmt.Init != nil {
		if init, err = s.compileStmt(forStmt.Init); err != nil {
			return nil, err
		}
	}
	if forStmt.Post != nil {
		if post, err = scope.compileStmt(forStmt.Post); err != nil {
			return nil, err
		}
	}
	if forStmt.Cond != nil {
		if cond, err = s.assembleOperand(forStmt.Cond, false); err != nil {
			return nil, err
		}
	}
	if body, err = scope.compileBlockStmt(forStmt.Body, true); err != nil {
		return nil, err
	}
	return et.NewForStmt(cond, init, post, body)
}
