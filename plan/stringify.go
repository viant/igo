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
	stringify(expr, &builder, &curr, depth)
	return builder.String()
}

func isDepthReached(current *int, depth int) bool {
	if depth == 0 {
		return false
	}
	return *current >= depth
}

func stringify(expr ast.Expr, builder *strings.Builder, current *int, depth int) {
	if isDepthReached(current, depth) {
		return
	}
	switch actual := expr.(type) {
	case *ast.BasicLit:
		builder.WriteString(actual.Value)
		*current++
	case *ast.Ident:
		builder.WriteString(actual.Name)
		*current++
	case *ast.IndexExpr:
		stringify(actual.X, builder, current, depth)
		if isDepthReached(current, depth) {
			return
		}
		builder.WriteString("[")
		stringify(actual.Index, builder, current, depth)
		builder.WriteString("]")
	case *ast.SelectorExpr:
		stringify(actual.X, builder, current, depth)
		if isDepthReached(current, depth) {
			return
		}
		builder.WriteString(".")
		stringify(actual.Sel, builder, current, depth)
	case *ast.ParenExpr:
		builder.WriteString("(")
		stringify(actual.X, builder, current, depth)
		builder.WriteString(")")
	case *ast.CallExpr:
		stringify(actual.Fun, builder, current, depth)
		builder.WriteString("(")
		for i := 0; i < len(actual.Args); i++ {
			if i > 0 {
				builder.WriteString(",")
			}
			stringify(actual.Args[i], builder, current, depth)
		}
		builder.WriteString(")")
	case *ast.BinaryExpr:
		stringify(actual.X, builder, current, depth)
		builder.WriteString(actual.Op.String())
		stringify(actual.Y, builder, current, depth)
	case *ast.UnaryExpr:
		builder.WriteString(actual.Op.String())
		stringify(actual.X, builder, current, depth)
	case *ast.ArrayType:
		stringify(actual.Elt, builder, current, depth)
	case *ast.StarExpr:
		builder.WriteString("*")
		stringify(actual.X, builder, current, depth)

	default:
		panic(fmt.Sprintf("%T", actual))
	}
}
