package plan

import (
	"fmt"
	"go/ast"
	"strings"
)

//stringifyExpr returns sting representation of expression
func stringifyExpr(expr ast.Expr, depth int) string {
	builder := strings.Builder{}
	curr := 0
	if err := stringify(expr, &builder, &curr, depth); err != nil {
		panic(err)
	}
	return builder.String()
}

//StringifyExpr returns sting representation of expression
func (s *Scope) StringifyExpr(expr ast.Expr) (string, error) {
	builder := strings.Builder{}
	curr := 0
	err := stringify(expr, &builder, &curr, 0)
	return builder.String(), err
}

func isDepthReached(current *int, depth int) bool {
	if depth == 0 {
		return false
	}
	return *current >= depth
}

func stringify(expr ast.Expr, builder *strings.Builder, current *int, depth int) error {
	if isDepthReached(current, depth) {
		return nil
	}
	switch actual := expr.(type) {
	case *ast.BasicLit:
		builder.WriteString(actual.Value)
		*current++
	case *ast.Ident:
		builder.WriteString(actual.Name)
		*current++
	case *ast.IndexExpr:
		if err := stringify(actual.X, builder, current, depth); err != nil {
			return err
		}
		if isDepthReached(current, depth) {
			return nil
		}
		builder.WriteString("[")
		if err := stringify(actual.Index, builder, current, depth); err != nil {
			return err
		}
		builder.WriteString("]")
	case *ast.SelectorExpr:
		if err := stringify(actual.X, builder, current, depth); err != nil {
			return err
		}
		if isDepthReached(current, depth) {
			return nil
		}
		builder.WriteString(".")
		return stringify(actual.Sel, builder, current, depth)
	case *ast.ParenExpr:
		builder.WriteString("(")
		if err := stringify(actual.X, builder, current, depth); err != nil {
			return err
		}
		builder.WriteString(")")
	case *ast.CallExpr:
		if err := stringify(actual.Fun, builder, current, depth); err != nil {
			return err
		}
		builder.WriteString("(")
		for i := 0; i < len(actual.Args); i++ {
			if i > 0 {
				builder.WriteString(",")
			}
			if err := stringify(actual.Args[i], builder, current, depth); err != nil {
				return err
			}
		}
		builder.WriteString(")")
	case *ast.BinaryExpr:
		if err := stringify(actual.X, builder, current, depth); err != nil {
			return err
		}
		builder.WriteString(actual.Op.String())
		return stringify(actual.Y, builder, current, depth)
	case *ast.UnaryExpr:
		builder.WriteString(actual.Op.String())
		return stringify(actual.X, builder, current, depth)
	case *ast.ArrayType:
		return stringify(actual.Elt, builder, current, depth)
	case *ast.StarExpr:
		builder.WriteString("*")
		return stringify(actual.X, builder, current, depth)
	default:
		return fmt.Errorf("unsupported node: %T", actual)
	}
	return nil
}
