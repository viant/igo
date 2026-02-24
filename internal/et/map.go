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

var mapStringIntType = reflect.TypeOf(map[string]int{})
var mapStringStringType = reflect.TypeOf(map[string]string{})
var mapStringBoolType = reflect.TypeOf(map[string]bool{})
var mapStringFloat64Type = reflect.TypeOf(map[string]float64{})
var mapStringTimeType = reflect.TypeOf(map[string]time.Time{})

var mapIntIntType = reflect.TypeOf(map[int]int{})
var mapIntStringType = reflect.TypeOf(map[int]string{})
var mapIntBoolType = reflect.TypeOf(map[int]bool{})
var mapIntFloat64Type = reflect.TypeOf(map[int]float64{})

var mapInt64IntType = reflect.TypeOf(map[int64]int{})
var mapInt64StringType = reflect.TypeOf(map[int64]string{})
var mapInt64BoolType = reflect.TypeOf(map[int64]bool{})
var mapInt64Float64Type = reflect.TypeOf(map[int64]float64{})
var mapUintFloat64Type = reflect.TypeOf(map[uint]float64{})
var mapUint64Float64Type = reflect.TypeOf(map[uint64]float64{})

type mapCase int

const (
	mapCaseUnknown mapCase = iota
	mapCaseStringInt
	mapCaseStringString
	mapCaseStringBool
	mapCaseStringFloat64
	mapCaseStringTime
	mapCaseIntInt
	mapCaseIntString
	mapCaseIntBool
	mapCaseIntFloat64
	mapCaseInt64Int
	mapCaseInt64String
	mapCaseInt64Bool
	mapCaseInt64Float64
	mapCaseUintFloat64
	mapCaseUint64Float64
)

type mapGetter interface {
	compute(ptr unsafe.Pointer) unsafe.Pointer
	existsPtr() *bool
}

type mapSetter interface {
	compute(ptr unsafe.Pointer) unsafe.Pointer
}

// NewMapGet creates a map index expression evaluator.
func NewMapGet(m, key *Operand) (New, reflect.Type, error) {
	if m == nil || m.Type == nil {
		return nil, nil, fmt.Errorf("map operand type was empty")
	}
	mapType := m.Type.Type()
	if mapType.Kind() != reflect.Map {
		return nil, nil, fmt.Errorf("operand %v is not a map", mapType.String())
	}
	if key == nil || key.Type == nil {
		return nil, nil, fmt.Errorf("map key operand type was empty")
	}
	return func(control *Control) (internal.Compute, error) {
		mapOp, err := m.NewOperand(control)
		if err != nil {
			return nil, err
		}
		keyOp, err := key.NewOperand(control)
		if err != nil {
			return nil, err
		}
		getter := newMapGetter(mapType, mapOp, keyOp)
		return getter.compute, nil
	}, mapType.Elem(), nil
}

// NewMapSet creates a map index assignment evaluator.
func NewMapSet(m, key, value *Operand) (New, error) {
	if m == nil || m.Type == nil {
		return nil, fmt.Errorf("map operand type was empty")
	}
	mapType := m.Type.Type()
	if mapType.Kind() != reflect.Map {
		return nil, fmt.Errorf("operand %v is not a map", mapType.String())
	}
	if key == nil || key.Type == nil {
		return nil, fmt.Errorf("map key operand type was empty")
	}
	if value == nil || value.Type == nil {
		return nil, fmt.Errorf("map value operand type was empty")
	}
	if value.Type.Type() != mapType.Elem() {
		return nil, fmt.Errorf("map value type mismatch, expected %v but had %v", mapType.Elem(), value.Type.Type())
	}
	return func(control *Control) (internal.Compute, error) {
		mapOp, err := m.NewOperand(control)
		if err != nil {
			return nil, err
		}
		keyOp, err := key.NewOperand(control)
		if err != nil {
			return nil, err
		}
		valueOp, err := value.NewOperand(control)
		if err != nil {
			return nil, err
		}
		setter := newMapSetter(mapType, mapOp, keyOp, valueOp)
		return setter.compute, nil
	}, nil
}

