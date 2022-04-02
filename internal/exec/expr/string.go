package expr

import (
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/state"
)

//String represents string result expression
type String struct {
	Vars    *state.State
	compute exec.Compute
}

//Compute computes string expr
func (e *String) Compute() string {
	return *(*string)(e.compute(e.Vars.Pointer()))
}

//NewString crates string expression
func NewString(variables *state.State, compute exec.Compute) *String {
	return &String{compute: compute, Vars: variables}
}
