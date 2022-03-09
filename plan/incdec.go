package plan

import (
	"github.com/viant/igo/exec/est"
	"go/ast"
)

func (s *Scope) compileIncDec(incDecStmt *ast.IncDecStmt) (est.New, error) {
	op, err := s.assembleOperand(incDecStmt.X, false)
	if err != nil {
		return nil, err
	}
	return est.NewIncDec(incDecStmt.Tok, op), nil

}
