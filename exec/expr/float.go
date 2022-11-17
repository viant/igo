package expr

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
)

//Float64 represents float64 result expression
type Float64 struct {
	State   *exec.State
	compute internal.Compute
}

//Compute computes float64 expression
func (e *Float64) Compute() float64 {
	return *(*float64)(e.compute(e.State.Pointer()))
}

//NewFloat64 creates a float64 expression
func NewFloat64(variables *exec.State, compute internal.Compute) *Float64 {
	return &Float64{compute: compute, State: variables}
}
