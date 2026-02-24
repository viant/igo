package et

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"unsafe"
)

type callExpr struct {
	operands []*exec.Operand
	caller   exec.Caller
}

func (e *callExpr) compute(ptr unsafe.Pointer) unsafe.Pointer {
	return e.caller.Call(ptr, e.operands)
}

type callStmt struct {
	operands   []*exec.Operand
	caller     exec.Caller
	voidCaller exec.VoidCaller
}

func (s *callStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	s.caller.Call(ptr, s.operands)
	return nil
}

func (s *callStmt) computeVoid(ptr unsafe.Pointer) unsafe.Pointer {
	s.voidCaller.CallVoid(ptr, s.operands)
	return nil
}

// NewCaller crates a caller
func NewCaller(caller exec.Caller, args []*Operand) (New, error) {
	expr := &callExpr{caller: caller}
	operands := Operands(args)
	return func(control *Control) (internal.Compute, error) {
		var err error
		if expr.operands, err = operands.operands(control); err != nil {
			return nil, err
		}
		return expr.compute, nil
	}, nil
}

// NewCallerStmt creates a caller node for expression statements (result ignored).
func NewCallerStmt(caller exec.Caller, args []*Operand) (New, error) {
	stmt := &callStmt{caller: caller}
	stmt.voidCaller, _ = caller.(exec.VoidCaller)
	operands := Operands(args)
	return func(control *Control) (internal.Compute, error) {
		var err error
		if stmt.operands, err = operands.operands(control); err != nil {
			return nil, err
		}
		if stmt.voidCaller != nil {
			return stmt.computeVoid, nil
		}
		return stmt.compute, nil
	}, nil
}
