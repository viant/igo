package et

import (
	"github.com/viant/igo/internal"
	"go/token"
	"unsafe"
)

const (
	forIntCallGuardKindNone = iota
	forIntCallGuardKindCmp
	forIntCallGuardKindModCmp
)

const (
	ForIntCallGuardKindNone   = forIntCallGuardKindNone
	ForIntCallGuardKindCmp    = forIntCallGuardKindCmp
	ForIntCallGuardKindModCmp = forIntCallGuardKindModCmp
)

type intForCallGuardStmt struct {
	init         func(ptr unsafe.Pointer) unsafe.Pointer
	iOffset      uintptr
	cmpOp        token.Token
	limit        int
	delta        int
	call         func(int)
	guardKind    int
	guardOp      token.Token
	guardConst   int
	guardModBase int
}

func (s *intForCallGuardStmt) checkGuard(i int) bool {
	lhs := i
	switch s.guardKind {
	case forIntCallGuardKindCmp:
		// use i as-is
	case forIntCallGuardKindModCmp:
		lhs = i % s.guardModBase
	default:
		return true
	}
	switch s.guardOp {
	case token.LSS:
		return lhs < s.guardConst
	case token.LEQ:
		return lhs <= s.guardConst
	case token.GTR:
		return lhs > s.guardConst
	case token.GEQ:
		return lhs >= s.guardConst
	case token.EQL:
		return lhs == s.guardConst
	case token.NEQ:
		return lhs != s.guardConst
	}
	return false
}

func (s *intForCallGuardStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	iPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.iOffset))
	switch s.cmpOp {
	case token.LSS:
		for ; *iPtr < s.limit; *iPtr += s.delta {
			if s.checkGuard(*iPtr) {
				s.call(*iPtr)
			}
		}
	case token.LEQ:
		for ; *iPtr <= s.limit; *iPtr += s.delta {
			if s.checkGuard(*iPtr) {
				s.call(*iPtr)
			}
		}
	case token.GTR:
		for ; *iPtr > s.limit; *iPtr += s.delta {
			if s.checkGuard(*iPtr) {
				s.call(*iPtr)
			}
		}
	case token.GEQ:
		for ; *iPtr >= s.limit; *iPtr += s.delta {
			if s.checkGuard(*iPtr) {
				s.call(*iPtr)
			}
		}
	case token.EQL:
		for ; *iPtr == s.limit; *iPtr += s.delta {
			if s.checkGuard(*iPtr) {
				s.call(*iPtr)
			}
		}
	case token.NEQ:
		for ; *iPtr != s.limit; *iPtr += s.delta {
			if s.checkGuard(*iPtr) {
				s.call(*iPtr)
			}
		}
	}
	return nil
}

// NewForIntCallGuardConstStmt creates fused int for-loop for guarded call body:
// `if guard(i) { fn(i) }`.
func NewForIntCallGuardConstStmt(init New, iOffset uintptr, cmpOp token.Token, limit, delta int, call func(int), guardKind int, guardOp token.Token, guardConst, guardModBase int) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &intForCallGuardStmt{
			iOffset:      iOffset,
			cmpOp:        cmpOp,
			limit:        limit,
			delta:        delta,
			call:         call,
			guardKind:    guardKind,
			guardOp:      guardOp,
			guardConst:   guardConst,
			guardModBase: guardModBase,
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
