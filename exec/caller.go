package exec

import (
	"unsafe"
)

//Caller converts internal call node to actual func call
type Caller interface {
	Call(ptr unsafe.Pointer, args []*Operand) unsafe.Pointer
}

//Func converts executor state to actual function
type Func interface {
	New(exec *Executor) interface{}
}
