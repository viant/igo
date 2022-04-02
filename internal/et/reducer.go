package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

//NewReducer creates reducer
func NewReducer(x *exec.Selector, params []*exec.Selector, results []*exec.Selector, init *Operand, body New) (New, reflect.Type, error) {
	if len(results) != 1 {
		return nil, nil, fmt.Errorf("invalid reducer signaure, expected return %v type", params[0].Type.String())
	}
	if results[0].Type != params[0].Type {
		return nil, nil, fmt.Errorf("invalid reducer signaure, invaid result type, expected: %s, but had: %s", params[0].Type.String(), results[0].Type.String())
	}
	destType := params[0].Type
	return func(control *Control) (internal.Compute, error) {
		var err error
		result := &reducer{x: x, slice: xunsafe.NewSlice(x.Type), retType: destType.Kind()}
		if init != nil {
			if result.init, err = init.NewOperand(control); err != nil {
				return nil, err
			}
		}
		if result.body, err = body(control); err != nil {
			return nil, err
		}
		result.accOffset = params[0].Offset()
		result.value = params[1]
		if len(params) > 2 {
			result.indexOffset = params[2].Offset()
		}
		result.isComponentPtr = x.Type.Elem().Kind() == reflect.Ptr
		switch destType.Kind() {
		case reflect.Float64:
			return result.computeFloat64, nil
		default:
			return result.computeInt, nil
		}
	}, destType, nil
}

type reducer struct {
	slice          *xunsafe.Slice
	init           *exec.Operand
	body           internal.Compute
	retType        reflect.Kind
	isComponentPtr bool
	x              *exec.Selector
	accOffset      uintptr
	value          *exec.Selector
	hasIndex       bool
	indexOffset    uintptr
}

func (s *reducer) computeFloat64(ptr unsafe.Pointer) unsafe.Pointer {
	slicePtr := s.x.Addr(ptr)
	sliceLen := s.slice.Len(slicePtr)
	acc := 0.0
	if s.init != nil {
		acc = xunsafe.AsFloat64(s.init.Compute(ptr))
	}

	if s.hasIndex {
		for i := 0; i < sliceLen; i++ {
			*(*float64)(unsafe.Pointer(uintptr(ptr) + s.accOffset)) = acc
			item := s.slice.ValuePointerAt(slicePtr, i)
			s.value.SetValue(s.value.Upstream(ptr), item)
			*(*int)(unsafe.Pointer(uintptr(ptr) + s.indexOffset)) = i
			acc = *(*float64)(s.body(ptr))
		}

	} else {
		for i := 0; i < sliceLen; i++ {
			*(*float64)(unsafe.Pointer(uintptr(ptr) + s.accOffset)) = acc
			item := s.slice.ValuePointerAt(slicePtr, i)
			s.value.SetValue(s.value.Upstream(ptr), item)
			acc = *(*float64)(s.body(ptr))
		}
	}
	return unsafe.Pointer(&acc)
}

func (s *reducer) computeInt(ptr unsafe.Pointer) unsafe.Pointer {
	slicePtr := s.x.Addr(ptr)
	sliceLen := s.slice.Len(slicePtr)
	acc := 0
	if s.init != nil {
		acc = xunsafe.AsInt(s.init.Compute(ptr))
	}

	if s.hasIndex {
		for i := 0; i < sliceLen; i++ {
			*(*int)(unsafe.Pointer(uintptr(ptr) + s.accOffset)) = acc
			item := s.slice.ValuePointerAt(slicePtr, i)
			s.value.SetValue(s.value.Upstream(ptr), item)
			*(*int)(unsafe.Pointer(uintptr(ptr) + s.indexOffset)) = i
			acc = *(*int)(s.body(ptr))
		}

	} else {
		for i := 0; i < sliceLen; i++ {
			*(*int)(unsafe.Pointer(uintptr(ptr) + s.accOffset)) = acc
			item := s.slice.ValuePointerAt(slicePtr, i)
			s.value.SetValue(s.value.Upstream(ptr), item)
			acc = *(*int)(s.body(ptr))
		}
	}
	return unsafe.Pointer(&acc)
}