// NewMapGet2Assign creates a map index value/ok assignment evaluator.
func NewMapGet2Assign(m, key, valueDest, okDest *Operand) (New, error) {
	if okDest == nil || okDest.Type == nil || okDest.Type.Type().Kind() != reflect.Bool {
		return nil, fmt.Errorf("value, ok assignment expects bool ok destination")
	}
	if m == nil || m.Type == nil {
		return nil, fmt.Errorf("map operand type was empty")
	}
	mapType := m.Type.Type()
	if mapType.Kind() != reflect.Map {
		return nil, fmt.Errorf("operand %v is not a map", mapType.String())
	}
	if key == nil || key.Type == nil {
		return nil, fmt.Errorf("map key operand type was empty")
	}
	if valueDest == nil || valueDest.Type == nil {
		return nil, fmt.Errorf("map value destination type was empty")
	}
	if valueDest.Type.Type() != mapType.Elem() {
		return nil, fmt.Errorf("map value destination type mismatch, expected %v but had %v", mapType.Elem(), valueDest.Type.Type())
	}
	return func(control *Control) (internal.Compute, error) {
		mapOp, err := m.NewOperand(control)
		if err != nil {
			return nil, err
		}
		keyOp, err := key.NewOperand(control)
		if err != nil {
			return nil, err
		}
		valDestOp, err := valueDest.NewOperand(control)
		if err != nil {
			return nil, err
		}
		okDestOp, err := okDest.NewOperand(control)
		if err != nil {
			return nil, err
		}
		assigner := &mapGet2Assign{
			mapGet:    newMapGetter(mapType, mapOp, keyOp),
			valDestOp: valDestOp,
			okDestOp:  okDestOp,
			valueSize: int(mapType.Elem().Size()),
		}
		assigner.exists = assigner.mapGet.existsPtr()
		return assigner.compute, nil
	}, nil
}

func detectMapCase(mapType reflect.Type) mapCase {
	switch mapType {
	case mapStringIntType:
		return mapCaseStringInt
	case mapStringStringType:
		return mapCaseStringString
	case mapStringBoolType:
		return mapCaseStringBool
	case mapStringFloat64Type:
		return mapCaseStringFloat64
	case mapStringTimeType:
		return mapCaseStringTime
	case mapIntIntType:
		return mapCaseIntInt
	case mapIntStringType:
		return mapCaseIntString
	case mapIntBoolType:
		return mapCaseIntBool
	case mapIntFloat64Type:
		return mapCaseIntFloat64
	case mapInt64IntType:
		return mapCaseInt64Int
	case mapInt64StringType:
		return mapCaseInt64String
	case mapInt64BoolType:
		return mapCaseInt64Bool
	case mapInt64Float64Type:
		return mapCaseInt64Float64
	case mapUintFloat64Type:
		return mapCaseUintFloat64
	case mapUint64Float64Type:
		return mapCaseUint64Float64
	default:
		return mapCaseUnknown
	}
}

func newMapGetter(mapType reflect.Type, mapOp, keyOp *exec.Operand) mapGetter {
	switch detectMapCase(mapType) {
	case mapCaseStringInt:
		return &mapGetStringInt{mapOp: mapOp, keyOp: keyOp}
	case mapCaseStringString:
		return &mapGetStringString{mapOp: mapOp, keyOp: keyOp}
	case mapCaseStringBool:
		return &mapGetStringBool{mapOp: mapOp, keyOp: keyOp}
	case mapCaseStringFloat64:
		return &mapGetStringFloat64{mapOp: mapOp, keyOp: keyOp}
	case mapCaseStringTime:
		return &mapGetStringTime{mapOp: mapOp, keyOp: keyOp}
	case mapCaseIntInt:
		return &mapGetIntInt{mapOp: mapOp, keyOp: keyOp}
	case mapCaseIntString:
		return &mapGetIntString{mapOp: mapOp, keyOp: keyOp}
	case mapCaseIntBool:
		return &mapGetIntBool{mapOp: mapOp, keyOp: keyOp}
	case mapCaseIntFloat64:
		return &mapGetIntFloat64{mapOp: mapOp, keyOp: keyOp}
	case mapCaseInt64Int:
		return &mapGetInt64Int{mapOp: mapOp, keyOp: keyOp}
	case mapCaseInt64String:
		return &mapGetInt64String{mapOp: mapOp, keyOp: keyOp}
	case mapCaseInt64Bool:
		return &mapGetInt64Bool{mapOp: mapOp, keyOp: keyOp}
	case mapCaseInt64Float64:
		return &mapGetInt64Float64{mapOp: mapOp, keyOp: keyOp}
	case mapCaseUintFloat64:
		return &mapGetUintFloat64{mapOp: mapOp, keyOp: keyOp}
	case mapCaseUint64Float64:
		return &mapGetUint64Float64{mapOp: mapOp, keyOp: keyOp}
	}
	getter := newReflectMapGet(mapType, mapOp, keyOp)
	return getter
}

