package plan

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func (s *Scope) compileForStmt(forStmt *ast.ForStmt, labelID uint32) (et.New, error) {
	var err error
	scope := s.SubScope()
	var cond *et.Operand
	var init, post, body et.New
	if forStmt.Init != nil {
		if init, err = s.compileStmt(forStmt.Init); err != nil {
			return nil, err
		}
	}
	if newFor, ok, err := s.tryCompileMergeSortLikeForStmt(scope, forStmt, init, labelID); ok || err != nil {
		return newFor, err
	}
	if newFor, ok, err := s.tryCompileSpecializedIntForStmt(scope, forStmt, init, labelID); ok || err != nil {
		return newFor, err
	}
	if forStmt.Post != nil {
		if post, err = scope.compileStmt(forStmt.Post); err != nil {
			return nil, err
		}
	}
	if forStmt.Cond != nil {
		if cond, err = s.assembleOperand(forStmt.Cond, false); err != nil {
			return nil, err
		}
	}
	if body, err = scope.compileBlockStmt(forStmt.Body, true); err != nil {
		return nil, err
	}
	return et.NewForStmt(cond, init, post, body, labelID)
}

func (s *Scope) tryCompileSpecializedIntForStmt(scope *Scope, forStmt *ast.ForStmt, init et.New, labelID uint32) (et.New, bool, error) {
	name, cmpOp, limit, ok := parseForIntCond(forStmt.Cond)
	if !ok {
		return nil, false, nil
	}
	delta, ok := parseForIntPost(forStmt.Post, name)
	if !ok {
		return nil, false, nil
	}
	sel, err := scope.selector(&ast.Ident{Name: name}, false)
	if err != nil {
		return nil, false, nil
	}
	if sel.Pathway != exec.PathwayDirect || sel.Type.Kind() != reflect.Int {
		return nil, false, nil
	}
	if labelID == 0 {
		if complexPattern, ok := parseComplexStateBodyPattern(forStmt.Body, name); ok {
			sumSel, err := scope.selector(complexPattern.sumExpr, false)
			if err == nil && sumSel.Pathway == exec.PathwayDirect && sumSel.Type.Kind() == reflect.Int {
				countSel, err := scope.selector(complexPattern.countExpr, false)
				if err == nil && countSel.Pathway == exec.PathwayDirect && countSel.Type.Kind() == reflect.Int {
					modeSel, err := scope.selector(complexPattern.modeExpr, false)
					if err == nil && modeSel.Pathway == exec.PathwayDirect && modeSel.Type.Kind() == reflect.Int {
						outputSel, err := scope.selector(complexPattern.outputExpr, false)
						if err == nil && outputSel.Pathway == exec.PathwayDirect && outputSel.Type.Kind() == reflect.Int {
							flagsSel, err := scope.selector(complexPattern.flagsExpr, false)
							if err == nil && flagsSel.Pathway == exec.PathwayDirect && flagsSel.Type.Kind() == reflect.Map {
								flagsKey := flagsSel.Type.Key()
								flagsElem := flagsSel.Type.Elem()
								if flagsKey.Kind() == reflect.String && flagsElem.Kind() == reflect.Bool {
									scoresSel, err := scope.selector(complexPattern.scoresExpr, false)
									if err == nil && scoresSel.Pathway == exec.PathwayDirect && scoresSel.Type.Kind() == reflect.Map {
										scoresKey := scoresSel.Type.Key()
										scoresElem := scoresSel.Type.Elem()
										if scoresKey.Kind() == reflect.Int && scoresElem.Kind() == reflect.Int {
											newFor, err := et.NewForComplexStateStmt(init, sel.Offset(), limit, sumSel.Offset(), countSel.Offset(), modeSel.Offset(), outputSel.Offset(), flagsSel.Offset(), scoresSel.Offset(), complexPattern.breakLimit, complexPattern.breakOnGE)
											return newFor, true, err
										}
									}
								}
							}
						}
					}
				}
			}
		}
		if callFn, ok := s.parseForBodyIntCall(forStmt.Body, name); ok {
			newFor, err := et.NewForIntCallConstStmt(init, sel.Offset(), cmpOp, limit, delta, callFn)
			return newFor, true, err
		}
		if callFn, guard, ok := s.parseForBodyIntCallGuard(forStmt.Body, name); ok {
			newFor, err := et.NewForIntCallGuardConstStmt(init, sel.Offset(), cmpOp, limit, delta, callFn, guard.kind, guard.op, guard.constV, guard.modBase)
			return newFor, true, err
		}
		if acc, ok := parseForBodyIntAccAssign(forStmt.Body, name); ok {
			destName := acc.dest
			destSel, err := scope.selector(&ast.Ident{Name: destName}, false)
			if err == nil && destSel.Pathway == exec.PathwayDirect && destSel.Type.Kind() == reflect.Int {
				if !acc.guardEnable && acc.assignOp == token.ADD_ASSIGN && acc.useLoopVar {
					newFor, err := et.NewForIntAddConstStmt(init, destSel.Offset(), sel.Offset(), cmpOp, limit, delta)
					return newFor, true, err
				}
				newFor, err := et.NewForIntAccConstStmt(init, destSel.Offset(), sel.Offset(), cmpOp, limit, delta, acc.assignOp, acc.useLoopVar, acc.rhsConst, acc.guardEnable, acc.guardOp, acc.guardConst)
				return newFor, true, err
			}
		}
	}
	body, err := scope.compileBlockStmt(forStmt.Body, true)
	if err != nil {
		return nil, true, err
	}
	newFor, err := et.NewForIntConstStmt(init, body, sel.Offset(), cmpOp, limit, delta, labelID)
	return newFor, true, err
}

