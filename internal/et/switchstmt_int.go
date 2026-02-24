package et

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"unsafe"
)

// IntSwitchCase represents a compiled int switch case.
type IntSwitchCase struct {
	Values []int
	Body   New
}

type intSwitchCase struct {
	values []int
	body   internal.Compute
}

type intSwitchStmt struct {
	tag     *exec.Operand
	cases   []intSwitchCase
	defBody internal.Compute
	labelID uint32
}

func (s *intSwitchStmt) consumeBreak(flow *internal.Flow) bool {
	if !flow.HasBreak() {
		return false
	}
	if label := flow.Label(); label != 0 && label != s.labelID {
		return false
	}
	flow.ClearBreak()
	return true
}

func (s *intSwitchStmt) computeNoControl(ptr unsafe.Pointer) unsafe.Pointer {
	value := *(*int)(s.tag.Compute(ptr))
	for i := range s.cases {
		aCase := &s.cases[i]
		for _, v := range aCase.values {
			if v == value {
				return aCase.body(ptr)
			}
		}
	}
	if s.defBody != nil {
		return s.defBody(ptr)
	}
	return nil
}

func (s *intSwitchStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	value := *(*int)(s.tag.Compute(ptr))
	flow := internal.AsFlow(ptr)
	for i := range s.cases {
		aCase := &s.cases[i]
		matched := false
		for _, v := range aCase.values {
			if v == value {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}
		result := aCase.body(ptr)
		if s.consumeBreak(flow) {
			return nil
		}
		return result
	}
	if s.defBody != nil {
		result := s.defBody(ptr)
		if s.consumeBreak(flow) {
			return nil
		}
		return result
	}
	return nil
}

// NewIntSwitchStmt creates a typed int switch statement node.
func NewIntSwitchStmt(tag *Operand, cases []IntSwitchCase, defaultBody New, labelID uint32) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		tagOp, err := tag.NewOperand(control)
		if err != nil {
			return nil, err
		}
		stmt := &intSwitchStmt{tag: tagOp, cases: make([]intSwitchCase, len(cases)), labelID: labelID}
		for i := range cases {
			body, err := cases[i].Body(control)
			if err != nil {
				return nil, err
			}
			stmt.cases[i] = intSwitchCase{
				values: cases[i].Values,
				body:   body,
			}
		}
		if defaultBody != nil {
			body, err := defaultBody(control)
			if err != nil {
				return nil, err
			}
			stmt.defBody = body
		}
		if !control.HasControl() {
			return stmt.computeNoControl, nil
		}
		return stmt.compute, nil
	}, nil
}
