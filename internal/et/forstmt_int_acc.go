package et

import (
	"github.com/viant/igo/internal"
	"go/token"
	"unsafe"
)

type intForAccStmt struct {
	init        func(ptr unsafe.Pointer) unsafe.Pointer
	destOffset  uintptr
	iOffset     uintptr
	cmpOp       token.Token
	limit       int
	delta       int
	assignOp    token.Token
	useLoopVar  bool
	rhsConst    int
	guardEnable bool
	guardOp     token.Token
	guardConst  int
}

func (s *intForAccStmt) guard(i int) bool {
	if !s.guardEnable {
		return true
	}
	switch s.guardOp {
	case token.LSS:
		return i < s.guardConst
	case token.LEQ:
		return i <= s.guardConst
	case token.GTR:
		return i > s.guardConst
	case token.GEQ:
		return i >= s.guardConst
	case token.EQL:
		return i == s.guardConst
	case token.NEQ:
		return i != s.guardConst
	}
	return false
}

func (s *intForAccStmt) apply(dest *int, i int) {
	rhs := s.rhsConst
	if s.useLoopVar {
		rhs = i
	}
	switch s.assignOp {
	case token.ADD_ASSIGN:
		*dest += rhs
	case token.SUB_ASSIGN:
		*dest -= rhs
	case token.MUL_ASSIGN:
		*dest *= rhs
	case token.QUO_ASSIGN:
		*dest /= rhs
	}
}

func (s *intForAccStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	destPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.destOffset))
	iPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.iOffset))
	switch s.cmpOp {
	case token.LSS:
		for ; *iPtr < s.limit; *iPtr += s.delta {
			if s.guard(*iPtr) {
				s.apply(destPtr, *iPtr)
			}
		}
	case token.LEQ:
		for ; *iPtr <= s.limit; *iPtr += s.delta {
			if s.guard(*iPtr) {
				s.apply(destPtr, *iPtr)
			}
		}
	case token.GTR:
		for ; *iPtr > s.limit; *iPtr += s.delta {
			if s.guard(*iPtr) {
				s.apply(destPtr, *iPtr)
			}
		}
	case token.GEQ:
		for ; *iPtr >= s.limit; *iPtr += s.delta {
			if s.guard(*iPtr) {
				s.apply(destPtr, *iPtr)
			}
		}
	case token.EQL:
		for ; *iPtr == s.limit; *iPtr += s.delta {
			if s.guard(*iPtr) {
				s.apply(destPtr, *iPtr)
			}
		}
	case token.NEQ:
		for ; *iPtr != s.limit; *iPtr += s.delta {
			if s.guard(*iPtr) {
				s.apply(destPtr, *iPtr)
			}
		}
	}
	return nil
}

// NewForIntAccConstStmt creates fused int for-loop for assignment bodies like
// `sum += i`, `sum -= i`, `sum += 3`, `sum *= i`, optionally guarded by `if i < K`.
func NewForIntAccConstStmt(init New, destOffset, iOffset uintptr, cmpOp token.Token, limit, delta int, assignOp token.Token, useLoopVar bool, rhsConst int, guardEnable bool, guardOp token.Token, guardConst int) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &intForAccStmt{
			destOffset:  destOffset,
			iOffset:     iOffset,
			cmpOp:       cmpOp,
			limit:       limit,
			delta:       delta,
			assignOp:    assignOp,
			useLoopVar:  useLoopVar,
			rhsConst:    rhsConst,
			guardEnable: guardEnable,
			guardOp:     guardOp,
			guardConst:  guardConst,
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
