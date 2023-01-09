package et

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

//Operand defines operand
type Operand struct {
	Key string
	Idx int
	*exec.Selector
	Type     *xunsafe.Type
	Value    interface{}
	ValuePtr unsafe.Pointer
	New
}

//Compute computes operand
func (o *Operand) Compute(control *Control) (internal.Compute, *exec.Selector, error) {
	if o.New != nil {
		comp, err := o.New(control)
		return comp, o.Selector, err
	}
	if o.ValuePtr != nil {
		return nil, nil, nil
	}
	if o.Selector != nil {
		return nil, o.Selector, nil
	}
	return nil, nil, nil
}

func (o *Operand) IsNilSelector() bool {
	if sel := o.Selector; sel != nil && sel.ID == "nil" {
		return true
	}
	return false
}

//NewOperand create an operand
func (o *Operand) NewOperand(control *Control) (*exec.Operand, error) {
	var err error
	offset := o.Offset()
	if o.Selector != nil && o.Pathway != exec.PathwayDirect {
		offset = 0
	}
	ptr := o.ValuePtr
	xType := o.Type
	var compute internal.Compute
	var selector *exec.Selector
	if compute, selector, err = o.Compute(control); err != nil {
		return nil, err
	}
	result := exec.NewOperand(xType, offset, compute, selector, ptr)
	return result, result.Validate()
}

//NewOperand crates a new operand
func NewOperand(sel *exec.Selector, oType reflect.Type, newFn New, value interface{}) *Operand {
	result := &Operand{Selector: sel, New: newFn, Value: value}
	if oType == nil && sel != nil {
		oType = sel.Type
	}
	if oType == nil && value != nil {
		oType = reflect.TypeOf(value)
	}
	if oType != nil {
		result.Type = xunsafe.NewType(oType)
	}
	if value != nil {
		result.ValuePtr = result.Type.Pointer(value)
	}
	return result
}

func (o *Operand) ensureNilType(y *Operand) {
	if sel := o.Selector; sel != nil && sel.ID == "nil" {
		if o.Type == nil {
			o.Type = y.Type
		}
		var ptr = unsafe.Pointer(nil)
		o.ValuePtr = unsafe.Pointer(&ptr)

		o.Selector = nil
	}
}
