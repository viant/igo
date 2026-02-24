package plan

import (
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustParseForStmt(t *testing.T, src string) *ast.ForStmt {
	t.Helper()
	expr, err := parser.ParseExpr("func(){" + src + "}")
	if !assert.NoError(t, err) {
		return nil
	}
	fn, ok := expr.(*ast.FuncLit)
	if !assert.True(t, ok) || !assert.Len(t, fn.Body.List, 1) {
		return nil
	}
	stmt, ok := fn.Body.List[0].(*ast.ForStmt)
	if !assert.True(t, ok) {
		return nil
	}
	return stmt
}

func TestParseForIntCond_DataDriven(t *testing.T) {
	testCases := []struct {
		name    string
		cond    string
		varName string
		op      token.Token
		limit   int
		ok      bool
	}{
		{name: "direct-lss", cond: "i < 10", varName: "i", op: token.LSS, limit: 10, ok: true},
		{name: "direct-leq", cond: "i <= 10", varName: "i", op: token.LEQ, limit: 10, ok: true},
		{name: "direct-gtr", cond: "i > 10", varName: "i", op: token.GTR, limit: 10, ok: true},
		{name: "direct-geq-negative", cond: "i >= -3", varName: "i", op: token.GEQ, limit: -3, ok: true},
		{name: "direct-eql", cond: "i == 7", varName: "i", op: token.EQL, limit: 7, ok: true},
		{name: "direct-neq", cond: "i != 7", varName: "i", op: token.NEQ, limit: 7, ok: true},
		{name: "symmetric-gtr", cond: "10 > i", varName: "i", op: token.LSS, limit: 10, ok: true},
		{name: "symmetric-leq", cond: "10 <= i", varName: "i", op: token.GEQ, limit: 10, ok: true},
		{name: "symmetric-lss", cond: "10 < i", varName: "i", op: token.GTR, limit: 10, ok: true},
		{name: "symmetric-geq", cond: "10 >= i", varName: "i", op: token.LEQ, limit: 10, ok: true},
		{name: "symmetric-eql", cond: "10 == i", varName: "i", op: token.EQL, limit: 10, ok: true},
		{name: "symmetric-neq", cond: "10 != i", varName: "i", op: token.NEQ, limit: 10, ok: true},
		{name: "unsupported-op", cond: "i + 1", ok: false},
		{name: "non-ident-var", cond: "(i+1) < 10", ok: false},
		{name: "non-int-limit", cond: "i < j", ok: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cond, err := parser.ParseExpr(tc.cond)
			if !assert.NoError(t, err) {
				return
			}
			varName, op, limit, ok := parseForIntCond(cond)
			assert.Equal(t, tc.ok, ok)
			if !tc.ok {
				return
			}
			assert.Equal(t, tc.varName, varName)
			assert.Equal(t, tc.op, op)
			assert.Equal(t, tc.limit, limit)
		})
	}
}

func TestParseForIntPost_DataDriven(t *testing.T) {
	testCases := []struct {
		name  string
		src   string
		delta int
		ok    bool
	}{
		{name: "inc", src: "for i:=0; i<10; i++ {}", delta: 1, ok: true},
		{name: "dec", src: "for i:=0; i<10; i-- {}", delta: -1, ok: true},
		{name: "add-assign", src: "for i:=0; i<10; i += 2 {}", delta: 2, ok: true},
		{name: "sub-assign", src: "for i:=0; i<10; i -= 3 {}", delta: -3, ok: true},
		{name: "assign-add", src: "for i:=0; i<10; i = i + 4 {}", delta: 4, ok: true},
		{name: "assign-sub", src: "for i:=0; i<10; i = i - 5 {}", delta: -5, ok: true},
		{name: "assign-other-var", src: "for i:=0; i<10; i = j + 1 {}", ok: false},
		{name: "assign-unsupported-op", src: "for i:=0; i<10; i = i * 2 {}", ok: false},
		{name: "zero-delta-add", src: "for i:=0; i<10; i += 0 {}", ok: false},
		{name: "zero-delta-sub", src: "for i:=0; i<10; i -= 0 {}", ok: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stmt := mustParseForStmt(t, tc.src)
			if stmt == nil {
				return
			}
			delta, ok := parseForIntPost(stmt.Post, "i")
			assert.Equal(t, tc.ok, ok)
			if tc.ok {
				assert.Equal(t, tc.delta, delta)
			}
		})
	}
}

