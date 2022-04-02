package et

import (
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/state"
	"unsafe"
)

type callExpr struct {
	operands []*state.Operand
	caller   state.Caller
}

func (e *callExpr) compute(ptr unsafe.Pointer) unsafe.Pointer {
	return e.caller.Call(ptr, e.operands)
}

//NewCaller crates a caller
func NewCaller(caller state.Caller, args []*Operand) (New, error) {
	expr := &callExpr{caller: caller}
	operands := Operands(args)
	return func(control *Control) (exec.Compute, error) {
		var err error
		if expr.operands, err = operands.operands(control); err != nil {
			return nil, err
		}
		return expr.compute, nil
	}, nil
}
