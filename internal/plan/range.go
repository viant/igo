package plan

import (
	"fmt"
	"github.com/viant/igo/internal/exec/et"
	"github.com/viant/igo/exec"
	"go/ast"
	"reflect"
)

func (s *Scope) compileRangeStmt(rangeStmt *ast.RangeStmt) (et.New, error) {
	var x *et.Operand
	var key, value *exec.Selector
	var body et.New
	var err error
	if x, err = s.assembleOperand(rangeStmt.X, false); err != nil {
		return nil, err
	}
	if rangeStmt.Value != nil {
		if value, err = s.selector(rangeStmt.Value, true); err != nil {
			return nil, err
		}
	}
	if rangeStmt.Key != nil {
		if key, err = s.selector(rangeStmt.Key, true); err != nil {
			return nil, err
		}
	}
	xType := derefType(x.Selector.Type)
	valueType := derefType(xType.Elem())
	switch xType.Kind() {
	case reflect.Slice:
		if key != nil {
			s.adjust(key, intType)
		}
		if value != nil {
			s.adjust(value, reflect.PtrTo(valueType))
		}
	default:
		return nil, fmt.Errorf("range not supported for type %s", xType.String())
	}

	body, err = s.compileBlockStmt(rangeStmt.Body, true)
	return et.NewRange(x, key, value, body)
}
