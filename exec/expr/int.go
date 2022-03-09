package expr

import (
	"github.com/viant/igo/exec"
)

//Int represents int result expression
type Int struct {
	Vars    *exec.Variables
	compute exec.Compute
}

//Compute computes int expr
func (e *Int) Compute() int {
	return *(*int)(e.compute(e.Vars.Pointer()))
}

//NewInt creates int expr
func NewInt(variables *exec.Variables, compute exec.Compute) *Int {
	return &Int{compute: compute, Vars: variables}
}
