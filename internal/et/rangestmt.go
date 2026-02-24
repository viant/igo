package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/xunsafe"
	"reflect"
	"time"
	"unsafe"
)

type rangeStmt struct {
	slice          *xunsafe.Slice
	body           internal.Compute
	isComponentPtr bool
	x              *exec.Operand
	key            *exec.Selector
	value          *exec.Selector
	labelID        uint32
	hasControl     bool
}

func (s *rangeStmt) step(ptr unsafe.Pointer) bool {
	s.body(ptr)
	if !s.hasControl {
		return true
	}
	flow := internal.AsFlow(ptr)
	if !flow.HasControl() {
		return true
	}
	if flow.HasBreak() {
		if label := flow.Label(); label == 0 || label == s.labelID {
			flow.ClearBreak()
		}
		return false
	}
	if flow.HasContinue() {
		if label := flow.Label(); label == 0 || label == s.labelID {
			flow.ClearContinue()
			return true
		}
		return false
	}
	return false
}

func (s *rangeStmt) computeRangeSlice(ptr unsafe.Pointer) unsafe.Pointer {
	slicePtr := s.x.Compute(ptr)
	s.slice.Range(slicePtr, func(index int, item interface{}) bool {
		if s.key != nil {
			*(*int)(s.key.Addr(ptr)) = index
		}
		if s.value != nil {
			itemPtr := xunsafe.AsPointer(item)
			*(*unsafe.Pointer)(s.value.Addr(ptr)) = itemPtr
		}
		if !s.step(ptr) {
			return false
		}
		return true
	})
	return nil
}

type mapRangeStmt struct {
	x          *exec.Operand
	mapType    reflect.Type
	mapCase    mapCase
	body       internal.Compute
	key        *exec.Selector
	value      *exec.Selector
	labelID    uint32
	hasControl bool
}

func (s *mapRangeStmt) step(ptr unsafe.Pointer) bool {
	s.body(ptr)
	if !s.hasControl {
		return true
	}
	flow := internal.AsFlow(ptr)
	if !flow.HasControl() {
		return true
	}
	if flow.HasBreak() {
		if label := flow.Label(); label == 0 || label == s.labelID {
			flow.ClearBreak()
		}
		return false
	}
	if flow.HasContinue() {
		if label := flow.Label(); label == 0 || label == s.labelID {
			flow.ClearContinue()
			return true
		}
		return false
	}
	return false
}

func (s *mapRangeStmt) computeMap(ptr unsafe.Pointer) unsafe.Pointer {
	switch s.mapCase {
	case mapCaseStringInt:
		return s.computeMapStringInt(ptr)
	case mapCaseStringString:
		return s.computeMapStringString(ptr)
	case mapCaseStringBool:
		return s.computeMapStringBool(ptr)
	case mapCaseStringFloat64:
		return s.computeMapStringFloat64(ptr)
	case mapCaseStringTime:
		return s.computeMapStringTime(ptr)
	case mapCaseIntInt:
		return s.computeMapIntInt(ptr)
	case mapCaseIntString:
		return s.computeMapIntString(ptr)
	case mapCaseIntBool:
		return s.computeMapIntBool(ptr)
	case mapCaseIntFloat64:
		return s.computeMapIntFloat64(ptr)
	case mapCaseInt64Int:
		return s.computeMapInt64Int(ptr)
	case mapCaseInt64String:
		return s.computeMapInt64String(ptr)
	case mapCaseInt64Bool:
		return s.computeMapInt64Bool(ptr)
	case mapCaseInt64Float64:
		return s.computeMapInt64Float64(ptr)
	case mapCaseUintFloat64:
		return s.computeMapUintFloat64(ptr)
	case mapCaseUint64Float64:
		return s.computeMapUint64Float64(ptr)
	default:
		return s.computeMapReflect(ptr)
	}
}

func (s *mapRangeStmt) computeMapStringInt(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]int)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*string)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*int)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapStringString(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]string)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*string)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*string)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapStringBool(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]bool)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*string)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*bool)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapStringFloat64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]float64)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*string)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*float64)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapStringTime(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]time.Time)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*string)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*time.Time)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapIntInt(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int]int)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*int)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*int)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapIntString(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int]string)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*int)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*string)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapIntBool(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int]bool)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*int)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*bool)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapIntFloat64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int]float64)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*int)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*float64)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapInt64Int(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int64]int)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*int64)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*int)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapInt64String(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int64]string)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*int64)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*string)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapInt64Bool(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int64]bool)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*int64)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*bool)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapInt64Float64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int64]float64)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*int64)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*float64)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapUintFloat64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[uint]float64)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*uint)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*float64)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) computeMapUint64Float64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[uint64]float64)(s.x.Compute(ptr))
	for k, v := range m {
		if s.key != nil {
			*(*uint64)(s.key.Addr(ptr)) = k
		}
		if s.value != nil {
			*(*float64)(s.value.Addr(ptr)) = v
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

func (s *mapRangeStmt) setReflect(sel *exec.Selector, ptr unsafe.Pointer, value reflect.Value) {
	dest := reflect.NewAt(sel.Type, sel.Addr(ptr)).Elem()
	if value.Type() == sel.Type {
		dest.Set(value)
		return
	}
	if value.Type().ConvertibleTo(sel.Type) {
		dest.Set(value.Convert(sel.Type))
	}
}

func (s *mapRangeStmt) computeMapReflect(ptr unsafe.Pointer) unsafe.Pointer {
	mapValue := reflect.NewAt(s.mapType, s.x.Compute(ptr)).Elem()
	iter := mapValue.MapRange()
	for iter.Next() {
		if s.key != nil {
			s.setReflect(s.key, ptr, iter.Key())
		}
		if s.value != nil {
			s.setReflect(s.value, ptr, iter.Value())
		}
		if !s.step(ptr) {
			return nil
		}
	}
	return nil
}

// NewRange creates a range stmt.
func NewRange(x *Operand, key, value *exec.Selector, body New, labelID uint32) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		var err error
		xType := x.Type.Type()
		if xType.Kind() == reflect.Ptr {
			xType = xType.Elem()
		}

		xOp, err := x.NewOperand(control)
		if err != nil {
			return nil, err
		}
		bodyOp, err := body(control)
		if err != nil {
			return nil, err
		}

		switch xType.Kind() {
		case reflect.Slice:
			stmt := &rangeStmt{
				slice:      xunsafe.NewSlice(xType),
				key:        key,
				value:      value,
				labelID:    labelID,
				hasControl: control.HasControl(),
			}
			stmt.isComponentPtr = xType.Elem().Kind() == reflect.Ptr
			stmt.x = xOp
			stmt.body = bodyOp
			return stmt.computeRangeSlice, nil
		case reflect.Map:
			stmt := &mapRangeStmt{
				x:          xOp,
				mapType:    xType,
				mapCase:    detectMapCase(xType),
				body:       bodyOp,
				key:        key,
				value:      value,
				labelID:    labelID,
				hasControl: control.HasControl(),
			}
			return stmt.computeMap, nil
		}
		return nil, fmt.Errorf("range not supported for type %s", xType.String())
	}, nil
}
