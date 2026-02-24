package plan

import (
	"github.com/viant/igo/internal/et"
	"go/ast"
)

func (s *Scope) compileBlockStmt(blockStmt *ast.BlockStmt, isForStmt bool) (et.New, error) {
	var statements = make([]et.New, len(blockStmt.List))
	labelPositions := map[uint32]int{}
	for i := range blockStmt.List {
		if labeled, ok := blockStmt.List[i].(*ast.LabeledStmt); ok {
			labelPositions[s.ensureLabelID(labeled.Label.Name)] = i
		}
	}
	var err error
	for i := range blockStmt.List {
		stmt := blockStmt.List[i]
		if labeled, ok := stmt.(*ast.LabeledStmt); ok {
			labelID := s.ensureLabelID(labeled.Label.Name)
			if statements[i], err = s.compileStmtWithLabel(labeled.Stmt, labelID); err != nil {
				return nil, err
			}
			continue
		}
		if statements[i], err = s.compileStmt(stmt); err != nil {
			return nil, err
		}
	}
	if len(labelPositions) > 0 {
		return et.NewJumpBlockStmt(statements, labelPositions, isForStmt), nil
	}
	return et.NewBlockStmt(statements, isForStmt), nil
}
