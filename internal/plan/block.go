package plan

import (
	"github.com/viant/igo/internal/et"
	"go/ast"
)

func (s *Scope) compileBlockStmt(blockStmt *ast.BlockStmt, isForStmt bool) (et.New, error) {
	var statements = make([]et.New, len(blockStmt.List))
	var err error
	for i := range blockStmt.List {
		if statements[i], err = s.compileStmt(blockStmt.List[i]); err != nil {
			return nil, err
		}
	}
	return et.NewBlockStmt(statements, isForStmt), nil
}
