package plan

import (
	"github.com/viant/igo/exec"
	"unsafe"
)

type variadicInterfaceIntErrCaller struct {
	fn  func(...interface{}) (int, error)
	a1  [1]interface{}
	a2  [2]interface{}
	a3  [3]interface{}
	a4  [4]interface{}
	i   int
	err error
	ret [2]unsafe.Pointer
}

func newVariadicInterfaceIntErrCaller(fn func(...interface{}) (int, error)) exec.Caller {
	caller := &variadicInterfaceIntErrCaller{fn: fn}
	caller.ret[0] = unsafe.Pointer(&caller.i)
	caller.ret[1] = unsafe.Pointer(&caller.err)
	return caller
}

func (c *variadicInterfaceIntErrCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	switch len(args) {
	case 0:
		c.i, c.err = c.fn()
	case 1:
		c.a1[0] = args[0].Interface(args[0].Compute(ptr))
		c.i, c.err = c.fn(c.a1[:]...)
	case 2:
		c.a2[0] = args[0].Interface(args[0].Compute(ptr))
		c.a2[1] = args[1].Interface(args[1].Compute(ptr))
		c.i, c.err = c.fn(c.a2[:]...)
	case 3:
		c.a3[0] = args[0].Interface(args[0].Compute(ptr))
		c.a3[1] = args[1].Interface(args[1].Compute(ptr))
		c.a3[2] = args[2].Interface(args[2].Compute(ptr))
		c.i, c.err = c.fn(c.a3[:]...)
	case 4:
		c.a4[0] = args[0].Interface(args[0].Compute(ptr))
		c.a4[1] = args[1].Interface(args[1].Compute(ptr))
		c.a4[2] = args[2].Interface(args[2].Compute(ptr))
		c.a4[3] = args[3].Interface(args[3].Compute(ptr))
		c.i, c.err = c.fn(c.a4[:]...)
	default:
		callArgs := make([]interface{}, len(args))
		for i := range args {
			callArgs[i] = args[i].Interface(args[i].Compute(ptr))
		}
		c.i, c.err = c.fn(callArgs...)
	}
	return unsafe.Pointer(&c.ret)
}

func (c *variadicInterfaceIntErrCaller) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	switch len(args) {
	case 0:
		c.fn()
	case 1:
		c.a1[0] = args[0].Interface(args[0].Compute(ptr))
		c.fn(c.a1[:]...)
	case 2:
		c.a2[0] = args[0].Interface(args[0].Compute(ptr))
		c.a2[1] = args[1].Interface(args[1].Compute(ptr))
		c.fn(c.a2[:]...)
	case 3:
		c.a3[0] = args[0].Interface(args[0].Compute(ptr))
		c.a3[1] = args[1].Interface(args[1].Compute(ptr))
		c.a3[2] = args[2].Interface(args[2].Compute(ptr))
		c.fn(c.a3[:]...)
	case 4:
		c.a4[0] = args[0].Interface(args[0].Compute(ptr))
		c.a4[1] = args[1].Interface(args[1].Compute(ptr))
		c.a4[2] = args[2].Interface(args[2].Compute(ptr))
		c.a4[3] = args[3].Interface(args[3].Compute(ptr))
		c.fn(c.a4[:]...)
	default:
		callArgs := make([]interface{}, len(args))
		for i := range args {
			callArgs[i] = args[i].Interface(args[i].Compute(ptr))
		}
		c.fn(callArgs...)
	}
}

type variadicStringVoidCaller struct {
	fn func(string, ...interface{})
	a1 [1]interface{}
	a2 [2]interface{}
	a3 [3]interface{}
	a4 [4]interface{}
}

func newVariadicStringVoidCaller(fn func(string, ...interface{})) exec.Caller {
	return &variadicStringVoidCaller{fn: fn}
}

func (c *variadicStringVoidCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	c.CallVoid(ptr, args)
	return nil
}

func (c *variadicStringVoidCaller) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	s := *(*string)(args[0].Compute(ptr))
	switch len(args) - 1 {
	case 0:
		c.fn(s)
	case 1:
		c.a1[0] = args[1].Interface(args[1].Compute(ptr))
		c.fn(s, c.a1[:]...)
	case 2:
		c.a2[0] = args[1].Interface(args[1].Compute(ptr))
		c.a2[1] = args[2].Interface(args[2].Compute(ptr))
		c.fn(s, c.a2[:]...)
	case 3:
		c.a3[0] = args[1].Interface(args[1].Compute(ptr))
		c.a3[1] = args[2].Interface(args[2].Compute(ptr))
		c.a3[2] = args[3].Interface(args[3].Compute(ptr))
		c.fn(s, c.a3[:]...)
	case 4:
		c.a4[0] = args[1].Interface(args[1].Compute(ptr))
		c.a4[1] = args[2].Interface(args[2].Compute(ptr))
		c.a4[2] = args[3].Interface(args[3].Compute(ptr))
		c.a4[3] = args[4].Interface(args[4].Compute(ptr))
		c.fn(s, c.a4[:]...)
	default:
		callArgs := make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			callArgs[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		c.fn(s, callArgs...)
	}
}

