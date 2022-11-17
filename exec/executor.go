package exec

import (
	"github.com/viant/igo/internal"
	"reflect"
	"sync"
)

//Executor abstraction holding execution plan (execution syntax tree) to execute it
type Executor struct {
	compute internal.Compute
	In      []string //params variable identity
	Out     []string //result variable identity
	in      []*Selector
	out     []*Selector
	Init    []interface{}
	fn      interface{}
	fnType  reflect.Type
	pool    *sync.Pool
}

//InAt returns in selector at position i
func (e *Executor) InAt(i int) *Selector {
	return e.in[i]
}

//OutAt returns out selector at position i
func (e *Executor) OutAt(i int) *Selector {
	return e.out[i]
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

func (e *Executor) NewState() *State {
	return e.pool.Get().(*State)
}

func (e *Executor) Func() (interface{}, error) {
	if e.fn != nil {
		return e.fn, nil
	}
	fnType := e.Signature()
	fnVal := reflect.New(fnType).Elem().Interface()
	if val := AsFunc(fnVal); val != nil {
		if fn, ok := val.(Func); ok {
			e.fn = fn.New(e)
			return e.fn, nil
		}
	}
	e.fn = reflect.MakeFunc(fnType, func(args []reflect.Value) (results []reflect.Value) {
		state := e.NewState()
		defer state.Release()
		for i, in := range e.In {
			if err := state.SetValue(in, args[i].Interface()); err != nil {
				panic(err)
			}
		}
		e.Exec(state)
		results = make([]reflect.Value, len(e.out))
		for i, out := range e.Out {
			value, err := state.Value(out)
			if err != nil {
				panic(err)
			}
			results[i] = reflect.ValueOf(value)
		}
		return results
	}).Interface()

	return e.fn, nil
}

func (e *Executor) Signature() reflect.Type {
	if e.fnType != nil {
		return e.fnType
	}
	var in, out []reflect.Type
	for i := range e.In {
		in = append(in, e.in[i].Type)
	}
	for i := range e.Out {
		out = append(out, e.out[i].Type)
	}
	e.fnType = reflect.FuncOf(in, out, false)
	return e.fnType
}

//NewExecution creates a new execution
func NewExecution(compute internal.Compute, pool *sync.Pool, in []*Selector, out []*Selector) *Executor {
	return &Executor{compute: compute, pool: pool, in: in, out: out,
		In:  Selectors(in).IDs(),
		Out: Selectors(out).IDs(),
	}
}
