package et

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"unsafe"
)

type ifElseStmt struct {
	cond       *exec.Operand
	whenBranch internal.Compute
	elseBranch internal.Compute
}

func (s *ifElseStmt) computeIf(ptr unsafe.Pointer) unsafe.Pointer {
	if *(*bool)(s.cond.Compute(ptr)) {
		return s.whenBranch(ptr)
	}
	return nil
}

func (s *ifElseStmt) computeIfElse(ptr unsafe.Pointer) unsafe.Pointer {
	if *(*bool)(s.cond.Compute(ptr)) {
		return s.whenBranch(ptr)
	}
	return s.elseBranch(ptr)
}

//NewIfStmt creates if stmt
func NewIfStmt(cond *Operand, whenBranch, elseBranch New) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		var err error
		result := &ifElseStmt{}
		if result.cond, err = cond.NewOperand(control); err != nil {
			return nil, err
		}
		if result.whenBranch, err = whenBranch(control); err != nil {
			return nil, err
		}
		if elseBranch == nil {
			return result.computeIf, nil
		}
		if result.elseBranch, err = elseBranch(control); err != nil {
			return nil, err
		}
		return result.computeIfElse, nil
	}, nil
}
