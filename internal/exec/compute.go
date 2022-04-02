package exec

import (
	"unsafe"
)

//Compute computes expression or run a statement, result pointer holds expression result or statement flow flag
type Compute func(ptr unsafe.Pointer) unsafe.Pointer
