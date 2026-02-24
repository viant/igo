package plan

import (
	"fmt"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func (s *Scope) compileSwitchStmt(switchStmt *ast.SwitchStmt, labelID uint32) (et.New, error) {
	var statements []et.New
	if switchStmt.Init != nil {
		initNew, err := s.compileStmt(switchStmt.Init)
		if err != nil {
			return nil, err
		}
		statements = append(statements, initNew)
	}

	tagExpr := switchStmt.Tag
	var intTag *et.Operand
	intSwitchEligible := false
	if tagExpr != nil {
		tagOp, err := s.assembleOperand(tagExpr, false)
		if err == nil && tagOp != nil && tagOp.Type != nil && tagOp.Type.Kind() == reflect.Int {
			intTag = tagOp
			intSwitchEligible = true
		}
	}
	var defaultBranch et.New
	cases := make([]et.SwitchCase, 0, len(switchStmt.Body.List))
	intCases := make([]et.IntSwitchCase, 0, len(switchStmt.Body.List))

	for _, item := range switchStmt.Body.List {
		caseClause, ok := item.(*ast.CaseClause)
		if !ok {
			return nil, fmt.Errorf("unsupported switch clause %T", item)
		}
		scope := s.SubScope()
		bodyItems := caseClause.Body
		hasFallthrough := false
		if len(bodyItems) > 0 {
			if branch, ok := bodyItems[len(bodyItems)-1].(*ast.BranchStmt); ok && branch.Tok == token.FALLTHROUGH {
				hasFallthrough = true
				bodyItems = bodyItems[:len(bodyItems)-1]
				intSwitchEligible = false
			}
		}

		bodyStatements := make([]et.New, len(bodyItems))
		var err error
		for _, bodyItem := range bodyItems {
			if labeled, ok := bodyItem.(*ast.LabeledStmt); ok {
				scope.ensureLabelID(labeled.Label.Name)
			}
		}
		for i := range bodyItems {
			bodyStatements[i], err = scope.compileStmt(bodyItems[i])
			if err != nil {
				return nil, err
			}
		}
		bodyNew := et.NewControlBlockStmt(bodyStatements)
		if len(caseClause.List) == 0 {
			if hasFallthrough {
				return nil, fmt.Errorf("cannot fallthrough from default case")
			}
			defaultBranch = bodyNew
			continue
		}
		if intSwitchEligible {
			values, ok := parseIntCaseValues(caseClause.List)
			if !ok {
				intSwitchEligible = false
			} else {
				intCases = append(intCases, et.IntSwitchCase{
					Values: values,
					Body:   bodyNew,
				})
			}
		}

		condExpr := switchCaseCondExpr(tagExpr, caseClause.List)
		cond, err := s.assembleOperand(condExpr, false)
		if err != nil {
			return nil, err
		}
		cases = append(cases, et.SwitchCase{
			Cond:           cond,
			Body:           bodyNew,
			HasFallthrough: hasFallthrough,
		})
	}

	var switchNew et.New
	var err error
	if intSwitchEligible && intTag != nil {
		switchNew, err = et.NewIntSwitchStmt(intTag, intCases, defaultBranch, labelID)
	} else {
		switchNew, err = et.NewSwitchStmt(cases, defaultBranch, labelID)
	}
	if err != nil {
		return nil, err
	}
	statements = append(statements, switchNew)
	return et.NewGroupStmt(statements, false), nil
}

func parseIntCaseValues(list []ast.Expr) ([]int, bool) {
	values := make([]int, 0, len(list))
	for _, item := range list {
		lit, ok := item.(*ast.BasicLit)
		if !ok || lit.Kind != token.INT {
			return nil, false
		}
		v, err := strconv.Atoi(lit.Value)
		if err != nil {
			return nil, false
		}
		values = append(values, v)
	}
	return values, true
}

func switchCaseCondExpr(tag ast.Expr, list []ast.Expr) ast.Expr {
	if len(list) == 0 {
		return &ast.Ident{Name: "false"}
	}
	var result ast.Expr
	for i, item := range list {
		caseExpr := item
		if tag != nil {
			caseExpr = &ast.BinaryExpr{X: tag, Op: token.EQL, Y: item}
		}
		if i == 0 {
			result = caseExpr
			continue
		}
		result = &ast.BinaryExpr{X: result, Op: token.LOR, Y: caseExpr}
	}
	return result
}
