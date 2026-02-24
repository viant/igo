package et

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/xunsafe"
	"reflect"
	"time"
	"unsafe"
)

// MapEntry represents a single map literal entry.
type MapEntry struct {
	Key   *Operand
	Value *Operand
}

// NewMapComposite returns map composite literal evaluator.
func NewMapComposite(cType reflect.Type, entries []MapEntry) New {
	mapType := cType
	isPtr := false
	if mapType.Kind() == reflect.Ptr {
		isPtr = true
		mapType = mapType.Elem()
	}
	return func(control *Control) (internal.Compute, error) {
		keys := make([]*exec.Operand, len(entries))
		values := make([]*exec.Operand, len(entries))
		for i := range entries {
			var err error
			if keys[i], err = entries[i].Key.NewOperand(control); err != nil {
				return nil, err
			}
			if values[i], err = entries[i].Value.NewOperand(control); err != nil {
				return nil, err
			}
		}
		composite := &compositeMap{
			keys:      keys,
			values:    values,
			mapType:   mapType,
			keyType:   mapType.Key(),
			valueType: mapType.Elem(),
			mapCase:   detectMapCase(mapType),
			isPtr:     isPtr,
		}
		if composite.mapCase != mapCaseUnknown {
			return composite.computeFast, nil
		}
		return composite.computeReflect, nil
	}
}

type compositeMap struct {
	keys      []*exec.Operand
	values    []*exec.Operand
	mapType   reflect.Type
	keyType   reflect.Type
	valueType reflect.Type
	mapCase   mapCase
	isPtr     bool
	mSI       map[string]int
	mSS       map[string]string
	mSB       map[string]bool
	mSF       map[string]float64
	mST       map[string]time.Time
	mII       map[int]int
	mIS       map[int]string
	mIB       map[int]bool
	mIF       map[int]float64
	mI64I     map[int64]int
	mI64S     map[int64]string
	mI64B     map[int64]bool
	mI64F     map[int64]float64
	mUF       map[uint]float64
	mU64F     map[uint64]float64
}

func (c *compositeMap) resultPtr(mapPtr unsafe.Pointer) unsafe.Pointer {
	if c.isPtr {
		return xunsafe.RefPointer(mapPtr)
	}
	return mapPtr
}

