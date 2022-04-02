package expr

import (
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/state"
)

//Bool represents bool result expression
type Bool struct {
	Vars    *state.State
	compute exec.Compute
}

//Compute computes bool expr
func (e *Bool) Compute() bool {
	return *(*bool)(e.compute(e.Vars.Pointer()))
}

//NewBool crates a bool expr
func NewBool(variables *state.State, compute exec.Compute) *Bool {
	return &Bool{compute: compute, Vars: variables}
}
