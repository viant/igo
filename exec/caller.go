package exec

import (
	"unsafe"
)

//Caller represents a caller interface
type Caller interface {
	Call(ptr unsafe.Pointer, args []*Operand) unsafe.Pointer
}
