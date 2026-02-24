package plan

import (
	"github.com/viant/igo/exec"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type compiledInlineFunc struct {
	caller   exec.Caller
	retTypes []reflect.Type
	fnType   reflect.Type
}

type inlineExecutorCaller struct {
	execution *exec.Executor
	state     *exec.State
	in        []*exec.Selector
	out       []*exec.Selector
	ret2      [2]unsafe.Pointer
	ret3      [3]unsafe.Pointer
	ret4      [4]unsafe.Pointer
	ret5      [5]unsafe.Pointer
	ret6      [6]unsafe.Pointer
}

func newInlineExecutorCaller(execution *exec.Executor) *inlineExecutorCaller {
	in := make([]*exec.Selector, len(execution.In))
	for i := range execution.In {
		in[i] = execution.InAt(i)
	}
	out := make([]*exec.Selector, len(execution.Out))
	for i := range execution.Out {
		out[i] = execution.OutAt(i)
	}
	return &inlineExecutorCaller{
		execution: execution,
		state:     execution.NewState(),
		in:        in,
		out:       out,
	}
}

func (c *inlineExecutorCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	if len(args) != len(c.in) {
		panic("invalid number of inline function arguments")
	}
	statePtr := c.state.Pointer()
	for i := range args {
		src := args[i].Compute(ptr)
		destSel := c.in[i]
		destPtr := destSel.Addr(statePtr)
		if argType := args[i].Type; argType != nil && argType.Type() == destSel.Type && destSel.Kind() != reflect.Interface {
			xunsafe.Copy(destPtr, src, int(destSel.Size()))
			continue
		}
		destSel.SetValue(destSel.Upstream(statePtr), args[i].Interface(src))
	}
	c.execution.Exec(c.state)

	switch len(c.out) {
	case 0:
		return nil
	case 1:
		return c.out[0].Addr(statePtr)
	case 2:
		c.ret2[0] = c.out[0].Addr(statePtr)
		c.ret2[1] = c.out[1].Addr(statePtr)
		return unsafe.Pointer(&c.ret2)
	case 3:
		c.ret3[0] = c.out[0].Addr(statePtr)
		c.ret3[1] = c.out[1].Addr(statePtr)
		c.ret3[2] = c.out[2].Addr(statePtr)
		return unsafe.Pointer(&c.ret3)
	case 4:
		c.ret4[0] = c.out[0].Addr(statePtr)
		c.ret4[1] = c.out[1].Addr(statePtr)
		c.ret4[2] = c.out[2].Addr(statePtr)
		c.ret4[3] = c.out[3].Addr(statePtr)
		return unsafe.Pointer(&c.ret4)
	case 5:
		c.ret5[0] = c.out[0].Addr(statePtr)
		c.ret5[1] = c.out[1].Addr(statePtr)
		c.ret5[2] = c.out[2].Addr(statePtr)
		c.ret5[3] = c.out[3].Addr(statePtr)
		c.ret5[4] = c.out[4].Addr(statePtr)
		return unsafe.Pointer(&c.ret5)
	case 6:
		c.ret6[0] = c.out[0].Addr(statePtr)
		c.ret6[1] = c.out[1].Addr(statePtr)
		c.ret6[2] = c.out[2].Addr(statePtr)
		c.ret6[3] = c.out[3].Addr(statePtr)
		c.ret6[4] = c.out[4].Addr(statePtr)
		c.ret6[5] = c.out[5].Addr(statePtr)
		return unsafe.Pointer(&c.ret6)
	}
	panic("too many inline return parameters")
}
