package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/exec/est"
	"go/ast"
	"reflect"
)

func (s *Scope) compileMapper(holder *exec.Selector, expr *ast.CallExpr) (est.New, reflect.Type, error) {
	args := expr.Args
	if len(args) == 0 {
		return nil, nil, fmt.Errorf("%v insuficient arguments", stringifyExpr(expr.Fun, 0))
	}
	funcLit, ok := args[0].(*ast.FuncLit)
	if !ok {
		return nil, nil, fmt.Errorf("%v invalid args", stringifyExpr(expr.Fun, 0))
	}
	scope := s.subScope()
	parameters, err := scope.paramSelectors(funcLit.Type.Params.List)
	if err != nil {
		return nil, nil, err
	}
	results, err := scope.paramSelectors(funcLit.Type.Results.List)
	if err != nil {
		return nil, nil, err
	}
	scope.out = &results
	body, err := scope.compileBlockStmt(funcLit.Body, false)
	return est.NewMapper(holder, parameters, results, body)
}
