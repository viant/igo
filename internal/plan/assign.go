package plan

import (
	"fmt"
	"github.com/viant/igo/internal/et"
	"github.com/viant/xunsafe"
	"go/ast"
	"go/token"
	"reflect"
)

func (s *Scope) compileAssignStmt(assignStmt *ast.AssignStmt) (et.New, error) {
	if helper, handled, err := s.compileInlineFuncAssign(assignStmt); handled || err != nil {
		return helper, err
	}
	if assertAssign, handled, err := s.compileTypeAssertAssignStmt(assignStmt); handled || err != nil {
		return assertAssign, err
	}
	if sliceAssign, handled, err := s.compileSliceAssignStmt(assignStmt); handled || err != nil {
		return sliceAssign, err
	}
	if mapAssign, handled, err := s.compileMapAssignStmt(assignStmt); handled || err != nil {
		return mapAssign, err
	}

	define := assignStmt.Tok == token.DEFINE
	var (
		callExpr            *ast.CallExpr
		isCallExprAssigment bool
	)
	if len(assignStmt.Rhs) == 1 {
		callExpr, isCallExprAssigment = assignStmt.Rhs[0].(*ast.CallExpr)
	}
	rOperands, err := s.assembleOperands(assignStmt.Rhs, false)
	if err != nil {
		return nil, err
	}
	lOperands, err := s.assembleOperands(assignStmt.Lhs, true)
	if err != nil {
		return nil, err
	}
	if isCallExprAssigment {
		return s.compileCallExprAssign(callExpr, lOperands)
	}
	if len(lOperands) != len(rOperands) {
		return nil, fmt.Errorf("assignment count mismatch: %d = %d", len(lOperands), len(rOperands))
	}
	var group = make([]et.New, len(lOperands))
	for i := 0; i < len(lOperands); i++ {
		if group[i], err = s.compileAssign(lOperands[i], rOperands[i], assignStmt.Tok, define); err != nil {
			return nil, err
		}
	}
	return et.NewGroupStmt(group, false), nil
}

func (s *Scope) compileInlineFuncAssign(assignStmt *ast.AssignStmt) (et.New, bool, error) {
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return nil, false, nil
	}
	fnLit, ok := assignStmt.Rhs[0].(*ast.FuncLit)
	if !ok {
		return nil, false, nil
	}
	ident, ok := assignStmt.Lhs[0].(*ast.Ident)
	if !ok {
		return nil, true, fmt.Errorf("inline function assignment requires identifier lhs")
	}
	if assignStmt.Tok != token.DEFINE && assignStmt.Tok != token.ASSIGN {
		return nil, true, fmt.Errorf("unsupported inline function assignment token: %v", assignStmt.Tok)
	}
	fn, err := s.functionFromFuncLit(fnLit)
	if err != nil {
		return nil, true, err
	}
	s.RegisterFunc(ident.Name, fn)
	return et.NewNop(), true, nil
}

func (s *Scope) compileAssign(left, right *et.Operand, assignToken token.Token, define bool) (et.New, error) {
	switch assignToken {
	case token.ADD_ASSIGN, token.SUB_ASSIGN, token.AND_ASSIGN, token.OR_ASSIGN, token.AND_NOT_ASSIGN, token.MUL_ASSIGN, token.REM_ASSIGN, token.QUO_ASSIGN:
		// Compound assign writes back to left operand; reusing left as z keeps pathway direct
		// and avoids degrading to generic binary dispatch.
		newFn, _ := et.NewBinaryExpr(assignToken, left, right, left)
		return newFn, nil
	default:
		if left.Selector.Type == nil {
			s.adjust(left.Selector, right.Type.Type())
			left.Type = right.Type
		}
		xPos := s.trackerXPos(left.Selector)
		return et.NewAssignExpr(xPos, left, right), nil
	}
}

