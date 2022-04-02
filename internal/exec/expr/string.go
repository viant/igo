package expr

import (
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/exec"
)

//String represents string result expression
type String struct {
	Vars    *exec.State
	compute exec.Compute
}

//Compute computes string expr
func (e *String) Compute() string {
	return *(*string)(e.compute(e.Vars.Pointer()))
}

//NewString crates string expression
func NewString(variables *exec.State, compute exec.Compute) *String {
	return &String{compute: compute, Vars: variables}
}
