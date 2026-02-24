package signature

import (
	"github.com/viant/igo/exec"
	"github.com/viant/xunsafe"
	"reflect"
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

type iisbFn func(int, []int) bool

func (f iisbFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsInt(args[0].Compute(ptr))
	var callArgs []int
	if len(args) == 2 {
		if args[1].Kind() == reflect.Int {
			callArgs = append(callArgs, xunsafe.AsInt(args[1].Compute(ptr)))
		} else {
			callArgs = xunsafe.AsInts(args[1].Compute(ptr))
		}
	} else {
		for i := 1; i < len(args); i++ {
			callArgs = append(callArgs, xunsafe.AsInt(args[i].Compute(ptr)))
		}
	}
	z := f(x, callArgs)
	return unsafe.Pointer(&z)
}

func (f iisbFn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.iisbFn
}

type sssbFn func(string, []string) bool

func (f sssbFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsString(args[0].Compute(ptr))
	var callArgs []string
	if len(args) == 2 {
		if args[1].Kind() == reflect.String {
			callArgs = append(callArgs, xunsafe.AsString(args[1].Compute(ptr)))
		} else {
			callArgs = xunsafe.AsStrings(args[1].Compute(ptr))
		}
	} else {
		for i := 1; i < len(args); i++ {
			callArgs = append(callArgs, xunsafe.AsString(args[i].Compute(ptr)))
		}
	}
	z := f(x, callArgs)
	return unsafe.Pointer(&z)
}

func (f sssbFn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.sssbFn
}

type ssbFn func(string, string) bool

func (f ssbFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsString(args[0].Compute(ptr))
	y := xunsafe.AsString(args[1].Compute(ptr))
	z := f(x, y)
	return unsafe.Pointer(&z)
}

func (f ssbFn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.ssbFn
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

type iiFn func(int) int

func (f iiFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsInt(args[0].Compute(ptr))
	z := f(x)
	return unsafe.Pointer(&z)
}

func (f iiFn) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	x := xunsafe.AsInt(args[0].Compute(ptr))
	f(x)
}

func (f iiFn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.iiFn
}

type ivFn func(int)

func (f ivFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsInt(args[0].Compute(ptr))
	f(x)
	return nil
}

func (f ivFn) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	x := xunsafe.AsInt(args[0].Compute(ptr))
	f(x)
}

func (f ivFn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.ivFn
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
			y[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		f(x, y...)
	}
	return nil
}

func (f svrFn) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	x := xunsafe.AsString(args[0].Compute(ptr))
	switch len(args) - 1 {
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
			y[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		f(x, y...)
	}
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
			y[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		i, err = f(x, y...)
	}
	ret := [2]unsafe.Pointer{
		unsafe.Pointer(&i),
		unsafe.Pointer(&err),
	}
	return unsafe.Pointer(&ret)
}

func (f svrieFn) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	x := xunsafe.AsString(args[0].Compute(ptr))
	switch len(args) - 1 {
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
			y[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		f(x, y...)
	}
}

type vrieFn func(...interface{}) (int, error)

func (f vrieFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	v := 0
	i := &v
	var vErr error
	err := &vErr
	switch len(args) { //avoid mem allocation with upto 3 var args
	case 0:
		*i, *err = f()
	case 1:
		a0 := args[0].Interface(args[0].Compute(ptr))
		callArgs := [1]interface{}{a0}
		*i, *err = f(callArgs[:]...)
	case 2:
		a0 := args[0].Interface(args[0].Compute(ptr))
		a1 := args[1].Interface(args[1].Compute(ptr))
		callArgs := [2]interface{}{a0, a1}
		*i, *err = f(callArgs[:]...)
	case 3:
		a0 := args[0].Interface(args[0].Compute(ptr))
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		callArgs := [3]interface{}{a0, a1, a2}
		*i, *err = f(callArgs[:]...)
	case 4:
		a0 := args[0].Interface(args[0].Compute(ptr))
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		a3 := args[3].Interface(args[3].Compute(ptr))
		callArgs := [4]interface{}{a0, a1, a2, a3}
		*i, *err = f(callArgs[:]...)
	default:
		var y = make([]interface{}, len(args))
		for i := 0; i < len(args); i++ {
			y[i] = args[i].Interface(args[i].Compute(ptr))
		}
		*i, *err = f(y...)
	}
	ret := [2]unsafe.Pointer{
		unsafe.Pointer(i),
		unsafe.Pointer(err),
	}
	return unsafe.Pointer(&ret)
}

