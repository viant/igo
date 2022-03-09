package expr

import (
	"github.com/viant/igo/exec"
)

//Bool represents bool result expression
type Bool struct {
	Vars    *exec.Variables
	compute exec.Compute
}

//Compute computes bool expr
func (e *Bool) Compute() bool {
	return *(*bool)(e.compute(e.Vars.Pointer()))
}

//NewBool crates a bool expr
func NewBool(variables *exec.Variables, compute exec.Compute) *Bool {
	return &Bool{compute: compute, Vars: variables}
}
