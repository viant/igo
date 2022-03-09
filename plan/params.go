package plan

import (
	"github.com/viant/igo/exec"
	"go/ast"
	"strconv"
)

func (s *Scope) paramSelectors(params []*ast.Field) ([]*exec.Selector, error) {
	var parameters []*exec.Selector
	for i, field := range params {
		paramType, err := s.discoverType(field.Type)
		if err != nil {
			return nil, err
		}
		if len(field.Names) == 0 {
			field.Names = []*ast.Ident{{Name: "Result" + strconv.Itoa(i)}}
		}
		for _, name := range field.Names {
			sel, err := s.selector(name, true)
			if err != nil {
				return nil, err
			}
			if err = s.adjust(sel, paramType); err != nil {
				return nil, err
			}
			parameters = append(parameters, sel)
		}
	}
	return parameters, nil
}
