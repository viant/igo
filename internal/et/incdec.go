package et

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"go/token"
	"unsafe"
)

type incStmt struct {
	*exec.Operand
}

func (s *incStmt) directInc(ptr unsafe.Pointer) unsafe.Pointer {
	x := unsafe.Pointer(uintptr(ptr) + s.Offset)
	*(*int)(x)++
	return x
}

func (s *incStmt) directDec(ptr unsafe.Pointer) unsafe.Pointer {
	x := unsafe.Pointer(uintptr(ptr) + s.Offset)
	*(*int)(x)--
	return x
}

func (s *incStmt) inc(ptr unsafe.Pointer) unsafe.Pointer {
	x := s.Compute(ptr)
	*(*int)(x)++
	return x
}

func (s *incStmt) dec(ptr unsafe.Pointer) unsafe.Pointer {
	x := s.Compute(ptr)
	*(*int)(x)--
	return x
}

//NewIncDec creates inc/dec stmt
func NewIncDec(tok token.Token, op *Operand) New {
	return func(control *Control) (internal.Compute, error) {
		operand, err := op.NewOperand(control)
		if err != nil {
			return nil, err
		}
		isDirect := op.Pathway == exec.PathwayDirect
		stmt := &incStmt{Operand: operand}
		if tok == token.INC {
			if isDirect {
				return stmt.directInc, nil
			}
			return stmt.inc, nil
		}
		if isDirect {
			return stmt.directDec, nil
		}
		return stmt.dec, nil
	}
}
