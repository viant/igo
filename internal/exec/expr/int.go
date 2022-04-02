package expr

import (
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/state"
)

//Int represents int result expression
type Int struct {
	Vars    *state.State
	compute exec.Compute
}

//Compute computes int expr
func (e *Int) Compute() int {
	return *(*int)(e.compute(e.Vars.Pointer()))
}

//NewInt creates int expr
func NewInt(variables *state.State, compute exec.Compute) *Int {
	return &Int{compute: compute, Vars: variables}
}
