package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	expr "github.com/viant/igo/exec/expr"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"reflect"
)

//IntExpression returns an int expression
func (s *Scope) IntExpression(exprStmt string) (*expr.Int, error) {
	exprNew, exprType, err := s.compileExprStmt(exprStmt)
	if err != nil {
		return nil, err
	}
	if exprType.Kind() != reflect.Int {
		return nil, fmt.Errorf("invalid expression type: %v", exprType.String())
	}
	variablesNew := exec.StateNew(s.mem.Type, *s.selectors, nil, nil)
	compute, err := exprNew(&et.Control{})
	if err != nil {
		return nil, err
	}
	return expr.NewInt(variablesNew, compute), nil
}

//BoolExpression returns  bool expression
func (s *Scope) BoolExpression(exprStmt string) (*expr.Bool, error) {
	sel, _ := s.DefineVariable("_boolExpr", reflect.TypeOf(true))
	output := []*exec.Selector{sel}
	s.out = &output

	exprNew, exprType, err := s.compileExprStmt(exprStmt)
	if err != nil {
		return nil, err
	}
	if exprType.Kind() != reflect.Bool {
		return nil, fmt.Errorf("invalid expression type: %v", exprType.String())
	}
	variablesNew := exec.StateNew(s.mem.Type, *s.selectors, nil, nil)
	control := *s.Control
	compute, err := exprNew(&control)
	if err != nil {
		return nil, err
	}
	return expr.NewBool(variablesNew, compute), nil
}

//Float64Expression returns  float64 expression
func (s *Scope) Float64Expression(exprStmt string) (*expr.Float64, error) {
	exprNew, exprType, err := s.compileExprStmt(exprStmt)
	if err != nil {
		return nil, err
	}
	if exprType.Kind() != reflect.Float64 {
		return nil, fmt.Errorf("invalid expression type: %v", exprType.String())
	}
	variablesNew := exec.StateNew(s.mem.Type, *s.selectors, nil, nil)
	control := *s.Control
	compute, err := exprNew(&control)
	if err != nil {
		return nil, err
	}
	return expr.NewFloat64(variablesNew, compute), nil
}

//StringExpression returns  string expression
func (s *Scope) StringExpression(exprStmt string) (*expr.String, error) {
	exprNew, exprType, err := s.compileExprStmt(exprStmt)
	if err != nil {
		return nil, err
	}
	if exprType.Kind() != reflect.String {
		return nil, fmt.Errorf("invalid expression type: %v", exprType.String())
	}
	variablesNew := exec.StateNew(s.mem.Type, *s.selectors, nil, nil)
	control := *s.Control
	compute, err := exprNew(&control)
	if err != nil {
		return nil, err
	}
	return expr.NewString(variablesNew, compute), nil
}

//compileExprStmt parses and compile expression
func (s *Scope) compileExprStmt(expr string) (et.New, reflect.Type, error) {
	fn, err := s.compileFunction(expr)
	if err != nil {
		return nil, nil, err
	}
	stmt, ok := fn.Body.List[0].(*ast.ExprStmt)
	if ok {
		return s.compileExpr(stmt.X)
	}
	newFn, err := s.compileBlockStmt(fn.Body, true)
	return newFn, reflect.TypeOf(true), err
}

func (s *Scope) compileExpr(expr ast.Expr) (et.New, reflect.Type, error) {
	if s.exprListener != nil {
		newExpr, err := s.exprListener(expr)
		if err != nil {
			return nil, nil, err
		}
		if newExpr != nil {
			expr = newExpr
		}
	}

	switch z := expr.(type) {
	case *ast.BinaryExpr:
		return s.compileBinaryExpr(z)
	case *ast.ParenExpr:
		return s.compileExpr(z.X)
	case *ast.CallExpr:
		newFn, types, err := s.compileCallExpr(z)
		if err != nil {
			return nil, nil, err
		}
		var rType reflect.Type
		if len(types) > 0 {
			rType = types[0]
		}
		return newFn, rType, nil
	case *ast.UnaryExpr:
		return s.compileUnaryExpr(z)
	case *ast.StarExpr:
		return s.compileUnaryStarExpr(z)

	case *ast.CompositeLit:
		return s.compileCompositeLiteral(0, z)
	}
	return nil, nil, fmt.Errorf("unsupported %T", expr)
}

func (s *Scope) compileBinaryExpr(binaryExpr *ast.BinaryExpr) (et.New, reflect.Type, error) {
	opX, err := s.assembleOperand(binaryExpr.X, false)
	if err != nil {
		return nil, nil, err
	}
	opY, err := s.assembleOperand(binaryExpr.Y, false)
	if err != nil {
		return nil, nil, err
	}
	z, _ := s.newTransient()
	opZ := &et.Operand{Selector: z}
	newFn, destType := et.NewBinaryExpr(binaryExpr.Op, opX, opY, opZ)
	if destType == nil {
		return nil, nil, fmt.Errorf("dest type was nil")
	}
	_ = s.adjust(z, destType)
	return newFn, destType, nil
}

func (s *Scope) assembleOperand(expr ast.Expr, defined bool) (*et.Operand, error) {
	var op *et.Operand
	if isSelector(expr) {
		sel, err := s.selector(expr, defined)
		if err != nil {
			return nil, err
		}
		op = et.NewOperand(sel, nil, nil, nil)
	} else if isBasicLit(expr) {
		return literalOperand(expr.(*ast.BasicLit))
	} else {
		newFn, eType, err := s.compileExpr(expr)
		if err != nil {
			return nil, err
		}
		op = et.NewOperand(nil, eType, newFn, nil)
	}
	return op, nil
}

//Parse returns an ast for the supploed expression
func (s *Scope) Parse(expr string) (ast.Stmt, error) {
	fn, err := s.compileFunction(expr)
	if err != nil {
		return nil, err
	}
	stmt := fn.Body.List[0]
	return stmt, nil
}

//expression returns ast for supplied expr or error
func expression(expr string) (ast.Expr, error) {
	scope := &Scope{}
	fn, err := scope.compileFunction(expr)
	if err != nil {
		return nil, err
	}
	stmt := fn.Body.List[0].(*ast.ExprStmt)
	return stmt.X, nil
}
