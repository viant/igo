package plan

import (
	"fmt"
	"github.com/viant/igo/internal"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
)

func (s *Scope) compileBranchStmt(actual *ast.BranchStmt) (et.New, error) {
	var labelID uint32
	if actual.Label != nil {
		var ok bool
		labelID, ok = s.labelID(actual.Label.Name)
		if !ok {
			return nil, fmt.Errorf("undefined label: %s", actual.Label.Name)
		}
	}
	switch actual.Tok {
	case token.BREAK:
		s.Control.Concat(internal.RtBreak)
	case token.CONTINUE:
		s.Control.Concat(internal.RtContinue)
	case token.GOTO:
		s.Control.Concat(internal.RtGoto)
	case token.FALLTHROUGH:
		s.Control.Concat(internal.RtFallthrough)
	default:
		return nil, fmt.Errorf("not yet sypported %s", actual.Tok)
	}
	return et.NewBranchStmt(actual.Tok, labelID)
}
