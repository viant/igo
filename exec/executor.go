package exec

import "github.com/viant/igo/internal"

//Executor abstraction holding execution plan (execution syntax tree) to execute it
type Executor struct {
	compute internal.Compute
	In      []string //params variable identity
	Out     []string //result variable identity
}

//Exec executes execution plan
func (e *Executor) Exec(state *State) {
	ptr := state.Pointer()
	internal.AsFlow(ptr).Reset()
	if trk := state.Tracker(); trk != nil {
		trk.Reset()
	}
	e.compute(state.Pointer())
}

//NewExecution creates a new execution
func NewExecution(compute internal.Compute) *Executor {
	return &Executor{compute: compute}
}
