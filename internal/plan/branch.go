package plan

import (
	"fmt"
	"github.com/viant/igo/internal"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
)

func (s *Scope) compileBranchStmt(actual *ast.BranchStmt) (et.New, error) {
	switch actual.Tok {
	case token.BREAK:
		s.Control.Concat(internal.RtBreak)
	case token.CONTINUE:
		s.Control.Concat(internal.RtContinue)
	default:
		return nil, fmt.Errorf("not yet sypported %s", actual.Tok)
	}
	return et.NewBranchStmt(actual.Tok, actual.Label)
}
