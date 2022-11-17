package plan

import (
	"fmt"
	"github.com/viant/igo/internal"
	"github.com/viant/igo/internal/et"
	"go/ast"
)

func (s *Scope) compileStmt(stmt ast.Stmt) (et.New, error) {
	if s.stmtListener != nil {
		newStmt, err := s.stmtListener(stmt)
		if err != nil {
			return nil, err
		}
		if newStmt != nil {
			stmt = newStmt
		}
	}
	switch actual := stmt.(type) {
	case *ast.AssignStmt:
		s.Metric.FlagAssign()
		return s.compileAssignStmt(actual)
	case *ast.BlockStmt:
		scope := s.SubScope()
		return scope.compileBlockStmt(actual, false)
	case *ast.IncDecStmt:
		s.Metric.FlagAssign()
		return s.compileIncDec(actual)
	case *ast.IfStmt:
		s.Metric.FlagIfElse()
		return s.compileIfStmt(actual)
	case *ast.ReturnStmt:
		s.Control.Concat(internal.RtReturn)
		operands, err := s.assembleOperands(actual.Results, false)
		if err != nil {
			return nil, err
		}
		return et.NewReturnStmt(operands, *s.out)
	case *ast.ForStmt:
		s.Metric.FlagFor()
		return s.compileForStmt(actual)
	case *ast.DeclStmt:
		s.Metric.FlagDeclare()
		return s.compileDeclStmt(actual.Decl)
	case *ast.BranchStmt:
		return s.compileBranchStmt(actual)
	case *ast.RangeStmt:
		return s.compileRangeStmt(actual)
	case *ast.ExprStmt:
		s.Metric.FlagExpr()
		n, _, err := s.compileExpr(actual.X)
		return n, err
	}
	return nil, fmt.Errorf("unsupported stmt: %Tgit init", stmt)
}
