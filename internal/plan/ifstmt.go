package plan

import (
	"github.com/viant/igo/internal/exec/et"
	"go/ast"
)

func (s *Scope) compileIfStmt(ifStmt *ast.IfStmt) (et.New, error) {
	cond, err := s.assembleOperand(ifStmt.Cond, false)
	if err != nil {
		return nil, err
	}
	scope := s.subScope()
	var whenBranch, elseBranch et.New
	whenBranch, err = scope.compileBlockStmt(ifStmt.Body, false)
	if err != nil {
		return nil, err
	}
	if ifStmt.Else != nil {
		switch actual := ifStmt.Else.(type) {
		case *ast.BlockStmt:
			elseWhenScope := s.subScope()
			if elseBranch, err = elseWhenScope.compileBlockStmt(actual, false); err != nil {
				return nil, err
			}
		case *ast.IfStmt:
			if elseBranch, err = s.compileIfStmt(actual); err != nil {
				return nil, err
			}
		}
	}
	return et.NewIfStmt(cond, whenBranch, elseBranch)
}
