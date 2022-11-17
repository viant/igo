package option

import "go/ast"

//StmtListener represents stmt listener
type StmtListener func(stmt ast.Stmt) (ast.Stmt, error)

//ExprListener represent expr listener
type ExprListener func(stmt ast.Expr) (ast.Expr, error)
