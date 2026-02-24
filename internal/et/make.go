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

type makeSlice struct {
	len      *exec.Operand
	cap      *exec.Operand
	destType reflect.Type
	mapCase  mapCase
	mSI      map[string]int
	mSS      map[string]string
	mSB      map[string]bool
	mSF      map[string]float64
	mST      map[string]time.Time
	mII      map[int]int
	mIS      map[int]string
	mIB      map[int]bool
	mIF      map[int]float64
	mI64I    map[int64]int
	mI64S    map[int64]string
	mI64B    map[int64]bool
	mI64F    map[int64]float64
	mUF      map[uint]float64
	mU64F    map[uint64]float64
}

// Compute computes make stmt
func (m *makeSlice) Compute(ptr unsafe.Pointer) unsafe.Pointer {
	size := 0
	if m.len != nil {
		size = *(*int)(m.len.Compute(ptr))
	}
	var value reflect.Value
	switch m.destType.Kind() {
	case reflect.Slice:
		sliceCap := size
		if m.cap != nil {
			sliceCap = *(*int)(m.cap.Compute(ptr))
		}
		value = reflect.MakeSlice(m.destType, size, sliceCap)
		return xunsafe.ValuePointer(&value)
	case reflect.Chan:
		value = reflect.MakeChan(m.destType, size)
	case reflect.Map:
		if m.mapCase != mapCaseUnknown {
			return m.computeFastMap(size)
		}
		if m.len != nil {
			value = reflect.MakeMapWithSize(m.destType, size)
		} else {
			value = reflect.MakeMap(m.destType)
		}
	}
	holder := reflect.New(m.destType)
	holder.Elem().Set(value)
	return unsafe.Pointer(holder.Pointer())
}

func (m *makeSlice) computeFastMap(size int) unsafe.Pointer {
	switch m.mapCase {
	case mapCaseStringInt:
		if m.mSI == nil {
			m.mSI = make(map[string]int, size)
		} else {
			for k := range m.mSI {
				delete(m.mSI, k)
			}
		}
		return unsafe.Pointer(&m.mSI)
	case mapCaseStringString:
		if m.mSS == nil {
			m.mSS = make(map[string]string, size)
		} else {
			for k := range m.mSS {
				delete(m.mSS, k)
			}
		}
		return unsafe.Pointer(&m.mSS)
	case mapCaseStringBool:
		if m.mSB == nil {
			m.mSB = make(map[string]bool, size)
		} else {
			for k := range m.mSB {
				delete(m.mSB, k)
			}
		}
		return unsafe.Pointer(&m.mSB)
	case mapCaseStringFloat64:
		if m.mSF == nil {
			m.mSF = make(map[string]float64, size)
		} else {
			for k := range m.mSF {
				delete(m.mSF, k)
			}
		}
		return unsafe.Pointer(&m.mSF)
	case mapCaseStringTime:
		if m.mST == nil {
			m.mST = make(map[string]time.Time, size)
		} else {
			for k := range m.mST {
				delete(m.mST, k)
			}
		}
		return unsafe.Pointer(&m.mST)
	case mapCaseIntInt:
		if m.mII == nil {
			m.mII = make(map[int]int, size)
		} else {
			for k := range m.mII {
				delete(m.mII, k)
			}
		}
		return unsafe.Pointer(&m.mII)
	case mapCaseIntString:
		if m.mIS == nil {
			m.mIS = make(map[int]string, size)
		} else {
			for k := range m.mIS {
				delete(m.mIS, k)
			}
		}
		return unsafe.Pointer(&m.mIS)
	case mapCaseIntBool:
		if m.mIB == nil {
			m.mIB = make(map[int]bool, size)
		} else {
			for k := range m.mIB {
				delete(m.mIB, k)
			}
		}
		return unsafe.Pointer(&m.mIB)
	case mapCaseIntFloat64:
		if m.mIF == nil {
			m.mIF = make(map[int]float64, size)
		} else {
			for k := range m.mIF {
				delete(m.mIF, k)
			}
		}
		return unsafe.Pointer(&m.mIF)
	case mapCaseInt64Int:
		if m.mI64I == nil {
			m.mI64I = make(map[int64]int, size)
		} else {
			for k := range m.mI64I {
				delete(m.mI64I, k)
			}
		}
		return unsafe.Pointer(&m.mI64I)
	case mapCaseInt64String:
		if m.mI64S == nil {
			m.mI64S = make(map[int64]string, size)
		} else {
			for k := range m.mI64S {
				delete(m.mI64S, k)
			}
		}
		return unsafe.Pointer(&m.mI64S)
	case mapCaseInt64Bool:
		if m.mI64B == nil {
			m.mI64B = make(map[int64]bool, size)
		} else {
			for k := range m.mI64B {
				delete(m.mI64B, k)
			}
		}
		return unsafe.Pointer(&m.mI64B)
	case mapCaseInt64Float64:
		if m.mI64F == nil {
			m.mI64F = make(map[int64]float64, size)
		} else {
			for k := range m.mI64F {
				delete(m.mI64F, k)
			}
		}
		return unsafe.Pointer(&m.mI64F)
	case mapCaseUintFloat64:
		if m.mUF == nil {
			m.mUF = make(map[uint]float64, size)
		} else {
			for k := range m.mUF {
				delete(m.mUF, k)
			}
		}
		return unsafe.Pointer(&m.mUF)
	case mapCaseUint64Float64:
		if m.mU64F == nil {
			m.mU64F = make(map[uint64]float64, size)
		} else {
			for k := range m.mU64F {
				delete(m.mU64F, k)
			}
		}
		return unsafe.Pointer(&m.mU64F)
	}
	return nil
}

// NewMake creates a make stmt
func NewMake(destType reflect.Type, args Operands) (New, reflect.Type, error) {
	switch destType.Kind() {
	case reflect.Slice, reflect.Chan, reflect.Map:
	default:
		return nil, nil, fmt.Errorf("invalid make type: %v", destType.String())
	}
	switch destType.Kind() {
	case reflect.Slice:
		if len(args) < 1 || len(args) > 2 {
			return nil, nil, fmt.Errorf("invalid make args for slice: %d", len(args))
		}
	case reflect.Chan:
		if len(args) > 1 {
			return nil, nil, fmt.Errorf("invalid make args for chan: %d", len(args))
		}
	case reflect.Map:
		if len(args) > 1 {
			return nil, nil, fmt.Errorf("invalid make args for map: %d", len(args))
		}
	}
	return func(control *Control) (internal.Compute, error) {
		operands, err := args.operands(control)
		if err != nil {
			return nil, err
		}
		result := &makeSlice{destType: destType}
		if destType.Kind() == reflect.Map {
			result.mapCase = detectMapCase(destType)
		}
		if len(operands) > 0 {
			result.len = operands[0]
		}
		if len(operands) > 1 {
			result.cap = operands[1]
		}
		return result.Compute, nil
	}, destType, nil
}
