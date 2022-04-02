package plan

import (
	"fmt"
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/internal/exec/et"
	"go/ast"
	"go/token"
)

func (s *Scope) compileBranchStmt(actual *ast.BranchStmt) (et.New, error) {
	switch actual.Tok {
	case token.BREAK:
		s.Control.Concat(exec.RtBreak)
	case token.CONTINUE:
		s.Control.Concat(exec.RtContinue)
	default:
		return nil, fmt.Errorf("not yet sypported %s", actual.Tok)
	}
	return et.NewBranchStmt(actual.Tok, actual.Label)
}
