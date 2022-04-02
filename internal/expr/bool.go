package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
)

//Bool represents bool result expression
type Bool struct {
	State   *exec.State
	compute internal.Compute
}

//Compute computes bool expr
func (e *Bool) Compute() bool {
	return *(*bool)(e.compute(e.State.Pointer()))
}

//NewBool crates a bool expr
func NewBool(variables *exec.State, compute internal.Compute) *Bool {
	return &Bool{compute: compute, State: variables}
}
