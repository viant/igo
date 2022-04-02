package plan

import (
	"fmt"
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/internal/exec/et"
	"go/ast"
)

func (s *Scope) compileStmt(stmt ast.Stmt) (et.New, error) {

	switch actual := stmt.(type) {
	case *ast.AssignStmt:
		return s.compileAssignStmt(actual)
	case *ast.BlockStmt:
		scope := s.subScope()
		return scope.compileBlockStmt(actual, false)
	case *ast.IncDecStmt:
		return s.compileIncDec(actual)
	case *ast.IfStmt:
		return s.compileIfStmt(actual)
	case *ast.ReturnStmt:
		s.Control.Concat(exec.RtReturn)
		operands, err := s.assembleOperands(actual.Results, false)
		if err != nil {
			return nil, err
		}
		return et.NewReturnStmt(operands, *s.out)
	case *ast.ForStmt:
		return s.compileForStmt(actual)
	case *ast.DeclStmt:
		return s.compileDeclStmt(actual.Decl)
	case *ast.BranchStmt:
		return s.compileBranchStmt(actual)
	case *ast.RangeStmt:
		return s.compileRangeStmt(actual)
	case *ast.ExprStmt:
		n, _, err := s.compileExpr(actual.X)
		return n, err
	}
	return nil, fmt.Errorf("unsupported stmt: %Tgit init", stmt)
}
