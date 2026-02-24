package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal/et"
	"github.com/viant/xunsafe"
	"go/ast"
	"reflect"
	"strings"
)

func (s *Scope) compileCallExpr(callExpr *ast.CallExpr) (et.New, []reflect.Type, error) {
	s.Metric.FlagCall()
	id := stringifyExpr(callExpr.Fun, 0)
	fnNew, fType, ok, err := s.compileBuildIn(id, callExpr)
	if ok {
		if fType == nil {
			return nil, nil, fmt.Errorf("builtin %s does not return a value", id)
		}

		return fnNew, []reflect.Type{fType}, err
	}
	aCaller, args, retTypes, err := s.compileCallExprSignature(callExpr)
	if err != nil {
		return nil, nil, err
	}
	newFn, err := et.NewCaller(aCaller, args)
	if err != nil {
		return nil, nil, err
	}
	return newFn, retTypes, nil
}

func (s *Scope) compileCallExprStmt(callExpr *ast.CallExpr) (et.New, error) {
	s.Metric.FlagCall()
	id := stringifyExpr(callExpr.Fun, 0)
	fnNew, _, ok, err := s.compileBuildIn(id, callExpr)
	if ok {
		return fnNew, err
	}
	caller, args, _, err := s.compileCallExprSignature(callExpr)
	if err != nil {
		return nil, err
	}
	return et.NewCallerStmt(caller, args)
}

func (s *Scope) compileCallExprAssign(callExpr *ast.CallExpr, dest []*et.Operand) (et.New, error) {
	s.Metric.FlagCall()
	id := stringifyExpr(callExpr.Fun, 0)
	fnNew, fType, ok, err := s.compileBuildIn(id, callExpr)
	if ok {
		if fType == nil {
			return nil, fmt.Errorf("builtin %s does not return a value", id)
		}
		_ = s.adjust(dest[0].Selector, fType)
		return et.NewAssignExpr(nil, dest[0], et.NewOperand(nil, fType, fnNew, nil)), nil
	}
	caller, args, ret, err := s.compileCallExprSignature(callExpr)
	if err != nil {
		return nil, err
	}
	if len(dest) == 1 && len(ret) == 1 && ret[0].Kind() == reflect.Interface &&
		dest[0] != nil && (dest[0].Type == nil || dest[0].Type.Kind() == reflect.Interface) {
		id := stringifyExpr(callExpr.Fun, 0)
		_, fn, err := s.discoverMethod(id)
		if err != nil {
			return nil, err
		}
		if fn == nil {
			var ok bool
			fn, ok = s.lookupFunction(id)
			if !ok {
				return nil, fmt.Errorf("failed to lookup function: %v", id)
			}
		}
		if variadic, ok := fn.(func(interface{}, ...interface{}) interface{}); ok {
			caller = newVariadicInterfaceInterfaceCaller(variadic)
		}
	}
	for i := range dest {
		if ret[i].Kind() == reflect.Interface {
			if dest[i].Type != nil {
				continue
			}
		}
		if err = s.adjust(dest[i].Selector, ret[i]); err != nil {
			return nil, err
		}
		dest[i].Type = xunsafe.NewType(ret[i])
	}
	return et.NewCallExprAssign(caller, args, dest)
}

var emptyInt16s = []uint16{}

func (s *Scope) updateTrackerLen(sel *exec.Selector) {
	if s.trackType == nil || (!strings.HasPrefix(sel.ID, s.trackRoot)) {
		return
	}
	positions := sel.XPos()[1:] //skip root level
	for _, pos := range positions {
		if int(pos) > s.trackLen {
			s.trackLen = int(pos)
		}
	}
}

func (s *Scope) trackerXPos(sel *exec.Selector) []uint16 {
	if s.trackType == nil || (!strings.HasPrefix(sel.ID, s.trackRoot)) {
		return emptyInt16s
	}
	return sel.XPos()[1:] //skip root level
}

