package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
)

//Int represents int result expression
type Int struct {
	State   *exec.State
	compute internal.Compute
}

//Compute computes int expr
func (e *Int) Compute() int {
	return *(*int)(e.compute(e.State.Pointer()))
}

//NewInt creates int expr
func NewInt(variables *exec.State, compute internal.Compute) *Int {
	return &Int{compute: compute, State: variables}
}
