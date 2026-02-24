package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"reflect"
	"unsafe"
)

var intSliceType = reflect.TypeOf([]int{})
var intType = reflect.TypeOf(0)
var intSize = unsafe.Sizeof(int(0))

type sliceGetInt struct {
	sliceOp *exec.Operand
	indexOp *exec.Operand
	value   int
}

func (s *sliceGetInt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	slicePtr := s.sliceOp.Compute(ptr)
	idx := *(*int)(s.indexOp.Compute(ptr))
	header := (*reflect.SliceHeader)(slicePtr)
	if idx < 0 || idx >= header.Len {
		panic(fmt.Sprintf("index out of range [%d] with length %d", idx, header.Len))
	}
	itemPtr := unsafe.Pointer(header.Data + uintptr(idx)*intSize)
	s.value = *(*int)(itemPtr)
	return unsafe.Pointer(&s.value)
}

// NewSliceGetInt creates []int index evaluator for expression slice[index].
func NewSliceGetInt(slice, index *Operand) (New, reflect.Type, error) {
	if slice == nil || slice.Type == nil || slice.Type.Type() != intSliceType {
		return nil, nil, fmt.Errorf("unsupported slice get type")
	}
	if index == nil || index.Type == nil || index.Type.Type() != intType {
		return nil, nil, fmt.Errorf("unsupported slice index type")
	}
	return func(control *Control) (internal.Compute, error) {
		sliceOp, err := slice.NewOperand(control)
		if err != nil {
			return nil, err
		}
		indexOp, err := index.NewOperand(control)
		if err != nil {
			return nil, err
		}
		getter := &sliceGetInt{sliceOp: sliceOp, indexOp: indexOp}
		return getter.compute, nil
	}, intType, nil
}

type sliceSetInt struct {
	sliceOp *exec.Operand
	indexOp *exec.Operand
	valueOp *exec.Operand
}

func (s *sliceSetInt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	slicePtr := s.sliceOp.Compute(ptr)
	idx := *(*int)(s.indexOp.Compute(ptr))
	header := (*reflect.SliceHeader)(slicePtr)
	if idx < 0 || idx >= header.Len {
		panic(fmt.Sprintf("index out of range [%d] with length %d", idx, header.Len))
	}
	itemPtr := unsafe.Pointer(header.Data + uintptr(idx)*intSize)
	*(*int)(itemPtr) = *(*int)(s.valueOp.Compute(ptr))
	return nil
}

// NewSliceSetInt creates []int index evaluator for assignment slice[index] = value.
func NewSliceSetInt(slice, index, value *Operand) (New, error) {
	if slice == nil || slice.Type == nil || slice.Type.Type() != intSliceType {
		return nil, fmt.Errorf("unsupported slice set type")
	}
	if index == nil || index.Type == nil || index.Type.Type() != intType {
		return nil, fmt.Errorf("unsupported slice index type")
	}
	if value == nil || value.Type == nil || value.Type.Type() != intType {
		return nil, fmt.Errorf("unsupported slice value type")
	}
	return func(control *Control) (internal.Compute, error) {
		sliceOp, err := slice.NewOperand(control)
		if err != nil {
			return nil, err
		}
		indexOp, err := index.NewOperand(control)
		if err != nil {
			return nil, err
		}
		valueOp, err := value.NewOperand(control)
		if err != nil {
			return nil, err
		}
		setter := &sliceSetInt{sliceOp: sliceOp, indexOp: indexOp, valueOp: valueOp}
		return setter.compute, nil
	}, nil
}