func (s *Scope) compileCallExprSignature(callExpr *ast.CallExpr) (exec.Caller, []*et.Operand, []reflect.Type, error) {
	id := stringifyExpr(callExpr.Fun, 0)
	var ok bool
	var holder *exec.Selector
	fn, ok := s.lookupFunction(id)
	if !ok {
		var err error
		holder, fn, err = s.discoverMethod(id)
		if err != nil {
			return nil, nil, nil, err
		}
		if fn == nil {
			return nil, nil, nil, fmt.Errorf("failed to lookup function: %v", id)
		}
	}

	args, err := s.compileCallArgs(callExpr, holder)
	if err != nil {
		return nil, nil, nil, err
	}
	if compiled, ok := fn.(*compiledInlineFunc); ok {
		return compiled.caller, args, compiled.retTypes, nil
	}
	funType := reflect.TypeOf(fn)
	var retTypes = make([]reflect.Type, funType.NumOut())
	for i := range retTypes {
		retTypes[i] = funType.Out(i)
	}
	if id == "append" && len(args) == 2 && args[0] != nil && args[0].Type != nil && args[1] != nil && args[1].Type != nil {
		sliceType := args[0].Type.Type()
		if sliceType.Kind() == reflect.Slice {
			elemType := sliceType.Elem()
			if elemType == args[1].Type.Type() && elemType.Kind() != reflect.Interface {
				return newAppendOneCaller(sliceType), args, retTypes, nil
			}
		}
	}

	if variadic, ok := fn.(func(...interface{}) (int, error)); ok {
		return newVariadicInterfaceIntErrCaller(variadic), args, retTypes, nil
	}
	if variadic, ok := fn.(func(string, ...interface{})); ok {
		return newVariadicStringVoidCaller(variadic), args, retTypes, nil
	}
	if variadic, ok := fn.(func(string, ...interface{}) string); ok {
		return newVariadicStringStringCaller(variadic), args, retTypes, nil
	}
	if variadic, ok := fn.(func(string, ...interface{}) (int, error)); ok {
		return newVariadicStringIntErrCaller(variadic), args, retTypes, nil
	}
	if unaryString, ok := fn.(func(string) string); ok {
		return &fnSStringCaller{fn: unaryString}, args, retTypes, nil
	}
	if binaryStringBool, ok := fn.(func(string, string) bool); ok {
		return &fnSSBoolCaller{fn: binaryStringBool}, args, retTypes, nil
	}
	aCaller := exec.AsCaller(fn)
	if aCaller == nil {
		aCaller = newCaller(fn)
	}

	return aCaller, args, retTypes, nil
}

func (s *Scope) discoverMethod(id string) (*exec.Selector, interface{}, error) {
	holder, method, err := s.selectorFun(id)
	if err != nil || holder == nil {
		return nil, nil, err
	}
	fn, err := s.lookupMethod(holder, method)
	return holder, fn, err
}

func (s *Scope) selectorFun(id string) (*exec.Selector, string, error) {
	pos := strings.LastIndex(id, ".")
	if pos == -1 {
		return nil, id, nil
	}
	holderName := id[:pos]
	fName := id[pos+1:]
	sel, err := s.Selector(holderName)
	if err != nil {
		return nil, "", err
	}
	return sel, fName, nil
}

func (s *Scope) compileCallArgs(callExpr *ast.CallExpr, holder *exec.Selector) ([]*et.Operand, error) {
	argLength := len(callExpr.Args)
	if holder != nil {
		argLength++
	}
	var args = make([]*et.Operand, argLength)
	if len(args) == 0 {
		return args, nil
	}
	var err error
	i := 0
	if holder != nil {
		args[i] = et.NewOperand(holder, nil, nil, nil)
		i++
	}
	for _, expr := range callExpr.Args {
		if args[i], err = s.assembleOperand(expr, false); err != nil {
			return nil, err
		}
		i++
	}
	return args, nil
}

func (s *Scope) lookupMethod(sel *exec.Selector, name string) (interface{}, error) {
	method, ok := sel.Type.MethodByName(name)
	if !ok {
		return nil, fmt.Errorf("failed to locate %v.%v", sel.Type.String(), name)
	}
	return method.Func.Interface(), nil
}
