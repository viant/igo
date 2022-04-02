package plan

import (
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
	"reflect"
)

func (s *Scope) compileCompositeLiteral(unaryToken token.Token, compositeLit *ast.CompositeLit) (et.New, reflect.Type, error) {
	litType, err := s.discoverType(compositeLit.Type)
	if err != nil {
		return nil, nil, err
	}
	destType := litType
	if unaryToken == token.AND {
		destType = reflect.PtrTo(litType)
	}
	var operands = make([]*et.Operand, len(compositeLit.Elts))
	for i, elt := range compositeLit.Elts {
		switch expr := elt.(type) {
		case *ast.KeyValueExpr:
			if operands[i], err = s.assembleOperand(expr.Value, false); err != nil {
				return nil, nil, err
			}
			operands[i].Key = stringifyExpr(expr.Key, 0)

		case *ast.UnaryExpr:
			if operands[i], err = s.assembleOperand(expr.X, false); err != nil {
				return nil, nil, err
			}
			operands[i].Idx = i
		default:
			if operands[i], err = s.assembleOperand(expr, false); err != nil {
				return nil, nil, err
			}
			operands[i].Idx = i
		}
	}
	return et.NewComposite(destType, operands), destType, nil
}
