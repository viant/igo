package et

import (
	"github.com/viant/igo/internal"
	"go/token"
	"unsafe"
)

type intForStmt struct {
	init    internal.Compute
	body    internal.Compute
	iOffset uintptr
	cmpOp   token.Token
	limit   int
	delta   int
	labelID uint32
}

const (
	forActionNone = iota
	forActionBreak
	forActionContinue
	forActionReturn
)

func (s *intForStmt) flowAction(flow *internal.Flow) int {
	if flow.HasBreak() {
		if label := flow.Label(); label != 0 && label != s.labelID {
			return forActionReturn
		}
		flow.ClearBreak()
		return forActionBreak
	}
	if flow.HasContinue() {
		if label := flow.Label(); label != 0 && label != s.labelID {
			return forActionReturn
		}
		flow.ClearContinue()
		return forActionContinue
	}
	if flow.HasGoto() || flow.HasReturn() {
		return forActionReturn
	}
	return forActionNone
}

func (s *intForStmt) computeNoFlow(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	iPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.iOffset))
	switch s.cmpOp {
	case token.LSS:
		for ; *iPtr < s.limit; *iPtr += s.delta {
			s.body(ptr)
		}
	case token.LEQ:
		for ; *iPtr <= s.limit; *iPtr += s.delta {
			s.body(ptr)
		}
	case token.GTR:
		for ; *iPtr > s.limit; *iPtr += s.delta {
			s.body(ptr)
		}
	case token.GEQ:
		for ; *iPtr >= s.limit; *iPtr += s.delta {
			s.body(ptr)
		}
	case token.EQL:
		for ; *iPtr == s.limit; *iPtr += s.delta {
			s.body(ptr)
		}
	case token.NEQ:
		for ; *iPtr != s.limit; *iPtr += s.delta {
			s.body(ptr)
		}
	}
	return nil
}

func (s *intForStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	iPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.iOffset))
	flow := internal.AsFlow(ptr)
	step := func() (unsafe.Pointer, int) {
		r := s.body(ptr)
		return r, s.flowAction(flow)
	}

	switch s.cmpOp {
	case token.LSS:
		for ; *iPtr < s.limit; *iPtr += s.delta {
			r, action := step()
			if action == forActionBreak {
				return nil
			}
			if action == forActionContinue {
				continue
			}
			if action == forActionReturn {
				return r
			}
		}
	case token.LEQ:
		for ; *iPtr <= s.limit; *iPtr += s.delta {
			r, action := step()
			if action == forActionBreak {
				return nil
			}
			if action == forActionContinue {
				continue
			}
			if action == forActionReturn {
				return r
			}
		}
	case token.GTR:
		for ; *iPtr > s.limit; *iPtr += s.delta {
			r, action := step()
			if action == forActionBreak {
				return nil
			}
			if action == forActionContinue {
				continue
			}
			if action == forActionReturn {
				return r
			}
		}
	case token.GEQ:
		for ; *iPtr >= s.limit; *iPtr += s.delta {
			r, action := step()
			if action == forActionBreak {
				return nil
			}
			if action == forActionContinue {
				continue
			}
			if action == forActionReturn {
				return r
			}
		}
	case token.EQL:
		for ; *iPtr == s.limit; *iPtr += s.delta {
			r, action := step()
			if action == forActionBreak {
				return nil
			}
			if action == forActionContinue {
				continue
			}
			if action == forActionReturn {
				return r
			}
		}
	case token.NEQ:
		for ; *iPtr != s.limit; *iPtr += s.delta {
			r, action := step()
			if action == forActionBreak {
				return nil
			}
			if action == forActionContinue {
				continue
			}
			if action == forActionReturn {
				return r
			}
		}
	}
	return nil
}

// NewForIntConstStmt creates a specialized for-loop with int counter and const limit.
func NewForIntConstStmt(init, body New, iOffset uintptr, cmpOp token.Token, limit, delta int, labelID uint32) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &intForStmt{
			iOffset: iOffset,
			cmpOp:   cmpOp,
			limit:   limit,
			delta:   delta,
			labelID: labelID,
		}
		var err error
		if init != nil {
			if stmt.init, err = init(control); err != nil {
				return nil, err
			}
		}
		if body != nil {
			if stmt.body, err = body(control); err != nil {
				return nil, err
			}
		}
		if !control.HasControl() {
			return stmt.computeNoFlow, nil
		}
		return stmt.compute, nil
	}, nil
}