func TestParseForBodyIntAccAssign_DataDriven(t *testing.T) {
	testCases := []struct {
		name       string
		src        string
		ok         bool
		dest       string
		assignOp   token.Token
		useLoopVar bool
		rhsConst   int
		guard      bool
		guardOp    token.Token
		guardConst int
	}{
		{name: "sum-add-loopvar", src: "for i:=0;i<10;i++ { sum += i }", ok: true, dest: "sum", assignOp: token.ADD_ASSIGN, useLoopVar: true},
		{name: "sum-sub-loopvar", src: "for i:=0;i<10;i++ { sum -= i }", ok: true, dest: "sum", assignOp: token.SUB_ASSIGN, useLoopVar: true},
		{name: "sum-mul-loopvar", src: "for i:=0;i<10;i++ { sum *= i }", ok: true, dest: "sum", assignOp: token.MUL_ASSIGN, useLoopVar: true},
		{name: "sum-div-const", src: "for i:=0;i<10;i++ { sum /= 2 }", ok: true, dest: "sum", assignOp: token.QUO_ASSIGN, useLoopVar: false, rhsConst: 2},
		{name: "sum-add-const", src: "for i:=0;i<10;i++ { sum += 3 }", ok: true, dest: "sum", assignOp: token.ADD_ASSIGN, useLoopVar: false, rhsConst: 3},
		{name: "guarded", src: "for i:=0;i<10;i++ { if i < 5 { sum += i } }", ok: true, dest: "sum", assignOp: token.ADD_ASSIGN, useLoopVar: true, guard: true, guardOp: token.LSS, guardConst: 5},
		{name: "guarded-leq", src: "for i:=0;i<10;i++ { if i <= 5 { sum += i } }", ok: true, dest: "sum", assignOp: token.ADD_ASSIGN, useLoopVar: true, guard: true, guardOp: token.LEQ, guardConst: 5},
		{name: "guarded-gtr", src: "for i:=0;i<10;i++ { if i > 5 { sum += i } }", ok: true, dest: "sum", assignOp: token.ADD_ASSIGN, useLoopVar: true, guard: true, guardOp: token.GTR, guardConst: 5},
		{name: "guarded-geq", src: "for i:=0;i<10;i++ { if i >= 5 { sum += i } }", ok: true, dest: "sum", assignOp: token.ADD_ASSIGN, useLoopVar: true, guard: true, guardOp: token.GEQ, guardConst: 5},
		{name: "guarded-eql", src: "for i:=0;i<10;i++ { if i == 5 { sum += i } }", ok: true, dest: "sum", assignOp: token.ADD_ASSIGN, useLoopVar: true, guard: true, guardOp: token.EQL, guardConst: 5},
		{name: "guarded-neq", src: "for i:=0;i<10;i++ { if i != 5 { sum += i } }", ok: true, dest: "sum", assignOp: token.ADD_ASSIGN, useLoopVar: true, guard: true, guardOp: token.NEQ, guardConst: 5},
		{name: "invalid-else", src: "for i:=0;i<10;i++ { if i < 5 { sum += i } else { sum += 1 } }", ok: false},
		{name: "invalid-guard-non-loop-var", src: "for i:=0;i<10;i++ { if j < 5 { sum += i } }", ok: false},
		{name: "invalid-nonloopvar", src: "for i:=0;i<10;i++ { sum += j }", ok: false},
		{name: "invalid-two-statements", src: "for i:=0;i<10;i++ { sum += i; sum += 1 }", ok: false},
		{name: "invalid-unsupported-assign-token", src: "for i:=0;i<10;i++ { sum = i }", ok: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stmt := mustParseForStmt(t, tc.src)
			if stmt == nil {
				return
			}
			pattern, ok := parseForBodyIntAccAssign(stmt.Body, "i")
			assert.Equal(t, tc.ok, ok)
			if !tc.ok {
				return
			}
			assert.Equal(t, tc.dest, pattern.dest)
			assert.Equal(t, tc.assignOp, pattern.assignOp)
			assert.Equal(t, tc.useLoopVar, pattern.useLoopVar)
			assert.Equal(t, tc.rhsConst, pattern.rhsConst)
			assert.Equal(t, tc.guard, pattern.guardEnable)
			if tc.guard {
				assert.Equal(t, tc.guardOp, pattern.guardOp)
				assert.Equal(t, tc.guardConst, pattern.guardConst)
			}
		})
	}
}

