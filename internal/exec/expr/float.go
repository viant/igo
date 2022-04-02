package expr

import (
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/state"
)

//Float64 represents float64 result expression
type Float64 struct {
	Vars    *state.State
	compute exec.Compute
}

//Compute computes float64 expression
func (e *Float64) Compute() float64 {
	return *(*float64)(e.compute(e.Vars.Pointer()))
}

//NewFloat64 creates a float64 expression
func NewFloat64(variables *state.State, compute exec.Compute) *Float64 {
	return &Float64{compute: compute, Vars: variables}
}
