package et

import (
	"github.com/viant/igo/internal"
	"go/token"
	"unsafe"
)

type intForCallStmt struct {
	init    func(ptr unsafe.Pointer) unsafe.Pointer
	iOffset uintptr
	cmpOp   token.Token
	limit   int
	delta   int
	call    func(int)
}

func (s *intForCallStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	iPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.iOffset))
	switch s.cmpOp {
	case token.LSS:
		for ; *iPtr < s.limit; *iPtr += s.delta {
			s.call(*iPtr)
		}
	case token.LEQ:
		for ; *iPtr <= s.limit; *iPtr += s.delta {
			s.call(*iPtr)
		}
	case token.GTR:
		for ; *iPtr > s.limit; *iPtr += s.delta {
			s.call(*iPtr)
		}
	case token.GEQ:
		for ; *iPtr >= s.limit; *iPtr += s.delta {
			s.call(*iPtr)
		}
	case token.EQL:
		for ; *iPtr == s.limit; *iPtr += s.delta {
			s.call(*iPtr)
		}
	case token.NEQ:
		for ; *iPtr != s.limit; *iPtr += s.delta {
			s.call(*iPtr)
		}
	}
	return nil
}

// NewForIntCallConstStmt creates fused int for-loop for `fn(i)` body.
func NewForIntCallConstStmt(init New, iOffset uintptr, cmpOp token.Token, limit, delta int, call func(int)) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &intForCallStmt{
			iOffset: iOffset,
			cmpOp:   cmpOp,
			limit:   limit,
			delta:   delta,
			call:    call,
		}
		if init != nil {
			initCompute, err := init(control)
			if err != nil {
				return nil, err
			}
			stmt.init = initCompute
		}
		return stmt.compute, nil
	}, nil
}