func newMapSetter(mapType reflect.Type, mapOp, keyOp, valueOp *exec.Operand) mapSetter {
	switch detectMapCase(mapType) {
	case mapCaseStringInt:
		return &mapSetStringInt{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseStringString:
		return &mapSetStringString{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseStringBool:
		return &mapSetStringBool{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseStringFloat64:
		return &mapSetStringFloat64{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseStringTime:
		return &mapSetStringTime{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseIntInt:
		return &mapSetIntInt{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseIntString:
		return &mapSetIntString{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseIntBool:
		return &mapSetIntBool{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseIntFloat64:
		return &mapSetIntFloat64{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseInt64Int:
		return &mapSetInt64Int{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseInt64String:
		return &mapSetInt64String{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseInt64Bool:
		return &mapSetInt64Bool{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseInt64Float64:
		return &mapSetInt64Float64{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseUintFloat64:
		return &mapSetUintFloat64{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	case mapCaseUint64Float64:
		return &mapSetUint64Float64{mapOp: mapOp, keyOp: keyOp, valueOp: valueOp}
	}
	setter := &reflectMapSet{
		mapOp:   mapOp,
		keyOp:   keyOp,
		valueOp: valueOp,
		mapType: mapType,
		keyType: mapType.Key(),
		valType: mapType.Elem(),
		keyKind: mapType.Key().Kind(),
	}
	setter.initKey()
	return setter
}

type mapGetStringInt struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  int
	exists bool
}

func (m *mapGetStringInt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]int)(m.mapOp.Compute(ptr))
	key := *(*string)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetStringInt) existsPtr() *bool { return &m.exists }

type mapGetStringString struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  string
	exists bool
}

func (m *mapGetStringString) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]string)(m.mapOp.Compute(ptr))
	key := *(*string)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetStringString) existsPtr() *bool { return &m.exists }

type mapGetStringBool struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  bool
	exists bool
}

func (m *mapGetStringBool) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]bool)(m.mapOp.Compute(ptr))
	key := *(*string)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetStringBool) existsPtr() *bool { return &m.exists }

type mapGetStringFloat64 struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  float64
	exists bool
}

func (m *mapGetStringFloat64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]float64)(m.mapOp.Compute(ptr))
	key := *(*string)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetStringFloat64) existsPtr() *bool { return &m.exists }

type mapGetStringTime struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  time.Time
	exists bool
}

func (m *mapGetStringTime) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]time.Time)(m.mapOp.Compute(ptr))
	key := *(*string)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetStringTime) existsPtr() *bool { return &m.exists }

type mapGetIntInt struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  int
	exists bool
}

func (m *mapGetIntInt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int]int)(m.mapOp.Compute(ptr))
	key := *(*int)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetIntInt) existsPtr() *bool { return &m.exists }

type mapGetIntString struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  string
	exists bool
}

func (m *mapGetIntString) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int]string)(m.mapOp.Compute(ptr))
	key := *(*int)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetIntString) existsPtr() *bool { return &m.exists }

type mapGetIntBool struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  bool
	exists bool
}

func (m *mapGetIntBool) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int]bool)(m.mapOp.Compute(ptr))
	key := *(*int)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetIntBool) existsPtr() *bool { return &m.exists }

type mapGetIntFloat64 struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  float64
	exists bool
}

func (m *mapGetIntFloat64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int]float64)(m.mapOp.Compute(ptr))
	key := *(*int)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetIntFloat64) existsPtr() *bool { return &m.exists }

type mapGetInt64Int struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  int
	exists bool
}

func (m *mapGetInt64Int) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int64]int)(m.mapOp.Compute(ptr))
	key := *(*int64)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetInt64Int) existsPtr() *bool { return &m.exists }

type mapGetInt64String struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  string
	exists bool
}

func (m *mapGetInt64String) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int64]string)(m.mapOp.Compute(ptr))
	key := *(*int64)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetInt64String) existsPtr() *bool { return &m.exists }

type mapGetInt64Bool struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  bool
	exists bool
}

func (m *mapGetInt64Bool) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int64]bool)(m.mapOp.Compute(ptr))
	key := *(*int64)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetInt64Bool) existsPtr() *bool { return &m.exists }

type mapGetInt64Float64 struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  float64
	exists bool
}

func (m *mapGetInt64Float64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int64]float64)(m.mapOp.Compute(ptr))
	key := *(*int64)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetInt64Float64) existsPtr() *bool { return &m.exists }

type mapGetUintFloat64 struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  float64
	exists bool
}

