package igo

import (
	"github.com/viant/igo/internal/plan"
	"github.com/viant/igo/option"
)

//Scope represents a scope
type Scope struct {
	*plan.Scope
}

func (s *Scope) SubScope(options ...option.Option) *Scope {
	return &Scope{Scope: s.Scope.SubScope(options...)}
}

//NewScope creates a scope
func NewScope(option ...option.Option) *Scope {
	return &Scope{Scope: plan.NewScope(option...)}
}