func (f vrieFn) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	switch len(args) {
	case 0:
		f()
	case 1:
		a0 := args[0].Interface(args[0].Compute(ptr))
		callArgs := [1]interface{}{a0}
		f(callArgs[:]...)
	case 2:
		a0 := args[0].Interface(args[0].Compute(ptr))
		a1 := args[1].Interface(args[1].Compute(ptr))
		callArgs := [2]interface{}{a0, a1}
		f(callArgs[:]...)
	case 3:
		a0 := args[0].Interface(args[0].Compute(ptr))
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		callArgs := [3]interface{}{a0, a1, a2}
		f(callArgs[:]...)
	case 4:
		a0 := args[0].Interface(args[0].Compute(ptr))
		a1 := args[1].Interface(args[1].Compute(ptr))
		a2 := args[2].Interface(args[2].Compute(ptr))
		a3 := args[3].Interface(args[3].Compute(ptr))
		callArgs := [4]interface{}{a0, a1, a2, a3}
		f(callArgs[:]...)
	default:
		var y = make([]interface{}, len(args))
		for i := 0; i < len(args); i++ {
			y[i] = args[i].Interface(args[i].Compute(ptr))
		}
		f(y...)
	}
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
			y[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		result = f(x, y...)
	}
	return unsafe.Pointer(&result)
}

func (f svrs) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	x := xunsafe.AsString(args[0].Compute(ptr))
	switch len(args) - 1 {
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
			y[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		f(x, y...)
	}
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

type vi64Fn func(interface{}) int64

func (f vi64Fn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := args[0].Interface(args[0].Compute(ptr))
	r := f(x)
	return unsafe.Pointer(&r)
}

func (f vi64Fn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.vi64Fn
}

type iibFn func(int, int) bool

func (f iibFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsInt(args[0].Compute(ptr))
	y := xunsafe.AsInt(args[1].Compute(ptr))
	z := f(x, y)
	return unsafe.Pointer(&z)
}

func (f iibFn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.iibFn
}

type ssiFn func(string, string) int

func (f ssiFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsString(args[0].Compute(ptr))
	y := xunsafe.AsString(args[1].Compute(ptr))
	z := f(x, y)
	return unsafe.Pointer(&z)
}

func (f ssiFn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.ssiFn
}

type i64i64i64Fn func(int64, int64) int64

func (f i64i64i64Fn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsInt64(args[0].Compute(ptr))
	y := xunsafe.AsInt64(args[1].Compute(ptr))
	z := f(x, y)
	return unsafe.Pointer(&z)
}

func (f i64i64i64Fn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.i64i64i64Fn
}

type i64i64bFn func(int64, int64) bool

func (f i64i64bFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	x := xunsafe.AsInt64(args[0].Compute(ptr))
	y := xunsafe.AsInt64(args[1].Compute(ptr))
	z := f(x, y)
	return unsafe.Pointer(&z)
}

func (f i64i64bFn) New(e *exec.Executor) interface{} {
	var adapter = &adapter{Executor: e}
	return adapter.i64i64bFn
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
			y[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		result = f(x, y...)
	}
	return xunsafe.AsPointer(result)
}

func (f vvsvFn) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	x := args[0].Interface(args[0].Compute(ptr))
	switch len(args) - 1 {
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
			y[i-1] = args[i].Interface(args[i].Compute(ptr))
		}
		f(x, y...)
	}
}