func (s *Scope) assembleOperands(expressions []ast.Expr, defined bool) ([]*et.Operand, error) {
	var result = make([]*et.Operand, len(expressions))
	var err error
	for i, expr := range expressions {
		if result[i], err = s.assembleOperand(expr, defined); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Scope) compileMapAssignStmt(assignStmt *ast.AssignStmt) (et.New, bool, error) {
	if len(assignStmt.Lhs) == 1 && len(assignStmt.Rhs) == 1 {
		indexExpr, ok := assignStmt.Lhs[0].(*ast.IndexExpr)
		if !ok {
			return nil, false, nil
		}
		mapOperand, keyOperand, mapType, isMap, err := s.assembleMapIndex(indexExpr)
		if err != nil || !isMap {
			return nil, isMap, err
		}
		if assignStmt.Tok != token.ASSIGN {
			return nil, true, fmt.Errorf("unsupported map assignment token: %v", assignStmt.Tok)
		}
		valueOperand, err := s.assembleOperand(assignStmt.Rhs[0], false)
		if err != nil {
			return nil, true, err
		}
		if valueOperand.Type == nil {
			return nil, true, fmt.Errorf("map assignment value type was empty")
		}
		if valueOperand.Type.Type() != mapType.Elem() {
			return nil, true, fmt.Errorf("incompatible map assignment type: %v cannot be assigned to %v", valueOperand.Type.Type(), mapType.Elem())
		}
		newFn, err := et.NewMapSet(mapOperand, keyOperand, valueOperand)
		return newFn, true, err
	}

	if len(assignStmt.Lhs) == 2 && len(assignStmt.Rhs) == 1 {
		indexExpr, ok := assignStmt.Rhs[0].(*ast.IndexExpr)
		if !ok {
			return nil, false, nil
		}
		mapOperand, keyOperand, mapType, isMap, err := s.assembleMapIndex(indexExpr)
		if err != nil || !isMap {
			return nil, isMap, err
		}
		if assignStmt.Tok != token.ASSIGN && assignStmt.Tok != token.DEFINE {
			return nil, true, fmt.Errorf("unsupported map value,ok token: %v", assignStmt.Tok)
		}
		lOperands, err := s.assembleOperands(assignStmt.Lhs, assignStmt.Tok == token.DEFINE)
		if err != nil {
			return nil, true, err
		}
		if lOperands[0].Selector.Type == nil {
			_ = s.adjust(lOperands[0].Selector, mapType.Elem())
			lOperands[0].Type = xunsafe.NewType(mapType.Elem())
		}
		if lOperands[1].Selector.Type == nil {
			_ = s.adjust(lOperands[1].Selector, boolType)
			lOperands[1].Type = xunsafe.NewType(boolType)
		}
		if lOperands[0].Type == nil {
			lOperands[0].Type = xunsafe.NewType(lOperands[0].Selector.Type)
		}
		if lOperands[1].Type == nil {
			lOperands[1].Type = xunsafe.NewType(lOperands[1].Selector.Type)
		}
		newFn, err := et.NewMapGet2Assign(mapOperand, keyOperand, lOperands[0], lOperands[1])
		return newFn, true, err
	}
	return nil, false, nil
}

func (s *Scope) compileSliceAssignStmt(assignStmt *ast.AssignStmt) (et.New, bool, error) {
	if len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		return nil, false, nil
	}
	indexExpr, ok := assignStmt.Lhs[0].(*ast.IndexExpr)
	if !ok {
		return nil, false, nil
	}
	sliceOperand, indexOperand, sliceType, isSlice, err := s.assembleSliceIndex(indexExpr)
	if err != nil || !isSlice {
		return nil, isSlice, err
	}
	if assignStmt.Tok != token.ASSIGN {
		return nil, true, nil
	}
	valueOperand, err := s.assembleOperand(assignStmt.Rhs[0], false)
	if err != nil {
		return nil, true, err
	}
	if valueOperand.Type == nil {
		return nil, true, fmt.Errorf("slice assignment value type was empty")
	}
	if valueOperand.Type.Type() != sliceType.Elem() {
		return nil, true, fmt.Errorf("incompatible slice assignment type: %v cannot be assigned to %v", valueOperand.Type.Type(), sliceType.Elem())
	}
	if sliceType == reflect.TypeOf([]int{}) {
		newFn, err := et.NewSliceSetInt(sliceOperand, indexOperand, valueOperand)
		return newFn, true, err
	}
	return nil, true, nil
}

func (s *Scope) assembleSliceIndex(indexExpr *ast.IndexExpr) (sliceOperand *et.Operand, indexOperand *et.Operand, sliceType reflect.Type, isSlice bool, err error) {
	sliceOperand, err = s.assembleOperand(indexExpr.X, false)
	if err != nil {
		return nil, nil, nil, false, err
	}
	if sliceOperand.Type == nil || sliceOperand.Type.Kind() != reflect.Slice {
		return nil, nil, nil, false, nil
	}
	sliceType = sliceOperand.Type.Type()
	indexOperand, err = s.assembleOperand(indexExpr.Index, false)
	if err != nil {
		return nil, nil, nil, false, err
	}
	if indexOperand.Type == nil {
		return nil, nil, nil, false, fmt.Errorf("slice index type was empty")
	}
	if indexOperand.Type.Type() != reflect.TypeOf(0) {
		return nil, nil, nil, false, fmt.Errorf("incompatible slice index type: %v, expected: int", indexOperand.Type.Type())
	}
	return sliceOperand, indexOperand, sliceType, true, nil
}

func (s *Scope) assembleMapIndex(indexExpr *ast.IndexExpr) (mapOperand *et.Operand, keyOperand *et.Operand, mapType reflect.Type, isMap bool, err error) {
	mapOperand, err = s.assembleOperand(indexExpr.X, false)
	if err != nil {
		return nil, nil, nil, false, err
	}
	if mapOperand.Type == nil || mapOperand.Type.Kind() != reflect.Map {
		return nil, nil, nil, false, nil
	}
	mapType = mapOperand.Type.Type()
	keyOperand, err = s.assembleOperand(indexExpr.Index, false)
	if err != nil {
		return nil, nil, nil, false, err
	}
	if keyOperand.Type == nil {
		return nil, nil, nil, false, fmt.Errorf("map key type was empty")
	}
	if keyOperand.Type.Type() != mapType.Key() {
		return nil, nil, nil, false, fmt.Errorf("incompatible map key type: %v, expected: %v", keyOperand.Type.Type(), mapType.Key())
	}
	return mapOperand, keyOperand, mapType, true, nil
}
