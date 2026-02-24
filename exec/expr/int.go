package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"unsafe"
)

// Int represents int result expression
type Int struct {
	State    *exec.State
	NewState exec.New
	compute  internal.Compute
	ptr      unsafe.Pointer
}

// Compute computes int expr
func (e *Int) Compute() int {
	return *(*int)(e.compute(e.ptr))
}

// ComputeWithState computes bool expr
func (e *Int) ComputeWithState(state *exec.State) bool {
	return *(*bool)(e.compute(state.Pointer()))
}

// NewInt creates int expr
func NewInt(newState exec.New, compute internal.Compute) *Int {
	state := newState()
	return &Int{compute: compute, NewState: newState, State: state, ptr: state.Pointer()}
}
