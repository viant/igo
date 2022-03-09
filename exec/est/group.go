package est

import (
	"github.com/viant/igo/exec"
	"unsafe"
)

func nop(ptr unsafe.Pointer) unsafe.Pointer {
	return ptr
}

type group2Stmt struct {
	s1, s2 exec.Compute
}

func (b *group2Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	b.s1(ptr)
	return b.s2(ptr)
}

func newGroup2Stmt(s1, s2 exec.Compute) *group2Stmt {
	return &group2Stmt{s1: s1, s2: s2}
}

type group3Stmt struct {
	g2 group2Stmt
	s3 exec.Compute
}

func (b *group3Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	b.g2.compute(ptr)
	return b.s3(ptr)
}

func newGroup3Stmt(s1, s2, s3 exec.Compute) *group3Stmt {
	return &group3Stmt{g2: *newGroup2Stmt(s1, s2), s3: s3}
}

type group4Stmt struct {
	g3 group3Stmt
	s4 exec.Compute
}

func (b *group4Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	b.g3.compute(ptr)
	return b.s4(ptr)
}

func newGroup4Stmt(s1, s2, s3, s4 exec.Compute) *group4Stmt {
	return &group4Stmt{g3: *newGroup3Stmt(s1, s2, s3), s4: s4}
}

type group5Stmt struct {
	g4 group4Stmt
	s5 exec.Compute
}

func (b *group5Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	b.g4.compute(ptr)
	return b.s5(ptr)
}

func newGroup5Stmt(s1, s2, s3, s4, s5 exec.Compute) *group5Stmt {
	return &group5Stmt{g4: *newGroup4Stmt(s1, s2, s3, s4), s5: s5}
}

type group6Stmt struct {
	g5 group5Stmt
	s6 exec.Compute
}

func (b *group6Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	b.g5.compute(ptr)
	return b.s6(ptr)
}

func newGroup6Stmt(s1, s2, s3, s4, s5, s6 exec.Compute) *group6Stmt {
	return &group6Stmt{g5: *newGroup5Stmt(s1, s2, s3, s4, s5), s6: s6}
}

type group7Stmt struct {
	g6 group6Stmt
	s7 exec.Compute
}

func (b *group7Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	b.g6.compute(ptr)
	return b.s7(ptr)
}

func newGroup7Stmt(s1, s2, s3, s4, s5, s6, s7 exec.Compute) *group7Stmt {
	return &group7Stmt{g6: *newGroup6Stmt(s1, s2, s3, s4, s5, s6), s7: s7}
}

type group8Stmt struct {
	g7 group7Stmt
	s8 exec.Compute
}

func (b *group8Stmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	b.g7.compute(ptr)
	return b.s8(ptr)
}

func newGroup8Stmt(s1, s2, s3, s4, s5, s6, s7, s8 exec.Compute) *group8Stmt {
	return &group8Stmt{g7: *newGroup7Stmt(s1, s2, s3, s4, s5, s6, s7), s8: s8}
}

type groupStmt struct {
	statements []exec.Compute
}

func (b *groupStmt) compute(ptr unsafe.Pointer) unsafe.Pointer {
	var result unsafe.Pointer
	for _, stmt := range b.statements {
		result = stmt(ptr)
	}
	return result
}

//NewGroupStmt crates group node
func NewGroupStmt(newStatementsFn []New, stmt bool) New {
	return func(control *Control) (exec.Compute, error) {
		var stmts = make([]exec.Compute, len(newStatementsFn))
		var err error
		for i := range newStatementsFn {
			newFn := newStatementsFn[i]
			if stmts[i], err = newFn(control); err != nil {
				return nil, err
			}
		}
		return newGroupStmt(stmts), nil
	}
}

func newGroupStmt(stmts []exec.Compute) exec.Compute {
	switch len(stmts) {
	case 0:
		return nop
	case 1:
		return stmts[0]
	case 2:
		return newGroup2Stmt(stmts[0], stmts[1]).compute
	case 3:
		return newGroup3Stmt(stmts[0], stmts[1], stmts[2]).compute
	case 4:
		return newGroup4Stmt(stmts[0], stmts[1], stmts[2], stmts[3]).compute
	case 5:
		return newGroup5Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4]).compute
	case 6:
		return newGroup6Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4], stmts[5]).compute
	case 7:
		return newGroup7Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4], stmts[5], stmts[6]).compute
	case 8:
		return newGroup8Stmt(stmts[0], stmts[1], stmts[2], stmts[3], stmts[4], stmts[5], stmts[6], stmts[7]).compute
	}
	grp := &groupStmt{statements: stmts}
	return grp.compute
}
