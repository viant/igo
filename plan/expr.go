package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/exec/est"
	"github.com/viant/igo/exec/expr"
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
	variablesNew := exec.VariablesNew(s.mem.Type, *s.selectors)
	compute, err := exprNew(&est.Control{})
	if err != nil {
		return nil, err
	}
	variables := variablesNew()
	return expr.NewInt(variables, compute), nil
}

//BoolExpression returns  bool expression
func (s *Scope) BoolExpression(exprStmt string) (*expr.Bool, error) {
	exprNew, exprType, err := s.compileExprStmt(exprStmt)
	if err != nil {
		return nil, err
	}
	if exprType.Kind() != reflect.Bool {
		return nil, fmt.Errorf("invalid expression type: %v", exprType.String())
	}
	variablesNew := exec.VariablesNew(s.mem.Type, *s.selectors)
	compute, err := exprNew(&est.Control{})
	if err != nil {
		return nil, err
	}
	variables := variablesNew()
	return expr.NewBool(variables, compute), nil
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
	variablesNew := exec.VariablesNew(s.mem.Type, *s.selectors)
	compute, err := exprNew(&est.Control{})
	if err != nil {
		return nil, err
	}
	variables := variablesNew()
	return expr.NewFloat64(variables, compute), nil
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
	variablesNew := exec.VariablesNew(s.mem.Type, *s.selectors)
	compute, err := exprNew(&est.Control{})
	if err != nil {
		return nil, err
	}
	variables := variablesNew()
	return expr.NewString(variables, compute), nil
}

//compileExprStmt parses and compile expression
func (s *Scope) compileExprStmt(expr string) (est.New, reflect.Type, error) {
	fn, err := s.compileFunction(expr)
	if err != nil {
		return nil, nil, err
	}
	stmt, ok := fn.Body.List[0].(*ast.ExprStmt)
	if !ok {
		return nil, nil, fmt.Errorf("expected %T, but had %T", stmt, fn.Body.List[0])
	}
	return s.compileExpr(stmt.X)
}

func (s *Scope) compileExpr(stmt ast.Expr) (est.New, reflect.Type, error) {
	switch z := stmt.(type) {
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
	case *ast.CompositeLit:
		return s.compileCompositeLiteral(0, z)
	}
	return nil, nil, fmt.Errorf("unsupported %T", stmt)
}

func (s *Scope) compileBinaryExpr(binaryExpr *ast.BinaryExpr) (est.New, reflect.Type, error) {
	opX, err := s.assembleOperand(binaryExpr.X, false)
	if err != nil {
		return nil, nil, err
	}
	opY, err := s.assembleOperand(binaryExpr.Y, false)
	if err != nil {
		return nil, nil, err
	}
	newFn, destType := est.NewBinaryExpr(binaryExpr.Op, opX, opY)
	if destType == nil {
		return nil, nil, fmt.Errorf("dest type was nil")
	}
	return newFn, destType, nil
}

func (s *Scope) assembleOperand(expr ast.Expr, defined bool) (*est.Operand, error) {
	var op *est.Operand
	if isSelector(expr) {
		sel, err := s.selector(expr, defined)
		if err != nil {
			return nil, err
		}
		op = est.NewOperand(sel, nil, nil, nil)
	} else if isBasicLit(expr) {
		return literalOperand(expr.(*ast.BasicLit))
	} else {
		newFn, eType, err := s.compileExpr(expr)
		if err != nil {
			return nil, err
		}
		op = est.NewOperand(nil, eType, newFn, nil)
	}
	return op, nil
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
