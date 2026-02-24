package exec

import (
	"unsafe"
)

// Caller converts internal call node to actual func call
type Caller interface {
	Call(ptr unsafe.Pointer, args []*Operand) unsafe.Pointer
}

// InterfaceCaller supports direct interface{} return without pointer reinterpretation.
type InterfaceCaller interface {
	Caller
	CallInterface(ptr unsafe.Pointer, args []*Operand) interface{}
}

// VoidCaller optionally optimizes call sites where return values are ignored.
type VoidCaller interface {
	Caller
	CallVoid(ptr unsafe.Pointer, args []*Operand)
}

// Func converts executor state to actual function
type Func interface {
	New(exec *Executor) interface{}
}
