package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/parser"
	"reflect"
	"strings"
	"sync"
)

// Function compile function
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

// Compile parses and compile simple golang expression into execution tree
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
	*s.defers = (*s.defers)[:0]
	fn, err := s.compileFunction(expr)
	if err != nil {
		return nil, err
	}
	return s.compileFromFuncLit(fn)
}

func (s *Scope) compileFromFuncLit(fn *ast.FuncLit) (et.New, error) {
	var err error
	if err = s.assignParams(s.in, fn.Type.Params); err != nil {
		return nil, err
	}
	if err = s.assignParams(s.out, fn.Type.Results); err != nil {
		return nil, err
	}
	body, err := s.compileBlockStmt(fn.Body, false)
	if err != nil {
		return nil, err
	}
	if len(*s.defers) == 0 {
		return body, nil
	}
	deferStmts := make([]et.New, 0, len(*s.defers))
	for i := len(*s.defers) - 1; i >= 0; i-- {
		deferStmts = append(deferStmts, (*s.defers)[i])
	}
	deferBlock := et.NewGroupStmt(deferStmts, false)
	group := []et.New{body, deferBlock}
	return et.NewGroupStmt(group, false), nil
}

func (s *Scope) functionFromFuncLit(fn *ast.FuncLit) (interface{}, error) {
	inner := NewScope()
	for name, t := range s.types {
		inner.RegisterNamedType(name, t)
	}
	for name, fnRef := range s.funcs {
		inner.RegisterFunc(name, fnRef)
	}

	newFn, err := inner.compileFromFuncLit(fn)
	if err != nil {
		return nil, err
	}
	stateNew := exec.StateNew(inner.mem.Type, *inner.selectors, nil, nil)
	pool := &sync.Pool{New: func() interface{} {
		return stateNew()
	}}
	compute, err := newFn(&et.Control{Flow: inner.Flow})
	if err != nil {
		return nil, err
	}
	execution := exec.NewExecution(compute, pool, *inner.in, *inner.out)
	retTypes := make([]reflect.Type, len(*inner.out))
	for i, out := range *inner.out {
		retTypes[i] = out.Type
	}
	return &compiledInlineFunc{
		caller:   newInlineExecutorCaller(execution),
		retTypes: retTypes,
		fnType:   execution.Signature(),
	}, nil
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
