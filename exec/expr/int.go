package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
)

//Int represents int result expression
type Int struct {
	State    *exec.State
	NewState exec.New
	compute  internal.Compute
}

//Compute computes int expr
func (e *Int) Compute() int {
	return *(*int)(e.compute(e.State.Pointer()))
}

//ComputeWithState computes bool expr
func (e *Int) ComputeWithState(state *exec.State) bool {
	return *(*bool)(e.compute(state.Pointer()))
}

//NewInt creates int expr
func NewInt(newState exec.New, compute internal.Compute) *Int {
	return &Int{compute: compute, NewState: newState, State: newState()}
}