func parseForIntCond(cond ast.Expr) (string, token.Token, int, bool) {
	binary, ok := cond.(*ast.BinaryExpr)
	if !ok {
		return "", token.ILLEGAL, 0, false
	}
	if left, ok := binary.X.(*ast.Ident); ok {
		limit, ok := parseIntConst(binary.Y)
		if !ok {
			return "", token.ILLEGAL, 0, false
		}
		switch binary.Op {
		case token.LSS, token.LEQ, token.GTR, token.GEQ, token.EQL, token.NEQ:
			return left.Name, binary.Op, limit, true
		}
		return "", token.ILLEGAL, 0, false
	}
	if right, ok := binary.Y.(*ast.Ident); ok {
		limit, ok := parseIntConst(binary.X)
		if !ok {
			return "", token.ILLEGAL, 0, false
		}
		switch binary.Op {
		case token.LSS:
			return right.Name, token.GTR, limit, true
		case token.LEQ:
			return right.Name, token.GEQ, limit, true
		case token.GTR:
			return right.Name, token.LSS, limit, true
		case token.GEQ:
			return right.Name, token.LEQ, limit, true
		case token.EQL, token.NEQ:
			return right.Name, binary.Op, limit, true
		}
	}
	return "", token.ILLEGAL, 0, false
}

func parseForIntPost(post ast.Stmt, name string) (int, bool) {
	if assign, ok := post.(*ast.AssignStmt); ok && len(assign.Lhs) == 1 && len(assign.Rhs) == 1 {
		ident, ok := assign.Lhs[0].(*ast.Ident)
		if !ok || ident.Name != name {
			return 0, false
		}
		switch assign.Tok {
		case token.ADD_ASSIGN:
			step, ok := parseIntConst(assign.Rhs[0])
			if !ok {
				return 0, false
			}
			return step, step != 0
		case token.SUB_ASSIGN:
			step, ok := parseIntConst(assign.Rhs[0])
			if !ok {
				return 0, false
			}
			return -step, step != 0
		case token.ASSIGN:
			rhs, ok := assign.Rhs[0].(*ast.BinaryExpr)
			if !ok {
				return 0, false
			}
			lhsIdent, ok := rhs.X.(*ast.Ident)
			if !ok || lhsIdent.Name != name {
				return 0, false
			}
			step, ok := parseIntConst(rhs.Y)
			if !ok || step == 0 {
				return 0, false
			}
			switch rhs.Op {
			case token.ADD:
				return step, true
			case token.SUB:
				return -step, true
			}
		}
	}
	incDec, ok := post.(*ast.IncDecStmt)
	if !ok {
		return 0, false
	}
	ident, ok := incDec.X.(*ast.Ident)
	if !ok || ident.Name != name {
		return 0, false
	}
	switch incDec.Tok {
	case token.INC:
		return 1, true
	case token.DEC:
		return -1, true
	}
	return 0, false
}

