package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type typeAssertExpr struct {
	x        *exec.Operand
	target   reflect.Type
	value    reflect.Value
	valuePtr unsafe.Pointer
}

func newTypeAssertExpr(x *exec.Operand, target reflect.Type) *typeAssertExpr {
	holder := reflect.New(target)
	return &typeAssertExpr{
		x:        x,
		target:   target,
		value:    holder,
		valuePtr: unsafe.Pointer(holder.Pointer()),
	}
}

func (e *typeAssertExpr) assert(ptr unsafe.Pointer) (unsafe.Pointer, bool) {
	v := reflect.NewAt(e.x.Type.Type(), e.x.Compute(ptr)).Elem().Interface()
	if v == nil {
		e.value.Elem().SetZero()
		return e.valuePtr, false
	}
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	ok := false
	if e.target.Kind() == reflect.Interface {
		ok = rt.Implements(e.target) || rt.AssignableTo(e.target)
	} else {
		ok = rt.AssignableTo(e.target)
	}
	if !ok {
		e.value.Elem().SetZero()
		return e.valuePtr, false
	}
	e.value.Elem().Set(rv)
	return e.valuePtr, true
}

func (e *typeAssertExpr) compute(ptr unsafe.Pointer) unsafe.Pointer {
	valuePtr, ok := e.assert(ptr)
	if !ok {
		v := reflect.NewAt(e.x.Type.Type(), e.x.Compute(ptr)).Elem().Interface()
		actual := "<nil>"
		if v != nil {
			actual = reflect.TypeOf(v).String()
		}
		panic(fmt.Errorf("interface conversion: interface {} is %s, not %s", actual, e.target.String()))
	}
	return valuePtr
}

type typeAssert2Assign struct {
	expr      *typeAssertExpr
	dest      *exec.Operand
	okDest    *exec.Operand
	valueSize int
}

func (e *typeAssert2Assign) compute(ptr unsafe.Pointer) unsafe.Pointer {
	valuePtr, ok := e.expr.assert(ptr)
	destPtr := e.dest.Compute(ptr)
	xunsafe.Copy(destPtr, valuePtr, e.valueSize)
	okPtr := e.okDest.Compute(ptr)
	*(*bool)(okPtr) = ok
	return nil
}

// NewTypeAssert creates type assertion expression evaluator.
func NewTypeAssert(x *Operand, target reflect.Type) (New, reflect.Type, error) {
	if x == nil || x.Type == nil {
		return nil, nil, fmt.Errorf("type assertion operand type was empty")
	}
	return func(control *Control) (internal.Compute, error) {
		xOp, err := x.NewOperand(control)
		if err != nil {
			return nil, err
		}
		expr := newTypeAssertExpr(xOp, target)
		return expr.compute, nil
	}, target, nil
}

// NewTypeAssert2Assign creates value,ok type assertion assignment evaluator.
func NewTypeAssert2Assign(x *Operand, target reflect.Type, valueDest, okDest *Operand) (New, error) {
	if okDest == nil || okDest.Type == nil || okDest.Type.Kind() != reflect.Bool {
		return nil, fmt.Errorf("type assertion value,ok expects bool destination")
	}
	return func(control *Control) (internal.Compute, error) {
		xOp, err := x.NewOperand(control)
		if err != nil {
			return nil, err
		}
		valueDestOp, err := valueDest.NewOperand(control)
		if err != nil {
			return nil, err
		}
		okDestOp, err := okDest.NewOperand(control)
		if err != nil {
			return nil, err
		}
		expr := newTypeAssertExpr(xOp, target)
		assign := &typeAssert2Assign{
			expr:      expr,
			dest:      valueDestOp,
			okDest:    okDestOp,
			valueSize: int(target.Size()),
		}
		return assign.compute, nil
	}, nil
}
