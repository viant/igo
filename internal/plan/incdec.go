package plan

import (
	"github.com/viant/igo/internal/et"
	"go/ast"
)

func (s *Scope) compileIncDec(incDecStmt *ast.IncDecStmt) (et.New, error) {
	op, err := s.assembleOperand(incDecStmt.X, false)
	if err != nil {
		return nil, err
	}
	return et.NewIncDec(incDecStmt.Tok, op), nil

}