type intAccPattern struct {
	dest        string
	assignOp    token.Token
	useLoopVar  bool
	rhsConst    int
	guardEnable bool
	guardOp     token.Token
	guardConst  int
}

func parseForBodyIntAccAssign(body *ast.BlockStmt, loopVar string) (intAccPattern, bool) {
	if body == nil || len(body.List) != 1 {
		return intAccPattern{}, false
	}
	switch stmt := body.List[0].(type) {
	case *ast.AssignStmt:
		p, ok := parseIntAccAssign(stmt, loopVar)
		return p, ok
	case *ast.IfStmt:
		if stmt.Else != nil || len(stmt.Body.List) != 1 {
			return intAccPattern{}, false
		}
		cond, ok := stmt.Cond.(*ast.BinaryExpr)
		if !ok {
			return intAccPattern{}, false
		}
		ident, ok := cond.X.(*ast.Ident)
		if !ok || ident.Name != loopVar {
			return intAccPattern{}, false
		}
		guardConst, ok := parseIntConst(cond.Y)
		if !ok {
			return intAccPattern{}, false
		}
		switch cond.Op {
		case token.LSS, token.LEQ, token.GTR, token.GEQ, token.EQL, token.NEQ:
		default:
			return intAccPattern{}, false
		}
		assign, ok := stmt.Body.List[0].(*ast.AssignStmt)
		if !ok {
			return intAccPattern{}, false
		}
		p, ok := parseIntAccAssign(assign, loopVar)
		if !ok {
			return intAccPattern{}, false
		}
		p.guardEnable = true
		p.guardOp = cond.Op
		p.guardConst = guardConst
		return p, true
	}
	return intAccPattern{}, false
}

func parseIntAccAssign(assign *ast.AssignStmt, loopVar string) (intAccPattern, bool) {
	if len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return intAccPattern{}, false
	}
	switch assign.Tok {
	case token.ADD_ASSIGN, token.SUB_ASSIGN, token.MUL_ASSIGN, token.QUO_ASSIGN:
	default:
		return intAccPattern{}, false
	}
	dest, ok := assign.Lhs[0].(*ast.Ident)
	if !ok {
		return intAccPattern{}, false
	}
	p := intAccPattern{dest: dest.Name, assignOp: assign.Tok}
	if src, ok := assign.Rhs[0].(*ast.Ident); ok {
		if src.Name != loopVar {
			return intAccPattern{}, false
		}
		p.useLoopVar = true
		return p, true
	}
	c, ok := parseIntConst(assign.Rhs[0])
	if !ok {
		return intAccPattern{}, false
	}
	p.useLoopVar = false
	p.rhsConst = c
	return p, true
}

type complexStateBodyPattern struct {
	sumExpr    ast.Expr
	countExpr  ast.Expr
	modeExpr   ast.Expr
	outputExpr ast.Expr
	flagsExpr  ast.Expr
	scoresExpr ast.Expr
	breakLimit int
	breakOnGE  bool
}

