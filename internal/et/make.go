package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type makeSlice struct {
	len      *exec.Operand
	cap      *exec.Operand
	dest     *exec.Selector
	destType reflect.Type
}

//Compute computes make stmt
func (m *makeSlice) Compute(ptr unsafe.Pointer) unsafe.Pointer {
	sliceLen := *(*int)(m.len.Compute(ptr))
	sliceCap := sliceLen
	if m.cap != nil {
		sliceCap = *(*int)(m.cap.Compute(ptr))
	}
	aSlice := reflect.MakeSlice(m.destType, sliceLen, sliceCap)
	return xunsafe.ValuePointer(&aSlice)
}

//NewMake creates a make stmt
func NewMake(destType reflect.Type, args Operands) (New, reflect.Type, error) {
	switch destType.Kind() {
	case reflect.Slice, reflect.Chan:
	default:
		return nil, nil, fmt.Errorf("invalid make type: %v", destType.String())
	}
	return func(control *Control) (internal.Compute, error) {
		operands, err := args.operands(control)
		if err != nil {
			return nil, err
		}
		result := &makeSlice{destType: destType, len: operands[0]}
		if len(operands) > 1 {
			result.len = operands[1]
		}
		return result.Compute, nil
	}, destType, nil
}
