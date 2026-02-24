package plan

import (
	"github.com/viant/igo/internal/et"
	"go/ast"
)

func (s *Scope) compileDeferStmt(deferStmt *ast.DeferStmt) (et.New, error) {
	newFn, err := s.compileCallExprStmt(deferStmt.Call)
	if err != nil {
		return nil, err
	}
	*s.defers = append(*s.defers, newFn)
	return et.NewNop(), nil
}