func parseComplexStateBodyPattern(body *ast.BlockStmt, loopVar string) (complexStateBodyPattern, bool) {
	if body == nil || (len(body.List) != 6 && len(body.List) != 5) {
		return complexStateBodyPattern{}, false
	}
	pattern, ok := parseComplexStateIfElse(body.List[0], loopVar)
	if !ok {
		return complexStateBodyPattern{}, false
	}
	modeExpr, ok := parseComplexStateSwitch(body.List[1], loopVar)
	if !ok {
		return complexStateBodyPattern{}, false
	}
	outputExpr, flagsExpr, scoresExpr, ok := parseComplexStateHotMap(body.List[2], loopVar)
	if !ok {
		return complexStateBodyPattern{}, false
	}
	breakLimit, breakOnGE, ok := parseComplexStateBreak(body.List[3], pattern.sumExpr)
	if !ok {
		return complexStateBodyPattern{}, false
	}
	var countExpr ast.Expr
	if len(body.List) == 6 {
		if !parseComplexStateContinue(body.List[4], loopVar) {
			return complexStateBodyPattern{}, false
		}
		countExpr, ok = parseComplexStateCount(body.List[5])
		if !ok {
			return complexStateBodyPattern{}, false
		}
	} else {
		countExpr, ok = parseComplexStateCountGuard(body.List[4], loopVar)
		if !ok {
			return complexStateBodyPattern{}, false
		}
	}
	pattern.countExpr = countExpr
	pattern.modeExpr = modeExpr
	pattern.outputExpr = outputExpr
	pattern.flagsExpr = flagsExpr
	pattern.scoresExpr = scoresExpr
	pattern.breakLimit = breakLimit
	pattern.breakOnGE = breakOnGE
	return pattern, true
}

func parseComplexStateIfElse(stmt ast.Stmt, loopVar string) (complexStateBodyPattern, bool) {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok || ifStmt.Init != nil {
		return complexStateBodyPattern{}, false
	}
	cond, ok := unparenExpr(ifStmt.Cond).(*ast.BinaryExpr)
	if !ok || cond.Op != token.EQL {
		return complexStateBodyPattern{}, false
	}
	condRhs, ok := parseIntConst(cond.Y)
	if !ok || condRhs != 0 {
		return complexStateBodyPattern{}, false
	}
	andExpr, ok := unparenExpr(cond.X).(*ast.BinaryExpr)
	if !ok || andExpr.Op != token.AND {
		return complexStateBodyPattern{}, false
	}
	andLoop, ok := unparenExpr(andExpr.X).(*ast.Ident)
	if !ok || andLoop.Name != loopVar {
		return complexStateBodyPattern{}, false
	}
	andConst, ok := parseIntConst(unparenExpr(andExpr.Y))
	if !ok || andConst != 1 {
		return complexStateBodyPattern{}, false
	}
	if len(ifStmt.Body.List) != 1 {
		return complexStateBodyPattern{}, false
	}
	thenAssign, ok := ifStmt.Body.List[0].(*ast.AssignStmt)
	if !ok || thenAssign.Tok != token.ADD_ASSIGN || len(thenAssign.Lhs) != 1 || len(thenAssign.Rhs) != 1 {
		return complexStateBodyPattern{}, false
	}
	rhsLoop, ok := thenAssign.Rhs[0].(*ast.Ident)
	if !ok || rhsLoop.Name != loopVar {
		return complexStateBodyPattern{}, false
	}
	elseBlock, ok := ifStmt.Else.(*ast.BlockStmt)
	if !ok || len(elseBlock.List) != 1 {
		return complexStateBodyPattern{}, false
	}
	elseAssign, ok := elseBlock.List[0].(*ast.AssignStmt)
	if !ok || elseAssign.Tok != token.SUB_ASSIGN || len(elseAssign.Lhs) != 1 || len(elseAssign.Rhs) != 1 {
		return complexStateBodyPattern{}, false
	}
	if stringifyExpr(thenAssign.Lhs[0], 0) != stringifyExpr(elseAssign.Lhs[0], 0) {
		return complexStateBodyPattern{}, false
	}
	divExpr, ok := elseAssign.Rhs[0].(*ast.BinaryExpr)
	if !ok || divExpr.Op != token.QUO {
		return complexStateBodyPattern{}, false
	}
	divX, ok := divExpr.X.(*ast.Ident)
	if !ok || divX.Name != loopVar {
		return complexStateBodyPattern{}, false
	}
	divBy, ok := parseIntConst(divExpr.Y)
	if !ok || divBy != 2 {
		return complexStateBodyPattern{}, false
	}
	sumExpr := thenAssign.Lhs[0]
	return complexStateBodyPattern{
		sumExpr: sumExpr,
	}, true
}

