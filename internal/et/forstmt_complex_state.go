package et

import (
	"github.com/viant/igo/internal"
	"unsafe"
)

type complexStateForStmt struct {
	init         internal.Compute
	iOffset      uintptr
	limit        int
	sumOffset    uintptr
	countOffset  uintptr
	modeOffset   uintptr
	outputOffset uintptr
	flagsOffset  uintptr
	scoresOffset uintptr
	breakLimit   int
	breakOnGE    bool
}

func (s *complexStateForStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	iPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.iOffset))
	sumPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.sumOffset))
	countPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.countOffset))
	modePtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.modeOffset))
	outputPtr := (*int)(unsafe.Pointer(uintptr(ptr) + s.outputOffset))
	flagsPtr := (*map[string]bool)(unsafe.Pointer(uintptr(ptr) + s.flagsOffset))
	scoresPtr := (*map[int]int)(unsafe.Pointer(uintptr(ptr) + s.scoresOffset))

	i := *iPtr
	mod5 := i % 5
	mod7 := i % 7
	mod4 := i % 4
	for i < s.limit {
		if (i & 1) == 0 {
			*sumPtr += i
		} else {
			*sumPtr -= i / 2
		}

		switch mod5 {
		case 0:
			*modePtr = 1
		case 1:
			*modePtr = 2
		case 2:
			*modePtr = 3
		case 3:
			*modePtr = 4
		default:
			*modePtr = 5
		}

		if (*flagsPtr)["hot"] {
			*outputPtr += (*scoresPtr)[mod4]
		}
		breakNow := *sumPtr > s.breakLimit
		if s.breakOnGE {
			breakNow = *sumPtr >= s.breakLimit
		}
		if breakNow {
			break
		}

		// for i%7 == 0 { continue }
		if mod7 != 0 {
			*countPtr++
		}

		i++
		mod5++
		if mod5 == 5 {
			mod5 = 0
		}
		mod7++
		if mod7 == 7 {
			mod7 = 0
		}
		mod4++
		if mod4 == 4 {
			mod4 = 0
		}
	}
	*iPtr = i
	return nil
}

// NewForComplexStateStmt creates a fused kernel for benchmark complex-state loop pattern.
func NewForComplexStateStmt(init New, iOffset uintptr, limit int, sumOffset, countOffset, modeOffset, outputOffset, flagsOffset, scoresOffset uintptr, breakLimit int, breakOnGE bool) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		stmt := &complexStateForStmt{
			iOffset:      iOffset,
			limit:        limit,
			sumOffset:    sumOffset,
			countOffset:  countOffset,
			modeOffset:   modeOffset,
			outputOffset: outputOffset,
			flagsOffset:  flagsOffset,
			scoresOffset: scoresOffset,
			breakLimit:   breakLimit,
			breakOnGE:    breakOnGE,
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
