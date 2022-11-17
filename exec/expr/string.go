package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
)

//String represents string result expression
type String struct {
	State    *exec.State
	NewState exec.New
	compute  internal.Compute
}

//Compute computes string expr
func (e *String) Compute() string {
	return *(*string)(e.compute(e.State.Pointer()))
}

//ComputeWithState computes bool expr
func (e *String) ComputeWithState(state *exec.State) bool {
	return *(*bool)(e.compute(state.Pointer()))
}

//NewString crates string expression
func NewString(newState exec.New, compute internal.Compute) *String {
	return &String{compute: compute, NewState: newState, State: newState()}
}
