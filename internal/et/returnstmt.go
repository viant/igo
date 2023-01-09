package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"unsafe"
)

type returnStmt struct {
	operands []*exec.Operand
	results  []*exec.Selector
}

func (s *returnStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	*(*internal.Flow)(ptr) |= internal.RtReturn
	return nil
}

func (s *returnStmt) computeR1(ptr unsafe.Pointer) unsafe.Pointer {
	retValuePtr := s.operands[0].Compute(ptr)
	s.results[0].SetValue(s.results[0].Upstream(ptr), s.operands[0].Interface(retValuePtr))
	*(*internal.Flow)(ptr) |= internal.RtReturn
	return retValuePtr
}

func (s *returnStmt) computeR2(ptr unsafe.Pointer) unsafe.Pointer {
	retValuePtr0 := s.operands[0].Compute(ptr)
	retValuePtr1 := s.operands[1].Compute(ptr)

	s.results[0].SetValue(s.results[0].Upstream(ptr), s.operands[0].Interface(retValuePtr0))
	s.results[1].SetValue(s.results[1].Upstream(ptr), s.operands[1].Interface(retValuePtr0))
	*(*internal.Flow)(ptr) |= internal.RtReturn
	return unsafe.Pointer(&[2]unsafe.Pointer{
		retValuePtr0,
		retValuePtr1,
	})
}

func (s *returnStmt) computeR3(ptr unsafe.Pointer) unsafe.Pointer {
	retValuePtr0 := s.operands[0].Compute(ptr)
	retValuePtr1 := s.operands[1].Compute(ptr)
	retValuePtr2 := s.operands[2].Compute(ptr)

	s.results[0].SetValue(s.results[0].Upstream(ptr), s.operands[0].Interface(retValuePtr0))
	s.results[1].SetValue(s.results[1].Upstream(ptr), s.operands[1].Interface(retValuePtr0))
	s.results[2].SetValue(s.results[2].Upstream(ptr), s.operands[2].Interface(retValuePtr0))
	*(*internal.Flow)(ptr) |= internal.RtReturn
	return unsafe.Pointer(&[3]unsafe.Pointer{
		retValuePtr0,
		retValuePtr1,
		retValuePtr2,
	})
}

func (s *returnStmt) computeR4(ptr unsafe.Pointer) unsafe.Pointer {
	retValuePtr0 := s.operands[0].Compute(ptr)
	retValuePtr1 := s.operands[1].Compute(ptr)
	retValuePtr2 := s.operands[2].Compute(ptr)
	retValuePtr3 := s.operands[3].Compute(ptr)

	s.results[0].SetValue(s.results[0].Upstream(ptr), s.operands[0].Interface(retValuePtr0))
	s.results[1].SetValue(s.results[1].Upstream(ptr), s.operands[1].Interface(retValuePtr0))
	s.results[2].SetValue(s.results[2].Upstream(ptr), s.operands[2].Interface(retValuePtr0))
	s.results[3].SetValue(s.results[3].Upstream(ptr), s.operands[3].Interface(retValuePtr0))

	*(*internal.Flow)(ptr) |= internal.RtReturn
	return unsafe.Pointer(&[4]unsafe.Pointer{
		retValuePtr0,
		retValuePtr1,
		retValuePtr2,
		retValuePtr3,
	})
}

//NewReturnStmt creates a return stmt
func NewReturnStmt(retOperands Operands, results []*exec.Selector) (New, error) {
	return func(control *Control) (internal.Compute, error) {

		operands, err := retOperands.operands(control)
		if err != nil {
			return nil, err
		}
		if len(operands) != len(results) {
			return nil, fmt.Errorf("invalid return args len, expected:%v, had: %v", len(results), len(operands))
		}
		returnStmt := &returnStmt{
			operands: operands,
			results:  results,
		}

		switch len(results) {
		case 0:
			return returnStmt.compute, nil
		case 1:
			return returnStmt.computeR1, nil
		case 2:
			return returnStmt.computeR2, nil
		case 3:
			return returnStmt.computeR3, nil
		case 4:
			return returnStmt.computeR4, nil
		default:
			return nil, fmt.Errorf("too many return values")
		}
	}, nil
}
