package signature

import (
	"github.com/viant/igo/exec"
)

type adapter struct {
	*exec.Executor
}

func (a *adapter) fnV() {
	state := a.Executor.NewState()
	defer state.Release()
	a.Exec(state)

}

func (a *adapter) iiiFn(x, y int) int {
	state := a.Executor.NewState()
	defer state.Release()
	ptr := state.Pointer()
	a.InAt(0).SetInt(ptr, x)
	a.InAt(1).SetInt(ptr, y)
	a.Exec(state)
	return a.OutAt(0).Int(ptr)
}

func (a *adapter) f64f64f64Fn(x, y float64) float64 {
	state := a.Executor.NewState()
	defer state.Release()
	ptr := state.Pointer()
	a.InAt(0).SetFloat64(ptr, x)
	a.InAt(1).SetFloat64(ptr, y)
	a.Exec(state)
	return a.OutAt(0).Float64(ptr)
}

func (a *adapter) sssFn(x, y string) string {
	state := a.Executor.NewState()
	defer state.Release()
	ptr := state.Pointer()
	a.InAt(0).SetString(ptr, x)
	a.InAt(1).SetString(ptr, y)
	a.Exec(state)
	return a.OutAt(0).String(ptr)
}
