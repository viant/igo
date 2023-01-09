package et

import (
	"fmt"
	"github.com/viant/igo/internal"
	"reflect"
	"unsafe"
)

func NewStarExpr(operand *Operand) (New, error) {
	d := &derefProvider{operand: operand}
	switch operand.Selector.Type.Elem().Kind() {
	case reflect.Bool:
		return d.newBoolDeref, nil
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
		return d.newIntDeref, nil
	case reflect.Float64:
		return d.newFloat64Deref, nil
	case reflect.Float32:
		return d.newFloat32Deref, nil
	case reflect.String:
		return d.newStringDeref, nil

	}
	return nil, fmt.Errorf("not yet supported for %v", operand.Selector.Type.String())
}

type (
	derefProvider struct {
		operand *Operand
	}
	deref struct {
		compute func(ptr unsafe.Pointer) unsafe.Pointer
	}
)

func (p *derefProvider) newIntDeref(control *Control) (internal.Compute, error) {
	op, err := p.operand.NewOperand(control)
	if err != nil {
		return nil, err
	}
	deref := &deref{compute: op.Compute}
	return deref.int, nil
}

func (p *derefProvider) newFloat64Deref(control *Control) (internal.Compute, error) {
	op, err := p.operand.NewOperand(control)
	if err != nil {
		return nil, err
	}
	deref := &deref{compute: op.Compute}
	return deref.float64, nil
}

func (p *derefProvider) newFloat32Deref(control *Control) (internal.Compute, error) {
	op, err := p.operand.NewOperand(control)
	if err != nil {
		return nil, err
	}
	deref := &deref{compute: op.Compute}
	return deref.float32, nil
}

func (p *derefProvider) newStringDeref(control *Control) (internal.Compute, error) {
	op, err := p.operand.NewOperand(control)
	if err != nil {
		return nil, err
	}
	deref := &deref{compute: op.Compute}
	return deref.string, nil
}

func (p *derefProvider) newBoolDeref(control *Control) (internal.Compute, error) {
	op, err := p.operand.NewOperand(control)
	if err != nil {
		return nil, err
	}
	deref := &deref{compute: op.Compute}
	return deref.bool, nil
}

func (d *deref) int(ptr unsafe.Pointer) unsafe.Pointer {
	x := d.compute(ptr)
	value := *(**int)(x)
	return unsafe.Pointer(value)
}

func (d *deref) float64(ptr unsafe.Pointer) unsafe.Pointer {
	x := d.compute(ptr)
	value := *(**float64)(x)
	return unsafe.Pointer(value)
}

func (d *deref) float32(ptr unsafe.Pointer) unsafe.Pointer {
	x := d.compute(ptr)
	value := *(**float64)(x)
	return unsafe.Pointer(value)
}

func (d *deref) bool(ptr unsafe.Pointer) unsafe.Pointer {
	x := d.compute(ptr)
	value := *(**float64)(x)
	return unsafe.Pointer(value)
}

func (d *deref) string(ptr unsafe.Pointer) unsafe.Pointer {
	x := d.compute(ptr)
	value := *(**float64)(x)
	return unsafe.Pointer(value)
}