func (m *mapGetUintFloat64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[uint]float64)(m.mapOp.Compute(ptr))
	key := *(*uint)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetUintFloat64) existsPtr() *bool { return &m.exists }

type mapGetUint64Float64 struct {
	mapOp  *exec.Operand
	keyOp  *exec.Operand
	value  float64
	exists bool
}

func (m *mapGetUint64Float64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[uint64]float64)(m.mapOp.Compute(ptr))
	key := *(*uint64)(m.keyOp.Compute(ptr))
	m.value, m.exists = mp[key]
	return unsafe.Pointer(&m.value)
}

func (m *mapGetUint64Float64) existsPtr() *bool { return &m.exists }

type mapSetStringInt struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetStringInt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]int)(m.mapOp.Compute(ptr))
	mp[*(*string)(m.keyOp.Compute(ptr))] = *(*int)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetStringString struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetStringString) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]string)(m.mapOp.Compute(ptr))
	mp[*(*string)(m.keyOp.Compute(ptr))] = *(*string)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetStringBool struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetStringBool) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]bool)(m.mapOp.Compute(ptr))
	mp[*(*string)(m.keyOp.Compute(ptr))] = *(*bool)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetStringFloat64 struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetStringFloat64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]float64)(m.mapOp.Compute(ptr))
	mp[*(*string)(m.keyOp.Compute(ptr))] = *(*float64)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetStringTime struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetStringTime) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[string]time.Time)(m.mapOp.Compute(ptr))
	mp[*(*string)(m.keyOp.Compute(ptr))] = *(*time.Time)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetIntInt struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetIntInt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int]int)(m.mapOp.Compute(ptr))
	mp[*(*int)(m.keyOp.Compute(ptr))] = *(*int)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetIntString struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetIntString) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int]string)(m.mapOp.Compute(ptr))
	mp[*(*int)(m.keyOp.Compute(ptr))] = *(*string)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetIntBool struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetIntBool) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int]bool)(m.mapOp.Compute(ptr))
	mp[*(*int)(m.keyOp.Compute(ptr))] = *(*bool)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetIntFloat64 struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetIntFloat64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int]float64)(m.mapOp.Compute(ptr))
	mp[*(*int)(m.keyOp.Compute(ptr))] = *(*float64)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetInt64Int struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetInt64Int) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int64]int)(m.mapOp.Compute(ptr))
	mp[*(*int64)(m.keyOp.Compute(ptr))] = *(*int)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetInt64String struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetInt64String) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int64]string)(m.mapOp.Compute(ptr))
	mp[*(*int64)(m.keyOp.Compute(ptr))] = *(*string)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetInt64Bool struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetInt64Bool) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int64]bool)(m.mapOp.Compute(ptr))
	mp[*(*int64)(m.keyOp.Compute(ptr))] = *(*bool)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetInt64Float64 struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetInt64Float64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[int64]float64)(m.mapOp.Compute(ptr))
	mp[*(*int64)(m.keyOp.Compute(ptr))] = *(*float64)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetUintFloat64 struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetUintFloat64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[uint]float64)(m.mapOp.Compute(ptr))
	mp[*(*uint)(m.keyOp.Compute(ptr))] = *(*float64)(m.valueOp.Compute(ptr))
	return nil
}

type mapSetUint64Float64 struct {
	mapOp   *exec.Operand
	keyOp   *exec.Operand
	valueOp *exec.Operand
}

func (m *mapSetUint64Float64) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mp := *(*map[uint64]float64)(m.mapOp.Compute(ptr))
	mp[*(*uint64)(m.keyOp.Compute(ptr))] = *(*float64)(m.valueOp.Compute(ptr))
	return nil
}

type reflectMapGet struct {
	mapOp      *exec.Operand
	keyOp      *exec.Operand
	mapType    reflect.Type
	keyType    reflect.Type
	keyKind    reflect.Kind
	keyScratch reflect.Value
	value      reflect.Value
	valuePtr   unsafe.Pointer
	exists     bool
}

func newReflectMapGet(mapType reflect.Type, mapOp, keyOp *exec.Operand) *reflectMapGet {
	holder := reflect.New(mapType.Elem())
	result := &reflectMapGet{
		mapOp:    mapOp,
		keyOp:    keyOp,
		mapType:  mapType,
		keyType:  mapType.Key(),
		keyKind:  mapType.Key().Kind(),
		value:    holder,
		valuePtr: unsafe.Pointer(holder.Pointer()),
	}
	result.initKey()
	return result
}

