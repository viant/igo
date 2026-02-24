package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
	"reflect"
)

var interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()

func (s *Scope) parseTypeSwitchAssign(assign ast.Stmt) (ast.Expr, string, error) {
	switch actual := assign.(type) {
	case *ast.ExprStmt:
		assertExpr, ok := actual.X.(*ast.TypeAssertExpr)
		if !ok || assertExpr.Type != nil {
			return nil, "", fmt.Errorf("invalid type switch assign expression")
		}
		return assertExpr.X, "", nil
	case *ast.AssignStmt:
		if len(actual.Rhs) != 1 || len(actual.Lhs) != 1 {
			return nil, "", fmt.Errorf("invalid type switch assign statement")
		}
		if actual.Tok != token.DEFINE && actual.Tok != token.ASSIGN {
			return nil, "", fmt.Errorf("invalid type switch assign token: %v", actual.Tok)
		}
		assertExpr, ok := actual.Rhs[0].(*ast.TypeAssertExpr)
		if !ok || assertExpr.Type != nil {
			return nil, "", fmt.Errorf("invalid type switch assignment expression")
		}
		ident, ok := actual.Lhs[0].(*ast.Ident)
		if !ok {
			return nil, "", fmt.Errorf("unsupported type switch lhs %T", actual.Lhs[0])
		}
		return assertExpr.X, ident.Name, nil
	default:
		return nil, "", fmt.Errorf("unsupported type switch assign %T", assign)
	}
}

func (s *Scope) compileTypeSwitchStmt(typeSwitchStmt *ast.TypeSwitchStmt, labelID uint32) (et.New, error) {
	var statements []et.New
	if typeSwitchStmt.Init != nil {
		initNew, err := s.compileStmt(typeSwitchStmt.Init)
		if err != nil {
			return nil, err
		}
		statements = append(statements, initNew)
	}

	xExpr, assignName, err := s.parseTypeSwitchAssign(typeSwitchStmt.Assign)
	if err != nil {
		return nil, err
	}
	x, err := s.assembleOperand(xExpr, false)
	if err != nil {
		return nil, err
	}

	var defaultBody et.New
	cases := make([]et.TypeSwitchCase, 0, len(typeSwitchStmt.Body.List))
	for _, item := range typeSwitchStmt.Body.List {
		caseClause, ok := item.(*ast.CaseClause)
		if !ok {
			return nil, fmt.Errorf("unsupported type switch clause %T", item)
		}
		scope := s.SubScope()
		var assignSelector *exec.Selector
		if assignName != "" && assignName != "_" {
			assignType := interfaceType
			if len(caseClause.List) == 1 {
				if ident, ok := caseClause.List[0].(*ast.Ident); ok && ident.Name == "nil" {
					assignType = interfaceType
				} else {
					assignType, err = s.discoverType(caseClause.List[0])
					if err != nil {
						return nil, err
					}
				}
			}
			assignSelector, err = scope.DefineVariable(assignName, assignType)
			if err != nil {
				return nil, err
			}
		}

		bodyStatements := make([]et.New, len(caseClause.Body))
		for _, bodyItem := range caseClause.Body {
			if labeled, ok := bodyItem.(*ast.LabeledStmt); ok {
				scope.ensureLabelID(labeled.Label.Name)
			}
		}
		for i := range caseClause.Body {
			if branch, ok := caseClause.Body[i].(*ast.BranchStmt); ok && branch.Tok == token.FALLTHROUGH {
				return nil, fmt.Errorf("cannot fallthrough in type switch")
			}
			bodyStatements[i], err = scope.compileStmt(caseClause.Body[i])
			if err != nil {
				return nil, err
			}
		}
		bodyNew := et.NewControlBlockStmt(bodyStatements)
		if len(caseClause.List) == 0 {
			defaultBody = bodyNew
			continue
		}
		types := make([]reflect.Type, len(caseClause.List))
		for i, tExpr := range caseClause.List {
			if ident, ok := tExpr.(*ast.Ident); ok && ident.Name == "nil" {
				types[i] = nil
				continue
			}
			types[i], err = s.discoverType(tExpr)
			if err != nil {
				return nil, err
			}
		}
		cases = append(cases, et.TypeSwitchCase{
			Types:  types,
			Body:   bodyNew,
			Assign: assignSelector,
		})
	}

	switchNew, err := et.NewTypeSwitchStmt(x, cases, defaultBody, labelID)
	if err != nil {
		return nil, err
	}
	statements = append(statements, switchNew)
	return et.NewGroupStmt(statements, false), nil
}