func TestParseIntConst_DataDriven(t *testing.T) {
	testCases := []struct {
		name  string
		expr  string
		value int
		ok    bool
	}{
		{name: "int", expr: "12", value: 12, ok: true},
		{name: "negative", expr: "-12", value: -12, ok: true},
		{name: "plus-unary-unsupported", expr: "+12", ok: false},
		{name: "identifier", expr: "i", ok: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := parser.ParseExpr(tc.expr)
			if !assert.NoError(t, err) {
				return
			}
			got, ok := parseIntConst(e)
			assert.Equal(t, tc.ok, ok)
			if tc.ok {
				assert.Equal(t, tc.value, got)
			}
		})
	}
}

func TestParseForBodyIntCall_DataDriven(t *testing.T) {
	scope := NewScope()
	scope.RegisterFunc("sinkInt", func(int) {})
	scope.RegisterFunc("mapInt", func(v int) int { return v + 1 })
	scope.RegisterFunc("sinkString", func(string) {})
	testCases := []struct {
		name string
		src  string
		ok   bool
	}{
		{name: "func-int", src: "for i:=0;i<10;i++ { sinkInt(i) }", ok: true},
		{name: "func-int-int", src: "for i:=0;i<10;i++ { mapInt(i) }", ok: true},
		{name: "wrong-arg-count", src: "for i:=0;i<10;i++ { sinkInt(i, i) }", ok: false},
		{name: "wrong-arg-var", src: "for i:=0;i<10;i++ { sinkInt(j) }", ok: false},
		{name: "unknown-func", src: "for i:=0;i<10;i++ { missing(i) }", ok: false},
		{name: "wrong-signature", src: "for i:=0;i<10;i++ { sinkString(i) }", ok: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stmt := mustParseForStmt(t, tc.src)
			if stmt == nil {
				return
			}
			fn, ok := scope.parseForBodyIntCall(stmt.Body, "i")
			assert.Equal(t, tc.ok, ok)
			if ok {
				assert.NotNil(t, fn)
				fn(1)
			}
		})
	}
}

func TestParseForBodyIntCallGuard_DataDriven(t *testing.T) {
	scope := NewScope()
	scope.RegisterFunc("sinkInt", func(int) {})
	scope.RegisterFunc("mapInt", func(v int) int { return v + 1 })
	testCases := []struct {
		name      string
		src       string
		ok        bool
		guardKind int
		guardOp   token.Token
		constV    int
		modBase   int
	}{
		{
			name:      "cmp-guard",
			src:       "for i:=0;i<10;i++ { if i < 5 { sinkInt(i) } }",
			ok:        true,
			guardKind: et.ForIntCallGuardKindCmp,
			guardOp:   token.LSS,
			constV:    5,
		},
		{
			name:      "mod-guard",
			src:       "for i:=0;i<10;i++ { if i % 3 != 0 { mapInt(i) } }",
			ok:        true,
			guardKind: et.ForIntCallGuardKindModCmp,
			guardOp:   token.NEQ,
			constV:    0,
			modBase:   3,
		},
		{
			name: "invalid-body-not-if",
			src:  "for i:=0;i<10;i++ { sinkInt(i) }",
			ok:   false,
		},
		{
			name: "invalid-else",
			src:  "for i:=0;i<10;i++ { if i < 5 { sinkInt(i) } else { sinkInt(i) } }",
			ok:   false,
		},
		{
			name: "invalid-guard-var",
			src:  "for i:=0;i<10;i++ { if j < 5 { sinkInt(i) } }",
			ok:   false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stmt := mustParseForStmt(t, tc.src)
			if stmt == nil {
				return
			}
			fn, guard, ok := scope.parseForBodyIntCallGuard(stmt.Body, "i")
			assert.Equal(t, tc.ok, ok)
			if !ok {
				return
			}
			assert.NotNil(t, fn)
			assert.Equal(t, tc.guardKind, guard.kind)
			assert.Equal(t, tc.guardOp, guard.op)
			assert.Equal(t, tc.constV, guard.constV)
			assert.Equal(t, tc.modBase, guard.modBase)
			fn(1)
		})
	}
}
