package plan

import (
	"fmt"
	"go/ast"
)

func (s *Scope) defineType(name string, expr ast.Expr) error {
	t, err := s.discoverType(expr)
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("not yet supported %T", expr)
	}
	s.RegisterNamedType(name, t)
	return nil
}
