package plan

import (
	"github.com/viant/igo/exec"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type caller struct {
	fn reflect.Value
}

func (c *caller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	var params = make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		vPtr := args[i].Compute(ptr)
		v := args[i].Interface(vPtr)
		params[i] = reflect.ValueOf(v)
	}
	out := c.fn.Call(params)
	switch len(out) {
	case 0:
		return nil
	case 1:
		return xunsafe.ValuePointer(&out[0])
	case 2:
		ret := [2]unsafe.Pointer{
			xunsafe.ValuePointer(&out[0]),
			xunsafe.ValuePointer(&out[1]),
		}
		return unsafe.Pointer(&ret)
	case 3:
		ret := [3]unsafe.Pointer{
			xunsafe.ValuePointer(&out[0]),
			xunsafe.ValuePointer(&out[1]),
			xunsafe.ValuePointer(&out[2]),
		}
		return unsafe.Pointer(&ret)
	case 4:
		ret := [4]unsafe.Pointer{
			xunsafe.ValuePointer(&out[0]),
			xunsafe.ValuePointer(&out[1]),
			xunsafe.ValuePointer(&out[2]),
			xunsafe.ValuePointer(&out[3]),
		}
		return unsafe.Pointer(&ret)
	case 5:
		ret := [5]unsafe.Pointer{
			xunsafe.ValuePointer(&out[0]),
			xunsafe.ValuePointer(&out[1]),
			xunsafe.ValuePointer(&out[2]),
			xunsafe.ValuePointer(&out[3]),
			xunsafe.ValuePointer(&out[4]),
		}
		return unsafe.Pointer(&ret)
	}
	panic("too many return parameters")
}

func newCaller(fn interface{}) exec.Caller {
	result := &caller{
		fn: reflect.ValueOf(fn),
	}
	return result
}
