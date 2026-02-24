package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"reflect"
	"strconv"
	"unsafe"
)

type int64CastExpr struct {
	operand *exec.Operand
	value   int64
}

func (e *int64CastExpr) compute(ptr unsafe.Pointer) unsafe.Pointer {
	vPtr := e.operand.Compute(ptr)
	switch e.operand.Kind() {
	case reflect.Int64:
		e.value = *(*int64)(vPtr)
	case reflect.Int:
		e.value = int64(*(*int)(vPtr))
	case reflect.Float64:
		e.value = int64(*(*float64)(vPtr))
	case reflect.Float32:
		e.value = int64(*(*float32)(vPtr))
	case reflect.String:
		i, err := strconv.Atoi(*(*string)(vPtr))
		if err != nil {
			panic(err)
		}
		e.value = int64(i)
	case reflect.Uint:
		e.value = int64(*(*uint)(vPtr))
	case reflect.Uint64:
		e.value = int64(*(*uint64)(vPtr))
	default:
		panic(fmt.Errorf("unsupported int64 cast: %s", e.operand.Type.Type().String()))
	}
	return unsafe.Pointer(&e.value)
}

// NewInt64Cast creates int64 cast evaluator.
func NewInt64Cast(x *Operand) (New, error) {
	if x == nil || x.Type == nil {
		return nil, fmt.Errorf("int64 cast operand type was empty")
	}
	return func(control *Control) (internal.Compute, error) {
		op, err := x.NewOperand(control)
		if err != nil {
			return nil, err
		}
		expr := &int64CastExpr{operand: op}
		return expr.compute, nil
	}, nil
}
