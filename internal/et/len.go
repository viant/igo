package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type lenExpr struct {
	operand *exec.Operand
	kind    reflect.Kind
	slice   *xunsafe.Slice
	mType   reflect.Type
	value   int
}

func (e *lenExpr) compute(ptr unsafe.Pointer) unsafe.Pointer {
	p := e.operand.Compute(ptr)
	switch e.kind {
	case reflect.String:
		e.value = len(*(*string)(p))
	case reflect.Slice:
		e.value = e.slice.Len(p)
	case reflect.Map:
		e.value = reflect.NewAt(e.mType, p).Elem().Len()
	default:
		e.value = 0
	}
	return unsafe.Pointer(&e.value)
}

// NewLen creates a built-in len expression.
func NewLen(x *Operand) (New, reflect.Type, error) {
	if x == nil || x.Type == nil {
		return nil, nil, fmt.Errorf("len operand type was empty")
	}
	xType := x.Type.Type()
	return func(control *Control) (internal.Compute, error) {
		operand, err := x.NewOperand(control)
		if err != nil {
			return nil, err
		}
		expr := &lenExpr{
			operand: operand,
			kind:    xType.Kind(),
			mType:   xType,
		}
		if xType.Kind() == reflect.Slice {
			expr.slice = xunsafe.NewSlice(xType)
		}
		return expr.compute, nil
	}, reflect.TypeOf(0), nil
}
