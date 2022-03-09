package est

import (
	"fmt"
	"github.com/viant/igo/exec"
	"go/ast"
	"go/token"
	"unsafe"
)

//NewBranchStmt create branch statment
func NewBranchStmt(tkn token.Token, label *ast.Ident) (New, error) {
	return func(control *Control) (exec.Compute, error) {
		switch tkn {
		case token.BREAK:
			return func(ptr unsafe.Pointer) unsafe.Pointer {
				*(*exec.Flow)(ptr) |= exec.RtBreak
				return nil
			}, nil

		case token.CONTINUE:
			return func(ptr unsafe.Pointer) unsafe.Pointer {
				*(*exec.Flow)(ptr) |= exec.RtContinue
				return nil
			}, nil
		}
		return nil, fmt.Errorf("unsupported branch token: %s", tkn)
	}, nil
}
