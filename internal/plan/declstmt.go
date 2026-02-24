package plan

import (
	"fmt"
	"github.com/viant/igo/internal"
	"github.com/viant/igo/internal/et"
	"github.com/viant/xunsafe"
	"go/ast"
	"reflect"
	"unsafe"
)

func (s *Scope) compileDeclStmt(decl ast.Decl) (et.New, error) {
	switch actual := decl.(type) {
	case *ast.GenDecl:
		for _, specItem := range actual.Specs {
			switch spec := specItem.(type) {
			case *ast.ValueSpec:
				if spec.Type != nil {
					return s.declareTypedVariables(spec)
				}
				return s.defineVariables(spec)
			case *ast.TypeSpec:
				name := stringifyExpr(spec.Name, 0)
				if err := s.defineType(name, spec.Type); err != nil {
					return nil, err
				}
				return func(control *et.Control) (internal.Compute, error) {
					return func(ptr unsafe.Pointer) unsafe.Pointer {
						return ptr
					}, nil
				}, nil
			default:
				return nil, fmt.Errorf("unsupported declaration spec %T", spec)
			}
		}
	}
	return nil, fmt.Errorf("unsupported declaration %T", decl)
}

func (s *Scope) declareTypedVariables(spec *ast.ValueSpec) (et.New, error) {
	sType, err := s.discoverType(spec.Type)
	if err != nil {
		return nil, err
	}
	var state []et.New
	for _, ident := range spec.Names {
		sel, err := s.selector(ident, true)
		if err != nil {
			return nil, err
		}
		s.adjust(sel, sType)
		left := et.NewOperand(sel, sType, nil, nil)
		zero := reflect.New(sType)
		right := et.NewOperand(nil, sType, nil, nil)
		right.ValuePtr = xunsafe.ValuePointer(&zero)
		state = append(state, et.NewAssignExpr(nil, left, right))
	}
	return et.NewGroupStmt(state, false), nil

}

func (s *Scope) defineVariables(spec *ast.ValueSpec) (et.New, error) {
	var declrs []et.New
	if len(spec.Names) == len(spec.Values) {
		for i, n := range spec.Names {
			if fnLit, ok := spec.Values[i].(*ast.FuncLit); ok {
				fn, err := s.functionFromFuncLit(fnLit)
				if err != nil {
					return nil, err
				}
				s.RegisterFunc(n.Name, fn)
				continue
			}
			sel, err := s.selector(n, true)
			if err != nil {
				return nil, err
			}
			declr, destType, err := s.compileExpr(spec.Values[i])
			if err != nil {
				return nil, err
			}
			_ = s.adjust(sel, destType)
			left := et.NewOperand(sel, destType, nil, nil)
			right := et.NewOperand(nil, destType, declr, nil)
			declrs = append(declrs, et.NewAssignExpr(nil, left, right))
		}
	}
	return et.NewGroupStmt(declrs, false), nil
}
