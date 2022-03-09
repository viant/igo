package plan

import (
	"github.com/viant/igo/exec/est"
	"go/ast"
)

func (s *Scope) compileBlockStmt(blockStmt *ast.BlockStmt, isForStmt bool) (est.New, error) {
	var statements = make([]est.New, len(blockStmt.List))
	var err error
	for i := range blockStmt.List {
		if statements[i], err = s.compileStmt(blockStmt.List[i]); err != nil {
			return nil, err
		}
	}
	return est.NewBlockStmt(statements, isForStmt), nil
}
