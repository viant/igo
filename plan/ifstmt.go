package plan

import (
	"github.com/viant/igo/exec/est"
	"go/ast"
)

func (s *Scope) compileIfStmt(ifStmt *ast.IfStmt) (est.New, error) {
	cond, err := s.assembleOperand(ifStmt.Cond, false)
	if err != nil {
		return nil, err
	}
	scope := s.subScope()
	var whenBranch, elseBranch est.New
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
	return est.NewIfStmt(cond, whenBranch, elseBranch)
}
