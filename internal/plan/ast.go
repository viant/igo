package plan

import "go/ast"

func isSelector(node ast.Expr) bool {
	switch node.(type) {
	case *ast.SelectorExpr, *ast.Ident, *ast.IndexExpr:
		return true
	}
	return false
}

func isCallExpr(node ast.Expr) bool {
	switch node.(type) {
	case *ast.CallExpr:
		return true
	}
	return false
}

func isBasicLit(node ast.Expr) bool {
	switch node.(type) {
	case *ast.BasicLit:
		return true
	}
	return false
}

func isArrayIdentifier(node ast.Node) bool {
	switch node.(type) {
	case *ast.ArrayType:
		return true
	}
	return false
}
