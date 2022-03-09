package exec

//Executor abstraction holding execution plan (execution syntax tree) to execute it
type Executor struct {
	compute Compute
	In      []string //params variable identity
	Out     []string //result variable identity
}

//Exec executes execution plan
func (e *Executor) Exec(vars *Variables) {
	ptr := vars.Pointer()
	AsFlow(ptr).Reset()
	e.compute(vars.Pointer())
}

//NewExecution creates a new execution
func NewExecution(compute Compute) *Executor {
	return &Executor{compute: compute}
}
