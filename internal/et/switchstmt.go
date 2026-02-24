package et

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"unsafe"
)

// SwitchCase represents a compiled switch case.
type SwitchCase struct {
	Cond           *Operand
	Body           New
	HasFallthrough bool
}

type switchCase struct {
	cond           *exec.Operand
	body           internal.Compute
	hasFallthrough bool
}

type switchStmt struct {
	cases   []switchCase
	defBody internal.Compute
	labelID uint32
}

func (s *switchStmt) consumeBreak(flow *internal.Flow) bool {
	if !flow.HasBreak() {
		return false
	}
	if label := flow.Label(); label != 0 && label != s.labelID {
		return false
	}
	flow.ClearBreak()
	return true
}

func (s *switchStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	flow := internal.AsFlow(ptr)
	var result unsafe.Pointer
	matched := false
	for i := range s.cases {
		aCase := &s.cases[i]
		if !matched {
			if !*(*bool)(aCase.cond.Compute(ptr)) {
				continue
			}
			matched = true
		}

		result = aCase.body(ptr)
		if s.consumeBreak(flow) {
			return nil
		}
		if flow.HasFallthrough() || aCase.hasFallthrough {
			flow.ClearFallthrough()
			continue
		}
		return result
	}
	if matched || s.defBody == nil {
		return result
	}
	result = s.defBody(ptr)
	if s.consumeBreak(flow) {
		return nil
	}
	return result
}

func (s *switchStmt) computeNoControl(ptr unsafe.Pointer) unsafe.Pointer {
	var result unsafe.Pointer
	matched := false
	for i := range s.cases {
		aCase := &s.cases[i]
		if !matched {
			if !*(*bool)(aCase.cond.Compute(ptr)) {
				continue
			}
			matched = true
		}
		result = aCase.body(ptr)
		if aCase.hasFallthrough {
			continue
		}
		return result
	}
	if matched || s.defBody == nil {
		return result
	}
	return s.defBody(ptr)
}

// NewSwitchStmt creates a switch statement node.
func NewSwitchStmt(cases []SwitchCase, defaultBody New, labelID uint32) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &switchStmt{cases: make([]switchCase, len(cases)), labelID: labelID}
		hasFallthrough := false
		for i := range cases {
			cond, err := cases[i].Cond.NewOperand(control)
			if err != nil {
				return nil, err
			}
			body, err := cases[i].Body(control)
			if err != nil {
				return nil, err
			}
			stmt.cases[i] = switchCase{
				cond:           cond,
				body:           body,
				hasFallthrough: cases[i].HasFallthrough,
			}
			hasFallthrough = hasFallthrough || cases[i].HasFallthrough
		}
		if defaultBody != nil {
			body, err := defaultBody(control)
			if err != nil {
				return nil, err
			}
			stmt.defBody = body
		}
		if !control.HasControl() && !hasFallthrough {
			return stmt.computeNoControl, nil
		}
		return stmt.compute, nil
	}, nil
}