type variadicStringStringCaller struct {
	fn     func(string, ...interface{}) string
	a1     [1]interface{}
	a2     [2]interface{}
	a3     [3]interface{}
	a4     [4]interface{}
	result string
}

func newVariadicStringStringCaller(fn func(string, ...interface{}) string) exec.Caller {
	return &variadicStringStringCaller{fn: fn}
}

func (c *variadicStringStringCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	s := *(*string)(args[0].Compute(ptr))
	switch len(args) - 1 {
	case 0:
		c.result = c.fn(s)
	case 1:
		c.a1[0] = args[1].Interface(args[1].Compute(ptr))
		c.result = c.fn(s, c.a1[:]...)
	case 2:
		c.a2[0] = args[1].Interface(args[1].Compute(ptr))
		c.a2[1] = args[2].Interface(args[2].Compute(ptr))
		c.result = c.fn(s, c.a2[:]...)
	case 3:
		c.a3[0] = args[1].Interface(args[1].Compute(ptr))
		c.a3[1] = args[2].Interface(args[2].Compute(ptr))
		c.a3[2] = args[3].Interface(args[3].Compute(ptr))
		c.result = c.fn(s, c.a3[:]...)
	case 4:
		c.a4[0] = args[1].Interface(args[1].Compute(ptr))
		c.a4[1] = args[2].Interface(args[2].Compute(ptr))
		c.a4[2] = args[3].Interface(args[3].Compute(ptr))
		c.a4[3] = args[4].Interface(args[4].Compute(ptr))
		c.result = c.fn(s, c.a4[:]...)
	default:
		callArgs := make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			callArgs[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		c.result = c.fn(s, callArgs...)
	}
	return unsafe.Pointer(&c.result)
}

func (c *variadicStringStringCaller) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	s := *(*string)(args[0].Compute(ptr))
	switch len(args) - 1 {
	case 0:
		c.fn(s)
	case 1:
		c.a1[0] = args[1].Interface(args[1].Compute(ptr))
		c.fn(s, c.a1[:]...)
	case 2:
		c.a2[0] = args[1].Interface(args[1].Compute(ptr))
		c.a2[1] = args[2].Interface(args[2].Compute(ptr))
		c.fn(s, c.a2[:]...)
	case 3:
		c.a3[0] = args[1].Interface(args[1].Compute(ptr))
		c.a3[1] = args[2].Interface(args[2].Compute(ptr))
		c.a3[2] = args[3].Interface(args[3].Compute(ptr))
		c.fn(s, c.a3[:]...)
	case 4:
		c.a4[0] = args[1].Interface(args[1].Compute(ptr))
		c.a4[1] = args[2].Interface(args[2].Compute(ptr))
		c.a4[2] = args[3].Interface(args[3].Compute(ptr))
		c.a4[3] = args[4].Interface(args[4].Compute(ptr))
		c.fn(s, c.a4[:]...)
	default:
		callArgs := make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			callArgs[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		c.fn(s, callArgs...)
	}
}

type variadicStringIntErrCaller struct {
	fn  func(string, ...interface{}) (int, error)
	a1  [1]interface{}
	a2  [2]interface{}
	a3  [3]interface{}
	a4  [4]interface{}
	i   int
	err error
	ret [2]unsafe.Pointer
}

func newVariadicStringIntErrCaller(fn func(string, ...interface{}) (int, error)) exec.Caller {
	caller := &variadicStringIntErrCaller{fn: fn}
	caller.ret[0] = unsafe.Pointer(&caller.i)
	caller.ret[1] = unsafe.Pointer(&caller.err)
	return caller
}

func (c *variadicStringIntErrCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	s := *(*string)(args[0].Compute(ptr))
	switch len(args) - 1 {
	case 0:
		c.i, c.err = c.fn(s)
	case 1:
		c.a1[0] = args[1].Interface(args[1].Compute(ptr))
		c.i, c.err = c.fn(s, c.a1[:]...)
	case 2:
		c.a2[0] = args[1].Interface(args[1].Compute(ptr))
		c.a2[1] = args[2].Interface(args[2].Compute(ptr))
		c.i, c.err = c.fn(s, c.a2[:]...)
	case 3:
		c.a3[0] = args[1].Interface(args[1].Compute(ptr))
		c.a3[1] = args[2].Interface(args[2].Compute(ptr))
		c.a3[2] = args[3].Interface(args[3].Compute(ptr))
		c.i, c.err = c.fn(s, c.a3[:]...)
	case 4:
		c.a4[0] = args[1].Interface(args[1].Compute(ptr))
		c.a4[1] = args[2].Interface(args[2].Compute(ptr))
		c.a4[2] = args[3].Interface(args[3].Compute(ptr))
		c.a4[3] = args[4].Interface(args[4].Compute(ptr))
		c.i, c.err = c.fn(s, c.a4[:]...)
	default:
		callArgs := make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			callArgs[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		c.i, c.err = c.fn(s, callArgs...)
	}
	return unsafe.Pointer(&c.ret)
}

func (c *variadicStringIntErrCaller) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	s := *(*string)(args[0].Compute(ptr))
	switch len(args) - 1 {
	case 0:
		c.fn(s)
	case 1:
		c.a1[0] = args[1].Interface(args[1].Compute(ptr))
		c.fn(s, c.a1[:]...)
	case 2:
		c.a2[0] = args[1].Interface(args[1].Compute(ptr))
		c.a2[1] = args[2].Interface(args[2].Compute(ptr))
		c.fn(s, c.a2[:]...)
	case 3:
		c.a3[0] = args[1].Interface(args[1].Compute(ptr))
		c.a3[1] = args[2].Interface(args[2].Compute(ptr))
		c.a3[2] = args[3].Interface(args[3].Compute(ptr))
		c.fn(s, c.a3[:]...)
	case 4:
		c.a4[0] = args[1].Interface(args[1].Compute(ptr))
		c.a4[1] = args[2].Interface(args[2].Compute(ptr))
		c.a4[2] = args[3].Interface(args[3].Compute(ptr))
		c.a4[3] = args[4].Interface(args[4].Compute(ptr))
		c.fn(s, c.a4[:]...)
	default:
		callArgs := make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			callArgs[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		c.fn(s, callArgs...)
	}
}

type variadicInterfaceInterfaceCaller struct {
	fn     func(interface{}, ...interface{}) interface{}
	a1     [1]interface{}
	a2     [2]interface{}
	a3     [3]interface{}
	a4     [4]interface{}
	result interface{}
}

func newVariadicInterfaceInterfaceCaller(fn func(interface{}, ...interface{}) interface{}) exec.Caller {
	return &variadicInterfaceInterfaceCaller{fn: fn}
}

func (c *variadicInterfaceInterfaceCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	c.result = c.CallInterface(ptr, args)
	return unsafe.Pointer(&c.result)
}

func (c *variadicInterfaceInterfaceCaller) CallInterface(ptr unsafe.Pointer, args []*exec.Operand) interface{} {
	v := args[0].Interface(args[0].Compute(ptr))
	switch len(args) - 1 {
	case 0:
		return c.fn(v)
	case 1:
		c.a1[0] = args[1].Interface(args[1].Compute(ptr))
		return c.fn(v, c.a1[:]...)
	case 2:
		c.a2[0] = args[1].Interface(args[1].Compute(ptr))
		c.a2[1] = args[2].Interface(args[2].Compute(ptr))
		return c.fn(v, c.a2[:]...)
	case 3:
		c.a3[0] = args[1].Interface(args[1].Compute(ptr))
		c.a3[1] = args[2].Interface(args[2].Compute(ptr))
		c.a3[2] = args[3].Interface(args[3].Compute(ptr))
		return c.fn(v, c.a3[:]...)
	case 4:
		c.a4[0] = args[1].Interface(args[1].Compute(ptr))
		c.a4[1] = args[2].Interface(args[2].Compute(ptr))
		c.a4[2] = args[3].Interface(args[3].Compute(ptr))
		c.a4[3] = args[4].Interface(args[4].Compute(ptr))
		return c.fn(v, c.a4[:]...)
	default:
		callArgs := make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			callArgs[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		return c.fn(v, callArgs...)
	}
}

func (c *variadicInterfaceInterfaceCaller) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	v := args[0].Interface(args[0].Compute(ptr))
	switch len(args) - 1 {
	case 0:
		c.fn(v)
	case 1:
		c.a1[0] = args[1].Interface(args[1].Compute(ptr))
		c.fn(v, c.a1[:]...)
	case 2:
		c.a2[0] = args[1].Interface(args[1].Compute(ptr))
		c.a2[1] = args[2].Interface(args[2].Compute(ptr))
		c.fn(v, c.a2[:]...)
	case 3:
		c.a3[0] = args[1].Interface(args[1].Compute(ptr))
		c.a3[1] = args[2].Interface(args[2].Compute(ptr))
		c.a3[2] = args[3].Interface(args[3].Compute(ptr))
		c.fn(v, c.a3[:]...)
	case 4:
		c.a4[0] = args[1].Interface(args[1].Compute(ptr))
		c.a4[1] = args[2].Interface(args[2].Compute(ptr))
		c.a4[2] = args[3].Interface(args[3].Compute(ptr))
		c.a4[3] = args[4].Interface(args[4].Compute(ptr))
		c.fn(v, c.a4[:]...)
	default:
		callArgs := make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			callArgs[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		c.fn(v, callArgs...)
	}
}
