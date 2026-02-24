package et

import (
	"fmt"
	"github.com/viant/igo/internal"
	"go/token"
	"unsafe"
)

// NewBranchStmt create branch statment
func NewBranchStmt(tkn token.Token, labelID uint32) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		switch tkn {
		case token.BREAK:
			return func(ptr unsafe.Pointer) unsafe.Pointer {
				internal.AsFlow(ptr).SetBranch(internal.RtBreak, labelID)
				return nil
			}, nil

		case token.CONTINUE:
			return func(ptr unsafe.Pointer) unsafe.Pointer {
				internal.AsFlow(ptr).SetBranch(internal.RtContinue, labelID)
				return nil
			}, nil
		case token.GOTO:
			return func(ptr unsafe.Pointer) unsafe.Pointer {
				internal.AsFlow(ptr).SetBranch(internal.RtGoto, labelID)
				return nil
			}, nil
		case token.FALLTHROUGH:
			return func(ptr unsafe.Pointer) unsafe.Pointer {
				internal.AsFlow(ptr).SetBranch(internal.RtFallthrough, 0)
				return nil
			}, nil
		}
		return nil, fmt.Errorf("unsupported branch token: %s", tkn)
	}, nil
}
