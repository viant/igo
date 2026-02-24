package plan

import (
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
)

type intCallGuardPattern struct {
	kind    int
	op      token.Token
	constV  int
	modBase int
}

func (s *Scope) parseForBodyIntCallGuard(body *ast.BlockStmt, loopVar string) (func(int), intCallGuardPattern, bool) {
	if body == nil || len(body.List) != 1 {
		return nil, intCallGuardPattern{}, false
	}
	ifStmt, ok := body.List[0].(*ast.IfStmt)
	if !ok || ifStmt.Init != nil || ifStmt.Else != nil || len(ifStmt.Body.List) != 1 {
		return nil, intCallGuardPattern{}, false
	}
	callFn, ok := s.parseCallStmtIntArg(ifStmt.Body.List[0], loopVar)
	if !ok {
		return nil, intCallGuardPattern{}, false
	}
	guard, ok := parseIntCallGuardCond(ifStmt.Cond, loopVar)
	if !ok {
		return nil, intCallGuardPattern{}, false
	}
	return callFn, guard, true
}

func parseIntCallGuardCond(cond ast.Expr, loopVar string) (intCallGuardPattern, bool) {
	binary, ok := unparenExpr(cond).(*ast.BinaryExpr)
	if !ok {
		return intCallGuardPattern{}, false
	}
	switch binary.Op {
	case token.LSS, token.LEQ, token.GTR, token.GEQ, token.EQL, token.NEQ:
	default:
		return intCallGuardPattern{}, false
	}
	if ident, ok := unparenExpr(binary.X).(*ast.Ident); ok && ident.Name == loopVar {
		c, ok := parseIntConst(binary.Y)
		if !ok {
			return intCallGuardPattern{}, false
		}
		return intCallGuardPattern{
			kind:   et.ForIntCallGuardKindCmp,
			op:     binary.Op,
			constV: c,
		}, true
	}
	modExpr, ok := unparenExpr(binary.X).(*ast.BinaryExpr)
	if !ok || modExpr.Op != token.REM {
		return intCallGuardPattern{}, false
	}
	modLoop, ok := unparenExpr(modExpr.X).(*ast.Ident)
	if !ok || modLoop.Name != loopVar {
		return intCallGuardPattern{}, false
	}
	modBase, ok := parseIntConst(modExpr.Y)
	if !ok || modBase == 0 {
		return intCallGuardPattern{}, false
	}
	c, ok := parseIntConst(binary.Y)
	if !ok {
		return intCallGuardPattern{}, false
	}
	return intCallGuardPattern{
		kind:    et.ForIntCallGuardKindModCmp,
		op:      binary.Op,
		constV:  c,
		modBase: modBase,
	}, true
}

func (s *Scope) parseCallStmtIntArg(stmt ast.Stmt, loopVar string) (func(int), bool) {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return nil, false
	}
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok || len(callExpr.Args) != 1 {
		return nil, false
	}
	arg, ok := callExpr.Args[0].(*ast.Ident)
	if !ok || arg.Name != loopVar {
		return nil, false
	}
	id := stringifyExpr(callExpr.Fun, 0)
	fn, found := s.lookupFunction(id)
	if !found {
		return nil, false
	}
	switch actual := fn.(type) {
	case func(int):
		return actual, true
	case func(int) int:
		return func(v int) { _ = actual(v) }, true
	}
	return nil, false
}
