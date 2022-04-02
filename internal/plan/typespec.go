package plan

import (
	"fmt"
	"go/ast"
	"reflect"
)

func (s *Scope) defineType(name string, expr ast.Expr) error {
	switch actual := expr.(type) {
	case *ast.StructType:
		var fields []reflect.StructField
		for _, field := range actual.Fields.List {
			fieldType, err := s.discoverType(field.Type)
			if err != nil {
				return err
			}
			for _, n := range field.Names {
				aField := reflect.StructField{
					Name: n.Name,
					Type: fieldType,
				}
				fields = append(fields, aField)
			}
		}
		s.RegisterNamedType(name, reflect.StructOf(fields))
		return nil

	}
	return fmt.Errorf("not yet sypported %T", expr)
}
