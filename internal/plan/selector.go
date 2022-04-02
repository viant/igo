package plan

import (
	"fmt"
	"github.com/viant/igo/state"
	"github.com/viant/xunsafe"
	"go/ast"
	"reflect"
	"strconv"
	"strings"
)

//Selector returns a selector
func (s *Scope) Selector(expr string) (*state.Selector, error) {
	selector := s.lookup(expr)
	if selector != nil {
		return selector, nil
	}
	panic(fmt.Sprintf("undefied variable: %v %v", expr, s.index))
	return nil, fmt.Errorf("undefied variable: %v %v", expr, s.index)
}

func (s *Scope) lookup(expr string) *state.Selector {
	if idx, ok := s.index[s.prefix+expr]; ok {
		return (*s.selectors)[idx]
	}
	//check selector in upstream scopes
	for i := len(s.upstream) - 1; i >= 0; i-- {
		if idx, ok := s.index[s.upstream[i]+expr]; ok {
			return (*s.selectors)[idx]
		}
	}
	//check global scope
	if idx, ok := s.index[expr]; ok {
		return (*s.selectors)[idx]
	}
	return nil
}

//DefineVariable defines variables
func (s *Scope) DefineVariable(name string, fType reflect.Type) (*state.Selector, error) {
	id := s.varID(name)
	if s.lookup(id) != nil {
		return nil, fmt.Errorf("variable %v already defined", fType.Name())
	}
	sel := &state.Selector{Field: &xunsafe.Field{Name: s.varName(name), Type: fType}, ID: id}
	if fType != nil {
		_ = s.RegisterType(fType)
	}
	sel.Ancestors = []*state.Selector{}
	var err error
	if fType != nil {
		err = s.appendSelector(sel)
	}
	return sel, err
}

func (s *Scope) newTransient() (*state.Selector, error) {
	id := "_transient" + strconv.Itoa(*s.transients)
	*s.transients++
	return s.DefineVariable(id, nil)
}


func (s *Scope) selector(expr ast.Expr, define bool) (*state.Selector, error) {
	id := stringifyExpr(expr, 0)
	if sel := s.lookup(id); sel != nil {
		return sel, nil
	}
	var parent *state.Selector
	i := 0
	if id == "nil" {
		return &state.Selector{ID: "nil", Field: &xunsafe.Field{Name: "nil"}}, nil
	}

	for ; ; i++ {
		prefix := stringifyExpr(expr, i+1)
		if sel := s.lookup(prefix); sel != nil {
			parent = sel
			continue
		}
		break
	}
	if parent == nil {
		if define {
			return s.DefineVariable(id, nil)
		}
		return s.Selector(id)
	}
	prefix := stringifyExpr(expr, i)
	index := strings.Index(id, prefix)
	if index == -1 {
		return s.Selector(id)
	}
	leaf := id[index+len(prefix):]
	switch leaf[0] {
	case '.':
		leaf = leaf[1:]
	case '[':
		leaf = "_" + leaf
	}
	leafExpr, err := expression(leaf)
	if err != nil {
		return nil, err
	}
	return s.addSelector(parent, leafExpr)
}

func (s *Scope) addSelector(parent *state.Selector, expr ast.Expr) (*state.Selector, error) {
	var err error
	sel := &state.Selector{Field: &xunsafe.Field{}}
	sel.Ancestors = make([]*state.Selector, len(parent.Ancestors)+1)
	if len(parent.Ancestors) > 0 {
		copy(sel.Ancestors, parent.Ancestors)
	}
	sel.Ancestors[len(sel.Ancestors)-1] = parent
	switch actual := expr.(type) {
	case *ast.Ident:
		sel.Field = xunsafe.FieldByName(parent.Type, actual.Name)
		if sel.Field == nil {
			return nil, fmt.Errorf("failed to lookup %v", actual.Name)
		}
		sel.ID = parent.ID + "." + sel.Name
		err = s.appendSelector(sel)
	case *ast.IndexExpr:
		return sel, s.addIndexSelector(parent, actual, sel)
	case *ast.SelectorExpr:
		return s.addSelectorNode(parent, actual)
	default:
		return nil, fmt.Errorf("unsupported selector %T", actual)
	}
	return sel, err
}

func (s *Scope) addSelectorNode(parent *state.Selector, actual *ast.SelectorExpr) (*state.Selector, error) {
	var err error
	parent, err = s.addSelector(parent, actual.X)
	if err != nil {
		return nil, err
	}
	return s.addSelector(parent, actual.Sel)
}

func (s *Scope) addIndexSelector(parent *state.Selector, actual *ast.IndexExpr, sel *state.Selector) error {
	sel.Type = parent.Type.Elem()
	sel.ID = parent.ID + "[" + stringifyExpr(actual.Index, 0) + "]"
	sel.Slice = xunsafe.NewSlice(parent.Type)
	operand, err := s.assembleOperand(actual.Index, false)
	if err != nil {
		return err
	}
	sel.Index, err = operand.NewOperand(nil)
	return err
}

func (s *Scope) varID(name string) string {
	return s.prefix + name
}

func (s *Scope) varName(name string) string {
	varName := name
	if !strings.HasPrefix(name, s.prefix) {
		varName = s.prefix + name
	}
	return varName
}

func (s *Scope) appendSelector(sel *state.Selector) error {
	index := len(*s.selectors)
	if sel.ID == "" {
		return fmt.Errorf("selector ID was empty")
	}
	if sel.Type == nil {
		return fmt.Errorf("selector %v type was empty", sel.Name)
	}
	if s.lookup(sel.ID) != nil {
		return fmt.Errorf("variable %v already defined", sel.Name)
	}
	s.index[sel.ID] = uint16(index)
	sel.Pos = uint16(index)
	if sel.Type.ConvertibleTo(errType) {
		sel.Type = errType
	}
	if len(sel.Ancestors) == 0 {
		sel.Field = s.mem.addField(sel.Name, sel.Type)
	}
	sel.Pathway = state.SelectorPathway(sel)
	if sel.Index != nil {
		if sel.Slice == nil {
			return fmt.Errorf("index slice was empty: %v", sel.Name)
		}
	}
	sel.IsErrorType = sel.Type.ConvertibleTo(errType)
	*s.selectors = append(*s.selectors, sel)
	return nil
}

//adjust selector type from inferred expression (t nsType)
func (s *Scope) adjust(selector *state.Selector, t reflect.Type) error {
	if selector.Type == nil {
		selector.Type = t
	}
	if s.lookup(selector.ID) == nil {
		if err := s.appendSelector(selector); err != nil {
			return err
		}
	}
	return nil
}
