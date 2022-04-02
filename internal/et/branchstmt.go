package et

import (
	"fmt"
	"github.com/viant/igo/internal"
	"go/ast"
	"go/token"
	"unsafe"
)

//NewBranchStmt create branch statment
func NewBranchStmt(tkn token.Token, label *ast.Ident) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		switch tkn {
		case token.BREAK:
			return func(ptr unsafe.Pointer) unsafe.Pointer {
				*(*internal.Flow)(ptr) |= internal.RtBreak
				return nil
			}, nil

		case token.CONTINUE:
			return func(ptr unsafe.Pointer) unsafe.Pointer {
				*(*internal.Flow)(ptr) |= internal.RtContinue
				return nil
			}, nil
		}
		return nil, fmt.Errorf("unsupported branch token: %s", tkn)
	}, nil
}
