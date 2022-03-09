package est

import (
	"github.com/viant/igo/exec"
	"unsafe"
)

type block2Stmt struct {
	s1, s2 exec.Compute
}

func (b *block2Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.s1(ptr)
	if exec.AsFlow(ptr).HasReturn() {
		return r
	}
	return b.s2(ptr)
}

func newBlock2Stmt(s1, s2 exec.Compute) *block2Stmt {
	return &block2Stmt{s1: s1, s2: s2}
}

type block3Stmt struct {
	g2 block2Stmt
	s3 exec.Compute
}

func (b *block3Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g2.compute(ptr)
	if exec.AsFlow(ptr).HasReturn() {
		return r
	}
	return b.s3(ptr)
}

func newBlock3Stmt(s1, s2, s3 exec.Compute) *block3Stmt {
	return &block3Stmt{g2: *newBlock2Stmt(s1, s2), s3: s3}
}

type block4Stmt struct {
	g3 block3Stmt
	s4 exec.Compute
}

func (b *block4Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g3.compute(ptr)
	if exec.AsFlow(ptr).HasReturn() {
		return r
	}
	return b.s4(ptr)
}

func newBlock4Stmt(s1, s2, s3, s4 exec.Compute) *block4Stmt {
	return &block4Stmt{g3: *newBlock3Stmt(s1, s2, s3), s4: s4}
}

type block5Stmt struct {
	g4 block4Stmt
	s5 exec.Compute
}

func (b *block5Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g4.compute(ptr)
	if exec.AsFlow(ptr).HasReturn() {
		return r
	}
	return b.s5(ptr)
}

func newBlock5Stmt(s1, s2, s3, s4, s5 exec.Compute) *block5Stmt {
	return &block5Stmt{g4: *newBlock4Stmt(s1, s2, s3, s4), s5: s5}
}

type block6Stmt struct {
	g5 block5Stmt
	s6 exec.Compute
}

func (b *block6Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g5.compute(ptr)
	if exec.AsFlow(ptr).HasReturn() {
		return r
	}
	return b.s6(ptr)
}

func newBlock6Stmt(s1, s2, s3, s4, s5, s6 exec.Compute) *block6Stmt {
	return &block6Stmt{g5: *newBlock5Stmt(s1, s2, s3, s4, s5), s6: s6}
}

type block7Stmt struct {
	g6 block6Stmt
	s7 exec.Compute
}

func (b *block7Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g6.compute(ptr)
	if exec.AsFlow(ptr).HasReturn() {
		return r
	}
	return b.s7(ptr)
}

func newBlock7Stmt(s1, s2, s3, s4, s5, s6, s7 exec.Compute) *block7Stmt {
	return &block7Stmt{g6: *newBlock6Stmt(s1, s2, s3, s4, s5, s6), s7: s7}
}

type block8Stmt struct {
	g7 block7Stmt
	s8 exec.Compute
}

func (b *block8Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	r := b.g7.compute(ptr)
	if exec.AsFlow(ptr).HasReturn() {
		return r
	}

	return b.s8(ptr)
}

func newBlock8Stmt(s1, s2, s3, s4, s5, s6, s7, s8 exec.Compute) *block8Stmt {
	return &block8Stmt{g7: *newBlock7Stmt(s1, s2, s3, s4, s5, s6, s7), s8: s8}
}

type blockStmt struct {
	statements []exec.Compute
}

func (b *blockStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	var result unsafe.Pointer
	for _, stmt := range b.statements {
		result = stmt(ptr)
		if exec.AsFlow(ptr).HasReturn() {
			return result
		}
	}
	return result
}

//NewBlockStmt crates block node
func NewBlockStmt(newStatementsFn []New, forStmt bool) New {
	return func(control *Control) (exec.Compute, error) {
		var stmts = make([]exec.Compute, len(newStatementsFn))
		var err error
		for i := range newStatementsFn {
			newFn := newStatementsFn[i]
			if stmts[i], err = newFn(control); err != nil {
				return nil, err
			}
		}

		if forStmt {
			if control.HasControl() {
				return newForBlockStmt(stmts), nil
			}
			return newGroupStmt(stmts), nil
		}
		if control.HasReturn() {
			return newBlockStmt(stmts), nil
		}
		return newGroupStmt(stmts), nil

	}
}

func newBlockStmt(stmts []exec.Compute) exec.Compute {
	switch len(stmts) {
	case 0:
		return nop
	case 1:
		return stmts[0]
	case 2:
		return newBlock2Stmt(stmts[0], stmts[1]).compute
	case 3:
		return newBlock3Stmt(stmts[0], stmts[1], stmts[2]).compute
	case 4:
		return newBlock4Stmt(stmts[0], stmts[1], stmts[2], stmts[3]).compute
	case 5:
		return newBlock5Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4]).compute
	case 6:
		return newBlock6Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4], stmts[5]).compute
	case 7:
		return newBlock7Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4], stmts[5], stmts[6]).compute
	case 8:
		return newBlock8Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4], stmts[5], stmts[6], stmts[7]).compute
	}
	grp := &blockStmt{statements: stmts}
	return grp.compute
}