func (m *reflectMapGet) initKey() {
	switch m.keyKind {
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.String:
		m.keyScratch = reflect.New(m.keyType).Elem()
	}
}

func (m *reflectMapGet) key(ptr unsafe.Pointer) reflect.Value {
	keyPtr := m.keyOp.Compute(ptr)
	switch m.keyKind {
	case reflect.Int:
		m.keyScratch.SetInt(int64(*(*int)(keyPtr)))
		return m.keyScratch
	case reflect.Int64:
		m.keyScratch.SetInt(*(*int64)(keyPtr))
		return m.keyScratch
	case reflect.Uint:
		m.keyScratch.SetUint(uint64(*(*uint)(keyPtr)))
		return m.keyScratch
	case reflect.Uint64:
		m.keyScratch.SetUint(*(*uint64)(keyPtr))
		return m.keyScratch
	case reflect.String:
		m.keyScratch.SetString(*(*string)(keyPtr))
		return m.keyScratch
	default:
		keyValue := reflect.ValueOf(m.keyOp.Type.Interface(keyPtr))
		if keyValue.Type() == m.keyType {
			return keyValue
		}
		if keyValue.Type().ConvertibleTo(m.keyType) {
			return keyValue.Convert(m.keyType)
		}
		return keyValue
	}
}

func (m *reflectMapGet) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mapPtr := m.mapOp.Compute(ptr)
	mapValue := reflect.NewAt(m.mapType, mapPtr).Elem()
	item := mapValue.MapIndex(m.key(ptr))
	m.exists = item.IsValid()
	if m.exists {
		m.value.Elem().Set(item)
	} else {
		m.value.Elem().SetZero()
	}
	return m.valuePtr
}

func (m *reflectMapGet) existsPtr() *bool { return &m.exists }

type reflectMapSet struct {
	mapOp      *exec.Operand
	keyOp      *exec.Operand
	valueOp    *exec.Operand
	mapType    reflect.Type
	keyType    reflect.Type
	valType    reflect.Type
	keyKind    reflect.Kind
	keyScratch reflect.Value
}

func (m *reflectMapSet) initKey() {
	switch m.keyKind {
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.String:
		m.keyScratch = reflect.New(m.keyType).Elem()
	}
}

func (m *reflectMapSet) key(ptr unsafe.Pointer) reflect.Value {
	keyPtr := m.keyOp.Compute(ptr)
	switch m.keyKind {
	case reflect.Int:
		m.keyScratch.SetInt(int64(*(*int)(keyPtr)))
		return m.keyScratch
	case reflect.Int64:
		m.keyScratch.SetInt(*(*int64)(keyPtr))
		return m.keyScratch
	case reflect.Uint:
		m.keyScratch.SetUint(uint64(*(*uint)(keyPtr)))
		return m.keyScratch
	case reflect.Uint64:
		m.keyScratch.SetUint(*(*uint64)(keyPtr))
		return m.keyScratch
	case reflect.String:
		m.keyScratch.SetString(*(*string)(keyPtr))
		return m.keyScratch
	default:
		keyValue := reflect.ValueOf(m.keyOp.Type.Interface(keyPtr))
		if keyValue.Type() == m.keyType {
			return keyValue
		}
		if keyValue.Type().ConvertibleTo(m.keyType) {
			return keyValue.Convert(m.keyType)
		}
		return keyValue
	}
}

func (m *reflectMapSet) compute(ptr unsafe.Pointer) unsafe.Pointer {
	mapPtr := m.mapOp.Compute(ptr)
	mapValue := reflect.NewAt(m.mapType, mapPtr).Elem()
	valuePtr := m.valueOp.Compute(ptr)
	value := reflect.NewAt(m.valType, valuePtr).Elem()
	mapValue.SetMapIndex(m.key(ptr), value)
	return nil
}

type mapGet2Assign struct {
	mapGet    mapGetter
	exists    *bool
	valDestOp *exec.Operand
	okDestOp  *exec.Operand
	valueSize int
}

func (m *mapGet2Assign) compute(ptr unsafe.Pointer) unsafe.Pointer {
	valuePtr := m.mapGet.compute(ptr)
	destPtr := m.valDestOp.Compute(ptr)
	xunsafe.Copy(destPtr, valuePtr, m.valueSize)
	okPtr := m.okDestOp.Compute(ptr)
	if m.exists != nil {
		*(*bool)(okPtr) = *m.exists
	} else {
		*(*bool)(okPtr) = false
	}
	return nil
}
