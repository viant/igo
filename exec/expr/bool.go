package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"unsafe"
)

// Bool represents bool result expression
type Bool struct {
	State    *exec.State
	NewState exec.New
	compute  internal.Compute
	ptr      unsafe.Pointer
}

// Compute computes bool expr
func (e *Bool) Compute() bool {
	return *(*bool)(e.compute(e.ptr))
}

// ComputeWithState computes bool expr
func (e *Bool) ComputeWithState(state *exec.State) bool {
	return *(*bool)(e.compute(state.Pointer()))
}

// NewBool crates a bool expr
func NewBool(newState exec.New, compute internal.Compute) *Bool {
	state := newState()
	return &Bool{compute: compute, NewState: newState, State: state, ptr: state.Pointer()}
}
