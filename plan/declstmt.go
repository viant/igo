package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/exec/est"
	"go/ast"
	"reflect"
	"unsafe"
)

func (s *Scope) compileDeclStmt(decl ast.Decl) (est.New, error) {
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
				return func(control *est.Control) (exec.Compute, error) {
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

func (s *Scope) declareTypedVariables(spec *ast.ValueSpec) (est.New, error) {
	sType, err := s.discoverType(spec.Type)
	if err != nil {
		return nil, err
	}
	var vars []est.New
	for _, ident := range spec.Names {
		sel, err := s.selector(ident, true)
		if err != nil {
			return nil, err
		}
		s.adjust(sel, sType)
		left := est.NewOperand(sel, sType, nil, nil)
		right := est.NewOperand(sel, sType, nil, reflect.New(sType).Elem().Interface())
		vars = append(vars, est.NewAssignExpr(left, right))
	}
	return est.NewGroupStmt(vars, false), nil

}

func (s *Scope) defineVariables(spec *ast.ValueSpec) (est.New, error) {
	var declrs []est.New
	if len(spec.Names) == len(spec.Values) {
		for i, n := range spec.Names {
			sel, err := s.selector(n, true)
			if err != nil {
				return nil, err
			}
			declr, destType, err := s.compileExpr(spec.Values[i])
			if err != nil {
				return nil, err
			}
			s.adjust(sel, destType)
			left := est.NewOperand(sel, destType, nil, nil)
			right := est.NewOperand(nil, destType, declr, nil)
			declrs = append(declrs, est.NewAssignExpr(left, right))
		}
	}
	return est.NewGroupStmt(declrs, false), nil
}
