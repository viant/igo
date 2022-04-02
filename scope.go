package igo

import "github.com/viant/igo/internal/plan"

//Scope represents a scope
type Scope struct {
	plan.Scope
}

//NewScope creates a scope
func NewScope() *Scope {
	return &Scope{*plan.NewScope()}
}