func parseComplexStateSwitch(stmt ast.Stmt, loopVar string) (ast.Expr, bool) {
	sw, ok := stmt.(*ast.SwitchStmt)
	if !ok || sw.Init != nil || sw.Tag == nil || len(sw.Body.List) != 5 {
		return nil, false
	}
	tagExpr, ok := sw.Tag.(*ast.BinaryExpr)
	if !ok || tagExpr.Op != token.REM {
		return nil, false
	}
	tagLoop, ok := tagExpr.X.(*ast.Ident)
	if !ok || tagLoop.Name != loopVar {
		return nil, false
	}
	tagConst, ok := parseIntConst(tagExpr.Y)
	if !ok || tagConst != 5 {
		return nil, false
	}
	var modeExpr ast.Expr
	seen := map[int]bool{}
	defaultSeen := false
	for _, item := range sw.Body.List {
		clause, ok := item.(*ast.CaseClause)
		if !ok || len(clause.Body) != 1 {
			return nil, false
		}
		assign, ok := clause.Body[0].(*ast.AssignStmt)
		if !ok || assign.Tok != token.ASSIGN || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			return nil, false
		}
		if modeExpr == nil {
			modeExpr = assign.Lhs[0]
		} else if stringifyExpr(modeExpr, 0) != stringifyExpr(assign.Lhs[0], 0) {
			return nil, false
		}
		assignVal, ok := parseIntConst(assign.Rhs[0])
		if !ok {
			return nil, false
		}
		if len(clause.List) == 0 {
			if assignVal != 5 || defaultSeen {
				return nil, false
			}
			defaultSeen = true
			continue
		}
		if len(clause.List) != 1 {
			return nil, false
		}
		caseVal, ok := parseIntConst(clause.List[0])
		if !ok || caseVal < 0 || caseVal > 3 || seen[caseVal] {
			return nil, false
		}
		if assignVal != caseVal+1 {
			return nil, false
		}
		seen[caseVal] = true
	}
	if !defaultSeen || len(seen) != 4 || modeExpr == nil {
		return nil, false
	}
	return modeExpr, true
}

func parseComplexStateHotMap(stmt ast.Stmt, loopVar string) (ast.Expr, ast.Expr, ast.Expr, bool) {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok || ifStmt.Init != nil || ifStmt.Else != nil || len(ifStmt.Body.List) != 1 {
		return nil, nil, nil, false
	}
	flagIndex, ok := ifStmt.Cond.(*ast.IndexExpr)
	if !ok {
		return nil, nil, nil, false
	}
	flagKey, ok := flagIndex.Index.(*ast.BasicLit)
	if !ok || flagKey.Kind != token.STRING || flagKey.Value != "\"hot\"" {
		return nil, nil, nil, false
	}
	assign, ok := ifStmt.Body.List[0].(*ast.AssignStmt)
	if !ok || assign.Tok != token.ADD_ASSIGN || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return nil, nil, nil, false
	}
	scoreIndex, ok := assign.Rhs[0].(*ast.IndexExpr)
	if !ok {
		return nil, nil, nil, false
	}
	scoreMod, ok := scoreIndex.Index.(*ast.BinaryExpr)
	if !ok || scoreMod.Op != token.REM {
		return nil, nil, nil, false
	}
	scoreLoop, ok := scoreMod.X.(*ast.Ident)
	if !ok || scoreLoop.Name != loopVar {
		return nil, nil, nil, false
	}
	scoreConst, ok := parseIntConst(scoreMod.Y)
	if !ok || scoreConst != 4 {
		return nil, nil, nil, false
	}
	return assign.Lhs[0], flagIndex.X, scoreIndex.X, true
}

