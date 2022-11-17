package signature

import (
	"github.com/viant/igo/exec"
	"github.com/viant/xunsafe"
	"unsafe"
)

type fnV func()

func (f fnV) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	f()
	return nil
}

func (f fnV) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.fnV
}

type iiiFn func(int, int) int

func (f iiiFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsInt(args[0].Compute(ptr))
	y := xunsafe.AsInt(args[1].Compute(ptr))
	z := f(x, y)
	return unsafe.Pointer(&z)
}

func (f iiiFn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.iiiFn
}

type f64f64f64Fn func(float64, float64) float64

func (f f64f64f64Fn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsFloat64(args[0].Compute(ptr))
	y := xunsafe.AsFloat64(args[1].Compute(ptr))
	z := f(x, y)
	return unsafe.Pointer(&z)
}

func (f f64f64f64Fn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.f64f64f64Fn
}

type sssFn func(string, string) string

func (f sssFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsString(args[0].Compute(ptr))
	y := xunsafe.AsString(args[1].Compute(ptr))
	z := f(x, y)
	return unsafe.Pointer(&z)
}

func (f sssFn) New(e *exec.Executor, initPairs ...interface{}) interface{} {
	var adapter = &adapter{e}
	return adapter.sssFn
}

type svrFn func(string, ...interface{})

func (f svrFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsString(args[0].Compute(ptr))
	switch len(args) - 1 { //avoid mem allocation with upto 3 var args
	case 0:
		f(x)
	case 1:
		a1 := args[1].Interface(args[1].Compute(ptr))
		f(x, a1)
	case 2:
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		f(x, a1, a2)
	case 3:
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		a3 := args[3].Interface(args[3].Compute(ptr))
		f(x, a1, a2, a3)
	default:
		var y = make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			y[i] = args[i].Interface(args[i].Compute(ptr))
		}
		f(x, y...)
	}
	return nil
}

type svrieFn func(string, ...interface{}) (int, error)

func (f svrieFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	var i int
	var err error
	x := xunsafe.AsString(args[0].Compute(ptr))

	switch len(args) - 1 { //avoid mem allocation with upto 3 var args
	case 0:
		i, err = f(x)
	case 1:
		a1 := args[1].Interface(args[1].Compute(ptr))
		i, err = f(x, a1)
	case 2:
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		i, err = f(x, a1, a2)
	case 3:
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		a3 := args[3].Interface(args[3].Compute(ptr))
		i, err = f(x, a1, a2, a3)
	default:
		var y = make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			y[i] = args[i].Interface(args[i].Compute(ptr))
		}
		i, err = f(x, y...)
	}
	ret := [2]unsafe.Pointer{
		unsafe.Pointer(&i),
		unsafe.Pointer(&err),
	}
	return unsafe.Pointer(&ret)
}

type vrieFn func(...interface{}) (int, error)

func (f vrieFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	v := 0
	i := &v
	var vErr error
	err := &vErr
	switch len(args) - 1 { //avoid mem allocation with upto 3 var args
	case 0:
		*i, *err = f()
	case 1:
		a0 := args[0].Interface(args[0].Compute(ptr))
		a1 := args[1].Interface(args[1].Compute(ptr))
		*i, *err = f(a0, a1)
	case 2:
		a0 := args[0].Interface(args[0].Compute(ptr))
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		*i, *err = f(a0, a1, a2)
	case 3:
		a0 := args[0].Interface(args[0].Compute(ptr))
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		a3 := args[3].Interface(args[3].Compute(ptr))
		*i, *err = f(a0, a1, a2, a3)
	default:
		var y = make([]interface{}, len(args))
		for i := 0; i < len(args); i++ {
			y[i] = args[i].Interface(args[i].Compute(ptr))
		}
		*i, *err = f(y...)
	}
	var ret = &[2]unsafe.Pointer{}
	ret[0] = unsafe.Pointer(i)
	ret[1] = unsafe.Pointer(err)
	return unsafe.Pointer(ret)
}

type svrs func(string, ...interface{}) string

func (f svrs) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	result := ""
	x := xunsafe.AsString(args[0].Compute(ptr))
	switch len(args) - 1 { //avoid mem allocation with upto 3 var args
	case 0:
		f(x)
	case 1:
		a1 := args[1].Interface(args[1].Compute(ptr))
		result = f(x, a1)
	case 2:
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		result = f(x, a1, a2)
	case 3:
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		a3 := args[3].Interface(args[3].Compute(ptr))
		result = f(x, a1, a2, a3)
	default:
		var y = make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			y[i] = args[i].Interface(args[i].Compute(ptr))
		}
		result = f(x, y...)
	}
	return unsafe.Pointer(&result)
}

type viFn func(interface{}) int

func (f viFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := args[0].Interface(args[0].Compute(ptr))
	r := f(x)
	return unsafe.Pointer(&r)
}

type vf32Fn func(interface{}) float32

func (f vf32Fn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := args[0].Interface(args[0].Compute(ptr))
	r := f(x)
	return unsafe.Pointer(&r)
}

type vf64Fn func(interface{}) float64

func (f vf64Fn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := args[0].Interface(args[0].Compute(ptr))
	r := f(x)
	return unsafe.Pointer(&r)

}

type vbFn func(interface{}) bool

func (f vbFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := args[0].Interface(args[0].Compute(ptr))
	r := f(x)
	return unsafe.Pointer(&r)

}

type vsFn func(interface{}) string

func (f vsFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := args[0].Interface(args[0].Compute(ptr))
	r := f(x)
	return unsafe.Pointer(&r)
}

type vvsvFn func(interface{}, ...interface{}) interface{}

func (f vvsvFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	var result interface{}
	x := args[0].Interface(args[0].Compute(ptr))
	switch len(args) - 1 { //avoid mem allocation with upto 3 var args
	case 0:
		f(x)
	case 1:
		a1 := args[1].Interface(args[1].Compute(ptr))
		result = f(x, a1)
	case 2:
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		result = f(x, a1, a2)
	case 3:
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		a3 := args[3].Interface(args[3].Compute(ptr))
		result = f(x, a1, a2, a3)
	default:
		var y = make([]interface{}, len(args)-1)
		for i := 1; i < len(args); i++ {
			y[i] = args[i].Interface(args[i].Compute(ptr))
		}
		result = f(x, y...)
	}
	return xunsafe.AsPointer(result)
}
