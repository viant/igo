package plan

import (
	"fmt"
	"github.com/viant/igo/internal"
	"github.com/viant/igo/internal/et"
	"go/ast"
)

func (s *Scope) compileStmt(stmt ast.Stmt) (et.New, error) {
	return s.compileStmtWithLabel(stmt, 0)
}

func (s *Scope) compileStmtWithLabel(stmt ast.Stmt, stmtLabel uint32) (et.New, error) {
	if s.options.StmtListener != nil {
		newStmt, err := s.options.StmtListener(stmt)
		if err != nil {
			return nil, err
		}
		if newStmt != nil {
			stmt = newStmt
		}
	}
	switch actual := stmt.(type) {
	case *ast.LabeledStmt:
		labelID := s.ensureLabelID(actual.Label.Name)
		return s.compileStmtWithLabel(actual.Stmt, labelID)
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
		return s.compileForStmt(actual, stmtLabel)
	case *ast.DeclStmt:
		s.Metric.FlagDeclare()
		return s.compileDeclStmt(actual.Decl)
	case *ast.BranchStmt:
		return s.compileBranchStmt(actual)
	case *ast.RangeStmt:
		return s.compileRangeStmt(actual, stmtLabel)
	case *ast.SwitchStmt:
		return s.compileSwitchStmt(actual, stmtLabel)
	case *ast.TypeSwitchStmt:
		return s.compileTypeSwitchStmt(actual, stmtLabel)
	case *ast.DeferStmt:
		return s.compileDeferStmt(actual)
	case *ast.ExprStmt:
		s.Metric.FlagExpr()
		if callExpr, ok := actual.X.(*ast.CallExpr); ok {
			return s.compileCallExprStmt(callExpr)
		}
		n, _, err := s.compileExpr(actual.X)
		return n, err
	case *ast.GoStmt:
		return nil, fmt.Errorf("go statement is not supported yet")
	case *ast.SelectStmt:
		return nil, fmt.Errorf("select statement is not supported yet")
	case *ast.SendStmt:
		return nil, fmt.Errorf("channel send is not supported yet")
	}
	return nil, fmt.Errorf("unsupported stmt: %T", stmt)
}