func parseComplexStateBreak(stmt ast.Stmt, sumExpr ast.Expr) (int, bool, bool) {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok || ifStmt.Init != nil || ifStmt.Else != nil || len(ifStmt.Body.List) != 1 {
		return 0, false, false
	}
	cond, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok || (cond.Op != token.GTR && cond.Op != token.GEQ) {
		return 0, false, false
	}
	if stringifyExpr(cond.X, 0) != stringifyExpr(sumExpr, 0) {
		return 0, false, false
	}
	breakLimit, ok := parseIntConst(cond.Y)
	if !ok {
		return 0, false, false
	}
	branch, ok := ifStmt.Body.List[0].(*ast.BranchStmt)
	if !ok || branch.Tok != token.BREAK || branch.Label != nil {
		return 0, false, false
	}
	return breakLimit, cond.Op == token.GEQ, true
}

func parseComplexStateContinue(stmt ast.Stmt, loopVar string) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok || ifStmt.Init != nil || ifStmt.Else != nil || len(ifStmt.Body.List) != 1 {
		return false
	}
	cond, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok || cond.Op != token.EQL {
		return false
	}
	modExpr, ok := cond.X.(*ast.BinaryExpr)
	if !ok || modExpr.Op != token.REM {
		return false
	}
	modLoop, ok := modExpr.X.(*ast.Ident)
	if !ok || modLoop.Name != loopVar {
		return false
	}
	modConst, ok := parseIntConst(modExpr.Y)
	if !ok || modConst != 7 {
		return false
	}
	condConst, ok := parseIntConst(cond.Y)
	if !ok || condConst != 0 {
		return false
	}
	branch, ok := ifStmt.Body.List[0].(*ast.BranchStmt)
	return ok && branch.Tok == token.CONTINUE && branch.Label == nil
}

func parseComplexStateCount(stmt ast.Stmt) (ast.Expr, bool) {
	inc, ok := stmt.(*ast.IncDecStmt)
	if !ok || inc.Tok != token.INC {
		return nil, false
	}
	return inc.X, true
}

func parseComplexStateCountGuard(stmt ast.Stmt, loopVar string) (ast.Expr, bool) {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok || ifStmt.Init != nil || ifStmt.Else != nil || len(ifStmt.Body.List) != 1 {
		return nil, false
	}
	cond, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok || cond.Op != token.NEQ {
		return nil, false
	}
	modExpr, ok := cond.X.(*ast.BinaryExpr)
	if !ok || modExpr.Op != token.REM {
		return nil, false
	}
	modLoop, ok := modExpr.X.(*ast.Ident)
	if !ok || modLoop.Name != loopVar {
		return nil, false
	}
	modConst, ok := parseIntConst(modExpr.Y)
	if !ok || modConst != 7 {
		return nil, false
	}
	condConst, ok := parseIntConst(cond.Y)
	if !ok || condConst != 0 {
		return nil, false
	}
	inc, ok := ifStmt.Body.List[0].(*ast.IncDecStmt)
	if !ok || inc.Tok != token.INC {
		return nil, false
	}
	return inc.X, true
}

func unparenExpr(expr ast.Expr) ast.Expr {
	for {
		if paren, ok := expr.(*ast.ParenExpr); ok {
			expr = paren.X
			continue
		}
		return expr
	}
}

func (s *Scope) parseForBodyIntCall(body *ast.BlockStmt, loopVar string) (func(int), bool) {
	if body == nil || len(body.List) != 1 {
		return nil, false
	}
	return s.parseCallStmtIntArg(body.List[0], loopVar)
}

func parseIntConst(expr ast.Expr) (int, bool) {
	switch actual := expr.(type) {
	case *ast.BasicLit:
		if actual.Kind != token.INT {
			return 0, false
		}
		v, err := strconv.Atoi(actual.Value)
		return v, err == nil
	case *ast.UnaryExpr:
		if actual.Op != token.SUB {
			return 0, false
		}
		v, ok := parseIntConst(actual.X)
		if !ok {
			return 0, false
		}
		return -v, true
	default:
		return 0, false
	}
}
