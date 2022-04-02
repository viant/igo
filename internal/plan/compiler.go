package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/parser"
	"reflect"
	"strings"
	"sync"
)

//Function compile function
func (s *Scope) Function(expr string) (interface{}, error) {
	newFn, err := s.compile(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to compile: %s, %w", expr, err)
	}
	stateNew := exec.StateNew(s.mem.Type, *s.selectors)
	compute, err := newFn(&et.Control{Flow: s.Flow})
	if err != nil {
		return nil, err
	}
	execution := internal.NewExecution(compute)
	execution.In = exec.Selectors(*s.in).IDs()
	execution.Out = exec.Selectors(*s.out).IDs()
	var in, out []reflect.Type
	for i := range execution.In {
		in = append(in, (*s.in)[i].Type)
	}
	for i := range execution.Out {
		out = append(out, (*s.out)[i].Type)
	}
	fnType := reflect.FuncOf(in, out, false)
	pool := sync.Pool{New: func() interface{} {
		return stateNew()
	}}
	return reflect.MakeFunc(fnType, func(args []reflect.Value) (results []reflect.Value) {
		state := pool.Get().(*exec.State)
		for i, in := range execution.In {
			if err := state.SetValue(in, args[i].Interface()); err != nil {
				panic(err)
			}
		}
		execution.Exec(state)
		results = make([]reflect.Value, len(out))
		for i, out := range execution.Out {
			value, err := state.Value(out)
			if err != nil {
				panic(err)
			}
			results[i] = reflect.ValueOf(value)
		}
		pool.Put(state)
		return results
	}).Interface(), nil
}

//Compile parses and compile simple golang expression into execution tree
func (s *Scope) Compile(expr string) (*internal.Executor, exec.New, error) {
	newFn, err := s.compile(expr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compile: %s, %w", expr, err)
	}
	variablesNew := exec.StateNew(s.mem.Type, *s.selectors)
	compute, err := newFn(&et.Control{Flow: s.Flow})
	if err != nil {
		return nil, nil, err
	}
	result := internal.NewExecution(compute)
	result.In = exec.Selectors(*s.in).IDs()
	result.Out = exec.Selectors(*s.out).IDs()
	return result, variablesNew, err
}

func (s *Scope) compile(expr string) (et.New, error) {
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
	if !strings.HasPrefix(code, "func") {
		codeExpr = `func() {` + code + `}`
	}
	tree, err := parser.ParseExpr(codeExpr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %s, %w", code, err)
	}
	return tree.(*ast.FuncLit), nil
}
