package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"unsafe"
)

// String represents string result expression
type String struct {
	State    *exec.State
	NewState exec.New
	compute  internal.Compute
	ptr      unsafe.Pointer
}

// Compute computes string expr
func (e *String) Compute() string {
	return *(*string)(e.compute(e.ptr))
}

// ComputeWithState computes bool expr
func (e *String) ComputeWithState(state *exec.State) bool {
	return *(*bool)(e.compute(state.Pointer()))
}

// NewString crates string expression
func NewString(newState exec.New, compute internal.Compute) *String {
	state := newState()
	return &String{compute: compute, NewState: newState, State: state, ptr: state.Pointer()}
}
