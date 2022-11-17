package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
)

//Bool represents bool result expression
type Bool struct {
	State    *exec.State
	NewState exec.New
	compute  internal.Compute
}

//Compute computes bool expr
func (e *Bool) Compute() bool {
	return *(*bool)(e.compute(e.State.Pointer()))
}

//ComputeWithState computes bool expr
func (e *Bool) ComputeWithState(state *exec.State) bool {
	return *(*bool)(e.compute(state.Pointer()))
}

//NewBool crates a bool expr
func NewBool(newState exec.New, compute internal.Compute) *Bool {
	return &Bool{compute: compute, NewState: newState, State: newState()}
}
