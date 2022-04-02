package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
)

//String represents string result expression
type String struct {
	State   *exec.State
	compute internal.Compute
}

//Compute computes string expr
func (e *String) Compute() string {
	return *(*string)(e.compute(e.State.Pointer()))
}

//NewString crates string expression
func NewString(variables *exec.State, compute internal.Compute) *String {
	return &String{compute: compute, State: variables}
}