func (c *compositeMap) computeFast(ptr unsafe.Pointer) unsafe.Pointer {
	switch c.mapCase {
	case mapCaseStringInt:
		if c.mSI == nil {
			c.mSI = make(map[string]int, len(c.keys))
		} else {
			for k := range c.mSI {
				delete(c.mSI, k)
			}
		}
		for i := range c.keys {
			c.mSI[*(*string)(c.keys[i].Compute(ptr))] = *(*int)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mSI))
	case mapCaseStringString:
		if c.mSS == nil {
			c.mSS = make(map[string]string, len(c.keys))
		} else {
			for k := range c.mSS {
				delete(c.mSS, k)
			}
		}
		for i := range c.keys {
			c.mSS[*(*string)(c.keys[i].Compute(ptr))] = *(*string)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mSS))
	case mapCaseStringBool:
		if c.mSB == nil {
			c.mSB = make(map[string]bool, len(c.keys))
		} else {
			for k := range c.mSB {
				delete(c.mSB, k)
			}
		}
		for i := range c.keys {
			c.mSB[*(*string)(c.keys[i].Compute(ptr))] = *(*bool)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mSB))
	case mapCaseStringFloat64:
		if c.mSF == nil {
			c.mSF = make(map[string]float64, len(c.keys))
		} else {
			for k := range c.mSF {
				delete(c.mSF, k)
			}
		}
		for i := range c.keys {
			c.mSF[*(*string)(c.keys[i].Compute(ptr))] = *(*float64)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mSF))
	case mapCaseStringTime:
		if c.mST == nil {
			c.mST = make(map[string]time.Time, len(c.keys))
		} else {
			for k := range c.mST {
				delete(c.mST, k)
			}
		}
		for i := range c.keys {
			c.mST[*(*string)(c.keys[i].Compute(ptr))] = *(*time.Time)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mST))
	case mapCaseIntInt:
		if c.mII == nil {
			c.mII = make(map[int]int, len(c.keys))
		} else {
			for k := range c.mII {
				delete(c.mII, k)
			}
		}
		for i := range c.keys {
			c.mII[*(*int)(c.keys[i].Compute(ptr))] = *(*int)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mII))
	case mapCaseIntString:
		if c.mIS == nil {
			c.mIS = make(map[int]string, len(c.keys))
		} else {
			for k := range c.mIS {
				delete(c.mIS, k)
			}
		}
		for i := range c.keys {
			c.mIS[*(*int)(c.keys[i].Compute(ptr))] = *(*string)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mIS))
	case mapCaseIntBool:
		if c.mIB == nil {
			c.mIB = make(map[int]bool, len(c.keys))
		} else {
			for k := range c.mIB {
				delete(c.mIB, k)
			}
		}
		for i := range c.keys {
			c.mIB[*(*int)(c.keys[i].Compute(ptr))] = *(*bool)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mIB))
	case mapCaseIntFloat64:
		if c.mIF == nil {
			c.mIF = make(map[int]float64, len(c.keys))
		} else {
			for k := range c.mIF {
				delete(c.mIF, k)
			}
		}
		for i := range c.keys {
			c.mIF[*(*int)(c.keys[i].Compute(ptr))] = *(*float64)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mIF))
	case mapCaseInt64Int:
		if c.mI64I == nil {
			c.mI64I = make(map[int64]int, len(c.keys))
		} else {
			for k := range c.mI64I {
				delete(c.mI64I, k)
			}
		}
		for i := range c.keys {
			c.mI64I[*(*int64)(c.keys[i].Compute(ptr))] = *(*int)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mI64I))
	case mapCaseInt64String:
		if c.mI64S == nil {
			c.mI64S = make(map[int64]string, len(c.keys))
		} else {
			for k := range c.mI64S {
				delete(c.mI64S, k)
			}
		}
		for i := range c.keys {
			c.mI64S[*(*int64)(c.keys[i].Compute(ptr))] = *(*string)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mI64S))
	case mapCaseInt64Bool:
		if c.mI64B == nil {
			c.mI64B = make(map[int64]bool, len(c.keys))
		} else {
			for k := range c.mI64B {
				delete(c.mI64B, k)
			}
		}
		for i := range c.keys {
			c.mI64B[*(*int64)(c.keys[i].Compute(ptr))] = *(*bool)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mI64B))
	case mapCaseInt64Float64:
		if c.mI64F == nil {
			c.mI64F = make(map[int64]float64, len(c.keys))
		} else {
			for k := range c.mI64F {
				delete(c.mI64F, k)
			}
		}
		for i := range c.keys {
			c.mI64F[*(*int64)(c.keys[i].Compute(ptr))] = *(*float64)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mI64F))
	case mapCaseUintFloat64:
		if c.mUF == nil {
			c.mUF = make(map[uint]float64, len(c.keys))
		} else {
			for k := range c.mUF {
				delete(c.mUF, k)
			}
		}
		for i := range c.keys {
			c.mUF[*(*uint)(c.keys[i].Compute(ptr))] = *(*float64)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mUF))
	case mapCaseUint64Float64:
		if c.mU64F == nil {
			c.mU64F = make(map[uint64]float64, len(c.keys))
		} else {
			for k := range c.mU64F {
				delete(c.mU64F, k)
			}
		}
		for i := range c.keys {
			c.mU64F[*(*uint64)(c.keys[i].Compute(ptr))] = *(*float64)(c.values[i].Compute(ptr))
		}
		return c.resultPtr(unsafe.Pointer(&c.mU64F))
	}
	return c.computeReflect(ptr)
}

func (c *compositeMap) computeReflect(ptr unsafe.Pointer) unsafe.Pointer {
	mapValue := reflect.MakeMapWithSize(c.mapType, len(c.keys))
	for i := range c.keys {
		keyPtr := c.keys[i].Compute(ptr)
		valuePtr := c.values[i].Compute(ptr)
		key := reflect.NewAt(c.keyType, keyPtr).Elem()
		value := reflect.NewAt(c.valueType, valuePtr).Elem()
		mapValue.SetMapIndex(key, value)
	}
	holder := reflect.New(c.mapType)
	holder.Elem().Set(mapValue)
	mapPtr := unsafe.Pointer(holder.Pointer())
	return c.resultPtr(mapPtr)
}
