package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/exec/est"
	"go/ast"
	"go/token"
)

func (s *Scope) compileBranchStmt(actual *ast.BranchStmt) (est.New, error) {
	switch actual.Tok {
	case token.BREAK:
		s.Control.Concat(exec.RtBreak)
	case token.CONTINUE:
		s.Control.Concat(exec.RtContinue)
	default:
		return nil, fmt.Errorf("not yet sypported %s", actual.Tok)
	}
	return est.NewBranchStmt(actual.Tok, actual.Label)
}
