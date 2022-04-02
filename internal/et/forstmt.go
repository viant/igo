package et

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"unsafe"
)

//NewForStmt creates a for stmt
func NewForStmt(cond *Operand, init, post, body New) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		var err error
		forKind := 0
		stmt := &forStmt{}
		if cond != nil {
			forKind = 1
			stmt.cond, err = cond.NewOperand(control)
		}
		if init != nil {
			forKind |= 1 << 1
			if stmt.init, err = init(control); err != nil {
				return nil, err
			}
		}
		if post != nil {
			forKind |= 1 << 2
			if stmt.post, err = post(control); err != nil {
				return nil, err
			}
		}

		if body != nil {
			forKind |= 1 << 3
			if stmt.body, err = body(control); err != nil {
				return nil, err
			}
		}
		switch forKind {
		case 0xF:
			return stmt.computeFor, nil
		}
		return stmt.compute, nil
	}, nil
}

type forStmt struct {
	init internal.Compute
	cond *exec.Operand
	body internal.Compute
	post internal.Compute
}

func (s *forStmt) computeFor(ptr unsafe.Pointer) unsafe.Pointer {
	for s.init(ptr); *(*bool)(s.cond.Compute(ptr)); s.post(ptr) {
		s.body(ptr)
	}
	return nil
}

func (s *forStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	if s.init != nil {
		s.init(ptr)
	}
	flow := internal.AsFlow(ptr)
	for {
		if s.cond != nil {
			if !*(*bool)(s.cond.Compute(ptr)) {
				return nil
			}
		}
		r := s.body(ptr)
		if flow.HasControl() {
			if flow.HasBreak() {
				flow.Reset()
				return nil
			}
			if flow.HasContinue() {
				flow.Reset()
				continue
			}
			if flow.HasReturn() {
				return r
			}
		}
	}
}

type forRorBlock2Stmt struct {
	s1, s2 internal.Compute
}

func (b *forRorBlock2Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.s1(ptr)
	if internal.AsFlow(ptr).HasBlockControl() {
		return r
	}
	return b.s2(ptr)
}

func newRorBlock2Stmt(s1, s2 internal.Compute) *forRorBlock2Stmt {
	return &forRorBlock2Stmt{s1: s1, s2: s2}
}

type forRorBlock3Stmt struct {
	g2 forRorBlock2Stmt
	s3 internal.Compute
}

func (b *forRorBlock3Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g2.compute(ptr)
	if internal.AsFlow(ptr).HasBlockControl() {
		return r
	}
	return b.s3(ptr)
}

func newRorBlock3Stmt(s1, s2, s3 internal.Compute) *forRorBlock3Stmt {
	return &forRorBlock3Stmt{g2: *newRorBlock2Stmt(s1, s2), s3: s3}
}

type forRorBlock4Stmt struct {
	g3 forRorBlock3Stmt
	s4 internal.Compute
}

func (b *forRorBlock4Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g3.compute(ptr)
	if internal.AsFlow(ptr).HasBlockControl() {
		return r
	}
	return b.s4(ptr)
}

func newRorBlock4Stmt(s1, s2, s3, s4 internal.Compute) *forRorBlock4Stmt {
	return &forRorBlock4Stmt{g3: *newRorBlock3Stmt(s1, s2, s3), s4: s4}
}

type forRorBlock5Stmt struct {
	g4 forRorBlock4Stmt
	s5 internal.Compute
}

func (b *forRorBlock5Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g4.compute(ptr)
	if internal.AsFlow(ptr).HasBlockControl() {
		return r
	}
	return b.s5(ptr)
}

func newRorBlock5Stmt(s1, s2, s3, s4, s5 internal.Compute) *forRorBlock5Stmt {
	return &forRorBlock5Stmt{g4: *newRorBlock4Stmt(s1, s2, s3, s4), s5: s5}
}

type forRorBlock6Stmt struct {
	g5 forRorBlock5Stmt
	s6 internal.Compute
}

func (b *forRorBlock6Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g5.compute(ptr)
	if internal.AsFlow(ptr).HasBlockControl() {
		return r
	}
	return b.s6(ptr)
}

func newRorBlock6Stmt(s1, s2, s3, s4, s5, s6 internal.Compute) *forRorBlock6Stmt {
	return &forRorBlock6Stmt{g5: *newRorBlock5Stmt(s1, s2, s3, s4, s5), s6: s6}
}

type forRorBlock7Stmt struct {
	g6 forRorBlock6Stmt
	s7 internal.Compute
}

func (b *forRorBlock7Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g6.compute(ptr)
	if internal.AsFlow(ptr).HasBlockControl() {
		return r
	}
	return b.s7(ptr)
}

func newRorBlock7Stmt(s1, s2, s3, s4, s5, s6, s7 internal.Compute) *forRorBlock7Stmt {
	return &forRorBlock7Stmt{g6: *newRorBlock6Stmt(s1, s2, s3, s4, s5, s6), s7: s7}
}

type forRorBlock8Stmt struct {
	g7 forRorBlock7Stmt
	s8 internal.Compute
}

func (b *forRorBlock8Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g7.compute(ptr)
	if internal.AsFlow(ptr).HasBlockControl() {
		return r
	}
	return b.s8(ptr)
}

func newRorBlock8Stmt(s1, s2, s3, s4, s5, s6, s7, s8 internal.Compute) *forRorBlock8Stmt {
	return &forRorBlock8Stmt{g7: *newRorBlock7Stmt(s1, s2, s3, s4, s5, s6, s7), s8: s8}
}

type forRorBlockStmt struct {
	statements []internal.Compute
}

func (b *forRorBlockStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	var result unsafe.Pointer
	flow := internal.AsFlow(ptr)
	for _, stmt := range b.statements {
		result = stmt(ptr)
		if flow.HasBlockControl() {
			return result
		}
	}
	return result
}

func newForBlockStmt(stmts []internal.Compute) internal.Compute {
	switch len(stmts) {
	case 0:
		return nop
	case 1:
		return stmts[0]
	case 2:
		return newRorBlock2Stmt(stmts[0], stmts[1]).compute
	case 3:
		return newRorBlock3Stmt(stmts[0], stmts[1], stmts[2]).compute
	case 4:
		return newRorBlock4Stmt(stmts[0], stmts[1], stmts[2], stmts[3]).compute
	case 5:
		return newRorBlock5Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4]).compute
	case 6:
		return newRorBlock6Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4], stmts[5]).compute
	case 7:
		return newRorBlock7Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4], stmts[5], stmts[6]).compute
	case 8:
		return newRorBlock8Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4], stmts[5], stmts[6], stmts[7]).compute
	}
	grp := &forRorBlockStmt{statements: stmts}
	return grp.compute
}
