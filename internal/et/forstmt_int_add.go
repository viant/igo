package et

import (
	"github.com/viant/igo/internal"
	"go/token"
	"unsafe"
)

type intForAddStmt struct {
	init       func(ptr unsafe.Pointer) unsafe.Pointer
	destOffset uintptr
	iOffset    uintptr
	cmpOp      token.Token
	limit      int
	delta      int
}

func (s *intForAddStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	destPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.destOffset))
	iPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.iOffset))
	switch s.cmpOp {
	case token.LSS:
		for ; *iPtr < s.limit; *iPtr += s.delta {
			*destPtr += *iPtr
		}
	case token.LEQ:
		for ; *iPtr <= s.limit; *iPtr += s.delta {
			*destPtr += *iPtr
		}
	case token.GTR:
		for ; *iPtr > s.limit; *iPtr += s.delta {
			*destPtr += *iPtr
		}
	case token.GEQ:
		for ; *iPtr >= s.limit; *iPtr += s.delta {
			*destPtr += *iPtr
		}
	case token.EQL:
		for ; *iPtr == s.limit; *iPtr += s.delta {
			*destPtr += *iPtr
		}
	case token.NEQ:
		for ; *iPtr != s.limit; *iPtr += s.delta {
			*destPtr += *iPtr
		}
	}
	return nil
}

// NewForIntAddConstStmt creates fused int for-loop for `dest += i` body.
func NewForIntAddConstStmt(init New, destOffset, iOffset uintptr, cmpOp token.Token, limit, delta int) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &intForAddStmt{
			destOffset: destOffset,
			iOffset:    iOffset,
			cmpOp:      cmpOp,
			limit:      limit,
			delta:      delta,
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
