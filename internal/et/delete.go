package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"reflect"
	"time"
	"unsafe"
)

type mapDelete struct {
	mapOp      *exec.Operand
	keyOp      *exec.Operand
	mapType    reflect.Type
	keyType    reflect.Type
	keyKind    reflect.Kind
	keyScratch reflect.Value
}

func (d *mapDelete) initKey() {
	switch d.keyKind {
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.String:
		d.keyScratch = reflect.New(d.keyType).Elem()
	}
}

func (d *mapDelete) key(ptr unsafe.Pointer) reflect.Value {
	keyPtr := d.keyOp.Compute(ptr)
	switch d.keyKind {
	case reflect.Int:
		d.keyScratch.SetInt(int64(*(*int)(keyPtr)))
		return d.keyScratch
	case reflect.Int64:
		d.keyScratch.SetInt(*(*int64)(keyPtr))
		return d.keyScratch
	case reflect.Uint:
		d.keyScratch.SetUint(uint64(*(*uint)(keyPtr)))
		return d.keyScratch
	case reflect.Uint64:
		d.keyScratch.SetUint(*(*uint64)(keyPtr))
		return d.keyScratch
	case reflect.String:
		d.keyScratch.SetString(*(*string)(keyPtr))
		return d.keyScratch
	default:
		keyValue := reflect.ValueOf(d.keyOp.Type.Interface(keyPtr))
		if keyValue.Type() == d.keyType {
			return keyValue
		}
		if keyValue.Type().ConvertibleTo(d.keyType) {
			return keyValue.Convert(d.keyType)
		}
		return keyValue
	}
}

func (d *mapDelete) computeReflect(ptr unsafe.Pointer) unsafe.Pointer {
	mapPtr := d.mapOp.Compute(ptr)
	mapValue := reflect.NewAt(d.mapType, mapPtr).Elem()
	mapValue.SetMapIndex(d.key(ptr), reflect.Value{})
	return nil
}

func (d *mapDelete) computeFastStringInt(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]int)(d.mapOp.Compute(ptr))
	delete(m, *(*string)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastStringString(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]string)(d.mapOp.Compute(ptr))
	delete(m, *(*string)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastStringBool(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]bool)(d.mapOp.Compute(ptr))
	delete(m, *(*string)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastStringFloat64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]float64)(d.mapOp.Compute(ptr))
	delete(m, *(*string)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastStringTime(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[string]time.Time)(d.mapOp.Compute(ptr))
	delete(m, *(*string)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastIntInt(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int]int)(d.mapOp.Compute(ptr))
	delete(m, *(*int)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastIntString(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int]string)(d.mapOp.Compute(ptr))
	delete(m, *(*int)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastIntBool(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int]bool)(d.mapOp.Compute(ptr))
	delete(m, *(*int)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastIntFloat64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int]float64)(d.mapOp.Compute(ptr))
	delete(m, *(*int)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastInt64Int(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int64]int)(d.mapOp.Compute(ptr))
	delete(m, *(*int64)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastInt64String(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int64]string)(d.mapOp.Compute(ptr))
	delete(m, *(*int64)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastInt64Bool(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int64]bool)(d.mapOp.Compute(ptr))
	delete(m, *(*int64)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastInt64Float64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[int64]float64)(d.mapOp.Compute(ptr))
	delete(m, *(*int64)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastUintFloat64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[uint]float64)(d.mapOp.Compute(ptr))
	delete(m, *(*uint)(d.keyOp.Compute(ptr)))
	return nil
}

func (d *mapDelete) computeFastUint64Float64(ptr unsafe.Pointer) unsafe.Pointer {
	m := *(*map[uint64]float64)(d.mapOp.Compute(ptr))
	delete(m, *(*uint64)(d.keyOp.Compute(ptr)))
	return nil
}

// NewDelete creates map delete evaluator.
func NewDelete(m, key *Operand) (New, error) {
	if m == nil || m.Type == nil {
		return nil, fmt.Errorf("delete map operand type was empty")
	}
	mapType := m.Type.Type()
	if mapType.Kind() != reflect.Map {
		return nil, fmt.Errorf("delete expects map, but got %s", mapType.String())
	}
	if key == nil || key.Type == nil {
		return nil, fmt.Errorf("delete key operand type was empty")
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
		stmt := &mapDelete{
			mapOp:   mapOp,
			keyOp:   keyOp,
			mapType: mapType,
			keyType: mapType.Key(),
			keyKind: mapType.Key().Kind(),
		}
		switch detectMapCase(mapType) {
		case mapCaseStringInt:
			return stmt.computeFastStringInt, nil
		case mapCaseStringString:
			return stmt.computeFastStringString, nil
		case mapCaseStringBool:
			return stmt.computeFastStringBool, nil
		case mapCaseStringFloat64:
			return stmt.computeFastStringFloat64, nil
		case mapCaseStringTime:
			return stmt.computeFastStringTime, nil
		case mapCaseIntInt:
			return stmt.computeFastIntInt, nil
		case mapCaseIntString:
			return stmt.computeFastIntString, nil
		case mapCaseIntBool:
			return stmt.computeFastIntBool, nil
		case mapCaseIntFloat64:
			return stmt.computeFastIntFloat64, nil
		case mapCaseInt64Int:
			return stmt.computeFastInt64Int, nil
		case mapCaseInt64String:
			return stmt.computeFastInt64String, nil
		case mapCaseInt64Bool:
			return stmt.computeFastInt64Bool, nil
		case mapCaseInt64Float64:
			return stmt.computeFastInt64Float64, nil
		case mapCaseUintFloat64:
			return stmt.computeFastUintFloat64, nil
		case mapCaseUint64Float64:
			return stmt.computeFastUint64Float64, nil
		default:
			stmt.initKey()
			return stmt.computeReflect, nil
		}
	}, nil
}
