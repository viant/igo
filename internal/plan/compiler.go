package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/parser"
	"strings"
	"sync"
)

//Function compile function
func (s *Scope) Function(expr string) (interface{}, error) {
	newFn, err := s.compile(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to compile: %s, %w", expr, err)
	}
	stateNew := exec.StateNew(s.mem.Type, *s.selectors, nil, nil)
	pool := &sync.Pool{New: func() interface{} {
		return stateNew()
	}}
	compute, err := newFn(&et.Control{Flow: s.Flow})
	if err != nil {
		return nil, err
	}
	execution := exec.NewExecution(compute, pool, *s.in, *s.out)
	return execution.Func()
}

//Compile parses and compile simple golang expression into execution tree
func (s *Scope) Compile(expr string) (*exec.Executor, error) {
	newFn, err := s.compile(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to compile: %s, %w", expr, err)
	}
	var tracker *exec.Tracker
	if s.trackType != nil {
		tracker = exec.NewTracker(s.trackLen)
	}
	pool := &sync.Pool{}
	stateNew := exec.StateNew(s.mem.Type, *s.selectors, tracker, pool)

	pool.New = func() interface{} {
		return stateNew()
	}
	compute, err := newFn(&et.Control{Flow: s.Flow})
	if err != nil {
		return nil, err
	}
	result := exec.NewExecution(compute, pool, *s.in, *s.out)
	return result, err
}

func (s *Scope) compile(expr string) (et.New, error) {
	s.Metric.Reset()
	fn, err := s.compileFunction(expr)
	if err != nil {
		return nil, err
	}
	if err = s.assignParams(s.in, fn.Type.Params); err != nil {
		return nil, err
	}
	if err = s.assignParams(s.out, fn.Type.Results); err != nil {
		return nil, err
	}
	return s.compileBlockStmt(fn.Body, false)
}

func (s *Scope) assignParams(dest *[]*exec.Selector, fieldList *ast.FieldList) error {
	if fieldList == nil || len(fieldList.List) == 0 {
		return nil
	}
	params, err := s.paramSelectors(fieldList.List)
	if err != nil {
		return err
	}
	*dest = params
	return nil
}

func (s *Scope) compileFunction(code string) (*ast.FuncLit, error) {
	codeExpr := code
	if !strings.HasPrefix(strings.TrimSpace(code), "func") {
		codeExpr = `func() {` + code + `}`
	}
	tree, err := parser.ParseExpr(codeExpr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %s, %w", code, err)
	}
	return tree.(*ast.FuncLit), nil
}
