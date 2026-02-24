package et

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"reflect"
	"unsafe"
)

// TypeSwitchCase represents a compiled type switch case.
type TypeSwitchCase struct {
	Types []reflect.Type
	Body  New
	// Assign, if present, receives asserted case value.
	Assign *exec.Selector
}

type typeSwitchCase struct {
	types  []reflect.Type
	assign *exec.Selector
	body   internal.Compute
}

type typeSwitchStmt struct {
	x       *exec.Operand
	cases   []typeSwitchCase
	defBody internal.Compute
	labelID uint32
}

func matchesType(actual, candidate reflect.Type) bool {
	if candidate == nil {
		return false
	}
	if candidate.Kind() == reflect.Interface {
		return actual.Implements(candidate) || actual.AssignableTo(candidate)
	}
	return actual.AssignableTo(candidate)
}

func (s *typeSwitchStmt) assignValue(dest *exec.Selector, ptr unsafe.Pointer, rv reflect.Value) {
	target := dest.Type
	if target.Kind() == reflect.Interface {
		dest.SetValue(dest.Upstream(ptr), rv.Interface())
		return
	}
	dst := reflect.NewAt(target, dest.Addr(ptr)).Elem()
	if rv.Type().AssignableTo(target) {
		dst.Set(rv)
		return
	}
	if rv.Type().ConvertibleTo(target) {
		dst.Set(rv.Convert(target))
	}
}

func (s *typeSwitchStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	flow := internal.AsFlow(ptr)
	xv := reflect.NewAt(s.x.Type.Type(), s.x.Compute(ptr)).Elem().Interface()
	if xv == nil {
		for i := range s.cases {
			aCase := &s.cases[i]
			matched := false
			for _, candidate := range aCase.types {
				if candidate == nil {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
			result := aCase.body(ptr)
			if flow.HasBreak() {
				if label := flow.Label(); label == 0 || label == s.labelID {
					flow.ClearBreak()
					return nil
				}
			}
			return result
		}
		if s.defBody != nil {
			result := s.defBody(ptr)
			if flow.HasBreak() {
				if label := flow.Label(); label == 0 || label == s.labelID {
					flow.ClearBreak()
					return nil
				}
			}
			return result
		}
		return nil
	}
	rv := reflect.ValueOf(xv)
	rt := rv.Type()
	for i := range s.cases {
		aCase := &s.cases[i]
		matched := false
		for _, candidate := range aCase.types {
			if matchesType(rt, candidate) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}
		if aCase.assign != nil {
			s.assignValue(aCase.assign, ptr, rv)
		}
		result := aCase.body(ptr)
		if flow.HasBreak() {
			if label := flow.Label(); label == 0 || label == s.labelID {
				flow.ClearBreak()
				return nil
			}
		}
		return result
	}
	if s.defBody != nil {
		result := s.defBody(ptr)
		if flow.HasBreak() {
			if label := flow.Label(); label == 0 || label == s.labelID {
				flow.ClearBreak()
				return nil
			}
		}
		return result
	}
	return nil
}

// NewTypeSwitchStmt creates a type switch statement.
func NewTypeSwitchStmt(x *Operand, cases []TypeSwitchCase, defaultBody New, labelID uint32) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		xOp, err := x.NewOperand(control)
		if err != nil {
			return nil, err
		}
		stmt := &typeSwitchStmt{x: xOp, cases: make([]typeSwitchCase, len(cases)), labelID: labelID}
		for i := range cases {
			body, err := cases[i].Body(control)
			if err != nil {
				return nil, err
			}
			stmt.cases[i] = typeSwitchCase{
				types:  cases[i].Types,
				assign: cases[i].Assign,
				body:   body,
			}
		}
		if defaultBody != nil {
			stmt.defBody, err = defaultBody(control)
			if err != nil {
				return nil, err
			}
		}
		return stmt.compute, nil
	}, nil
}
