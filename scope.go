package igo

import (
	"github.com/viant/igo/internal/plan"
	"github.com/viant/igo/option"
)

//Scope represents a scope
type Scope struct {
	plan.Scope
}

//NewScope creates a scope
func NewScope(option ...option.Option) *Scope {
	return &Scope{*plan.NewScope(option...)}
}
