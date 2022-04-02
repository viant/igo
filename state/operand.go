package state

import (
	"fmt"
	"github.com/viant/xunsafe"
	"unsafe"
)

//Operand represents typed compute
type Operand struct {
	*xunsafe.Type //result pointer type
	Offset        uintptr
	compute       func(ptr unsafe.Pointer) unsafe.Pointer
	selector      *Selector
	useSelector   bool
	ptr           unsafe.Pointer
}

//Validate checks if operand is valid
func (o *Operand) Validate() error {
	if o.Offset == 0 && o.ptr == nil && o.compute == nil && o.selector == nil {

		return fmt.Errorf("invalid operand")
	}
	return nil
}

func (o *Operand) computePtr(ptr unsafe.Pointer) unsafe.Pointer {
	x := o.ptr
	if o.Offset != 0 {
		x = unsafe.Pointer(uintptr(ptr) + o.Offset)
	}
	return x
}

//SetValuePtr set value ptr
func (o *Operand) SetValuePtr(ptr unsafe.Pointer, src unsafe.Pointer) {
	value := o.Type.Interface(src)
	oPtr := unsafe.Pointer(uintptr(o.Compute(ptr)) - o.selector.Field.Offset)
	o.selector.SetValue(oPtr, value)
}

//SetValue sets value
func (o *Operand) SetValue(ptr unsafe.Pointer, src interface{}) {
	oPtr := unsafe.Pointer(uintptr(o.Compute(ptr)) - o.selector.Field.Offset)
	o.selector.SetValue(oPtr, src)
}

//Compute compute operan value
func (o *Operand) Compute(ptr unsafe.Pointer) unsafe.Pointer {
	x := o.computePtr(ptr)
	if o.useSelector {
		x = o.selector.Addr(ptr)
		return x
	}
	if o.compute != nil {
		x = o.compute(ptr)
	}
	return x
}

//NewOperand creates an operand
func NewOperand(xType *xunsafe.Type, offset uintptr, compute func(ptr unsafe.Pointer) unsafe.Pointer, selector *Selector, ptr unsafe.Pointer) *Operand {
	return &Operand{
		Type:        xType,
		Offset:      offset,
		compute:     compute,
		selector:    selector,
		useSelector: offset == 0 && selector != nil,
		ptr:         ptr,
	}
}
