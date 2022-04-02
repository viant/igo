package internal

import "github.com/viant/igo/exec"

//Executor abstraction holding execution plan (execution syntax tree) to execute it
type Executor struct {
	compute Compute
	In      []string //params variable identity
	Out     []string //result variable identity
}

//Exec executes execution plan
func (e *Executor) Exec(state *exec.State) {
	ptr := state.Pointer()
	AsFlow(ptr).Reset()
	e.compute(state.Pointer())
}

//NewExecution creates a new execution
func NewExecution(compute Compute) *Executor {
	return &Executor{compute: compute}
}
