package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"go/token"
	"reflect"
	"unsafe"
)

var trueValue = true
var falseValue = false
var trueValuePointer = unsafe.Pointer(&trueValue)
var falseValuePointer = unsafe.Pointer(&falseValue)

// NewBinaryExpr creates binary expr
func NewBinaryExpr(op token.Token, operands ...*Operand) (New, reflect.Type) {
	opType := operands[0].Type.Type()
	destType := opType
	switch op {
	case token.EQL, token.NEQ, token.GTR, token.LSS, token.GEQ, token.LEQ, token.LOR, token.LAND:
		destType = boolType
	}
	assignToken := op

	return func(exec *Control) (internal.Compute, error) {
		ops := Operands(operands)
		if err := ops.Validate(); err != nil {
			return nil, err
		}
		switch assignToken {
		case token.REM:
			return newRem(opType, operands, exec)
		case token.ADD:
			return newAdd(opType, operands, exec)
		case token.ADD_ASSIGN:
			return newAddAssign(opType, operands, exec)
		case token.SUB:
			return newSub(opType, operands, exec)
		case token.SUB_ASSIGN:
			return newSubAssign(opType, operands, exec)
		case token.MUL:
			return newMul(opType, operands, exec)
		case token.MUL_ASSIGN:
			return newMulAssign(opType, operands, exec)
		case token.QUO:
			return newQuo(opType, operands, exec)
		case token.QUO_ASSIGN:
			return newQuoAssign(opType, operands, exec)
		case token.GTR:
			return newGtr(opType, operands, exec)
		case token.LSS:
			return newLss(opType, operands, exec)
		case token.GEQ:
			return newGeq(opType, operands, exec)
		case token.LEQ:
			return newLeq(opType, operands, exec)
		case token.SHR:
			return newShr(opType, operands, exec)
		case token.SHL:
			return newShl(opType, operands, exec)
		case token.AND, token.LAND:
			return newAnd(opType, operands, exec)
		case token.AND_ASSIGN:
			return newAndAssign(opType, operands, exec)
		case token.OR, token.LOR:
			return newOr(opType, operands, exec)
		case token.OR_ASSIGN:
			return newOrAssign(opType, operands, exec)
		case token.XOR:
			return newXor(opType, operands, exec)
		case token.XOR_ASSIGN:
			return newXorAssign(opType, operands, exec)
		case token.AND_NOT:
			return newAndNot(opType, operands, exec)
		case token.AND_NOT_ASSIGN:
			return newAndNotAssign(opType, operands, exec)
		case token.EQL:
			return newEql(operands, exec)
		case token.NEQ:
			return newNeq(operands, exec)
		}
		return Unsupported(fmt.Sprintf("unsupported  %v %v %v", opType.Name(), op, opType.Name()))
	}, destType
}

func newEql(operands Operands, control *Control) (internal.Compute, error) {
	if xOffset, y, ok := intDirectRightConst(operands); ok {
		return intEqlConstRight(xOffset, y), nil
	}
	if yOffset, x, ok := intDirectLeftConst(operands); ok {
		return intEqlConstLeft(x, yOffset), nil
	}

	cmp, err := newComparison(operands, control)
	if err != nil {
		return nil, err
	}
	if operands.HasNil() {
		return cmp.NilEql, nil
	}

	switch operands.pathway() {
	case exec.PathwayDirect:
		return cmp.DirectEql, nil
	default:
		return cmp.Eql, nil
	}
}

func newNeq(operands Operands, control *Control) (internal.Compute, error) {
	if xOffset, y, ok := intDirectRightConst(operands); ok {
		return intNeqConstRight(xOffset, y), nil
	}
	if yOffset, x, ok := intDirectLeftConst(operands); ok {
		return intNeqConstLeft(x, yOffset), nil
	}
	cmp, err := newComparison(operands, control)
	if err != nil {
		return nil, err
	}
	if operands.HasNil() {
		return cmp.NilNeq, nil
	}

	if err != nil {
		return nil, err
	}
	switch operands.pathway() {
	case exec.PathwayDirect:
		return cmp.DirectNeq, nil
	default:
		return cmp.Neq, nil
	}
}

func newShr(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intShr, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intShr, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v >> %v", opType.Name(), opType.Name()))

}

func newShl(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intShl, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intShr, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v << %v", opType.Name(), opType.Name()))
}

func newAnd(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	if opType.Kind() == reflect.Int {
		if xOffset, y, ok := intDirectRightConst(operands); ok {
			return intAndConstRight(xOffset, y, operands[zOp].Offset()), nil
		}
	}
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAnd, nil
		case reflect.Bool:
			return expr.lAnd, nil

		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAnd, nil
		case reflect.Bool:
			return expr.and, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v & %v", opType.Name(), opType.Name()))
}

func newAndAssign(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAndAssign, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAndAssign, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v & %v", opType.Name(), opType.Name()))
}

func newOr(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intOr, nil
		case reflect.Bool:
			return expr.lOr, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intOr, nil
		case reflect.Bool:
			return expr.or, nil

		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v | %v", opType.Name(), opType.Name()))
}

func newOrAssign(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intOrAssign, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intOrAssign, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v | %v", opType.Name(), opType.Name()))
}

func newXor(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intXor, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intXor, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v ^ %v", opType.Name(), opType.Name()))
}

func newXorAssign(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intXorAssign, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intXorAssign, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v ^ %v", opType.Name(), opType.Name()))
}

func newAndNot(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAndNot, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAndNot, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v &^ %v", opType.Name(), opType.Name()))
}

func newAndNotAssign(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAndNotAssign, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAndNotAssign, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v &^ %v", opType.Name(), opType.Name()))
}

func newRem(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	if opType.Kind() == reflect.Int {
		if xOffset, y, ok := intDirectRightConst(operands); ok {
			return intRemConstRight(xOffset, y, operands[zOp].Offset()), nil
		}
	}
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intRem, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intRem, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v + %v", opType.Name(), opType.Name()))
}

func newAdd(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAdd, nil
		case reflect.Float64:
			return expr.float64Add, nil
		case reflect.String:
			return expr.stringAdd, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAdd, nil
		case reflect.Float64:
			return expr.float64Add, nil
		case reflect.String:
			return expr.stringAdd, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v + %v", opType.Name(), opType.Name()))
}

func newAddAssign(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAddAssign, nil
		case reflect.Float64:
			return expr.float64AddAssign, nil
		case reflect.String:
			return expr.stringAddAssign, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intAddAssign, nil
		case reflect.Float64:
			return expr.float64AddAssign, nil
		case reflect.String:
			return expr.stringAddAssign, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v + %v", opType.Name(), opType.Name()))
}

func newSub(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intSub, nil
		case reflect.Float64:
			return expr.float64Sub, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intSub, nil
		case reflect.Float64:
			return expr.float64Sub, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v - %v p:%v", opType.Name(), opType.Name(), operands.pathway()))
}

func newSubAssign(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intSubAssign, nil
		case reflect.Float64:
			return expr.float64SubAssign, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intSubAssign, nil
		case reflect.Float64:
			return expr.float64SubAssign, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v - %v p:%v", opType.Name(), opType.Name(), operands.pathway()))
}

func newMul(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intMul, nil
		case reflect.Float64:
			return expr.float64Mul, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intMul, nil
		case reflect.Float64:
			return expr.float64Mul, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v * %v", opType.Name(), opType.Name()))
}

func newMulAssign(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intMulAssign, nil
		case reflect.Float64:
			return expr.float64MulAssign, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intMulAssign, nil
		case reflect.Float64:
			return expr.float64MulAssign, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v * %v", opType.Name(), opType.Name()))
}

func newQuo(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intQuo, nil
		case reflect.Float64:
			return expr.float64Quo, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intQuo, nil
		case reflect.Float64:
			return expr.float64Quo, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v * %v", opType.Name(), opType.Name()))
}

func newQuoAssign(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intQuoAssign, nil
		case reflect.Float64:
			return expr.float64QuoAssign, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intQuoAssign, nil
		case reflect.Float64:
			return expr.float64QuoAssign, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v * %v", opType.Name(), opType.Name()))
}

func newGtr(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	if opType.Kind() == reflect.Int {
		if xOffset, y, ok := intDirectRightConst(operands); ok {
			return intGtrConstRight(xOffset, y), nil
		}
		if yOffset, x, ok := intDirectLeftConst(operands); ok {
			return intGtrConstLeft(x, yOffset), nil
		}
	}
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intGtr, nil
		case reflect.Float64:
			return expr.float64Gtr, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intGtr, nil
		case reflect.Float64:
			return expr.float64Gtr, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v * %v", opType.Name(), opType.Name()))
}

func newLss(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	if opType.Kind() == reflect.Int {
		if xOffset, y, ok := intDirectRightConst(operands); ok {
			return intLssConstRight(xOffset, y), nil
		}
		if yOffset, x, ok := intDirectLeftConst(operands); ok {
			return intLssConstLeft(x, yOffset), nil
		}
	}
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intLss, nil
		case reflect.Float64:
			return expr.float64Lss, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intLss, nil
		case reflect.Float64:
			return expr.float64Lss, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v * %v", opType.Name(), opType.Name()))
}

func intLiteralValue(op *Operand) (int, bool) {
	if op == nil || op.Selector != nil || op.New != nil {
		return 0, false
	}
	value, ok := op.Value.(int)
	return value, ok
}

func intDirectRightConst(operands Operands) (uintptr, int, bool) {
	x := operands[xOp]
	if x == nil || x.Selector == nil || x.Pathway != exec.PathwayDirect || x.Type.Kind() != reflect.Int {
		return 0, 0, false
	}
	y, ok := intLiteralValue(operands[yOp])
	if !ok {
		return 0, 0, false
	}
	return x.Offset(), y, true
}

func intDirectLeftConst(operands Operands) (uintptr, int, bool) {
	y := operands[yOp]
	if y == nil || y.Selector == nil || y.Pathway != exec.PathwayDirect || y.Type.Kind() != reflect.Int {
		return 0, 0, false
	}
	x, ok := intLiteralValue(operands[xOp])
	if !ok {
		return 0, 0, false
	}
	return y.Offset(), x, true
}

func intLssConstRight(xOffset uintptr, y int) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if *(*int)(unsafe.Pointer(uintptr(ptr) + xOffset)) < y {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intLssConstLeft(x int, yOffset uintptr) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if x < *(*int)(unsafe.Pointer(uintptr(ptr) + yOffset)) {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intLeqConstRight(xOffset uintptr, y int) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if *(*int)(unsafe.Pointer(uintptr(ptr) + xOffset)) <= y {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intLeqConstLeft(x int, yOffset uintptr) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if x <= *(*int)(unsafe.Pointer(uintptr(ptr) + yOffset)) {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intGtrConstRight(xOffset uintptr, y int) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if *(*int)(unsafe.Pointer(uintptr(ptr) + xOffset)) > y {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intGtrConstLeft(x int, yOffset uintptr) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if x > *(*int)(unsafe.Pointer(uintptr(ptr) + yOffset)) {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intGeqConstRight(xOffset uintptr, y int) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if *(*int)(unsafe.Pointer(uintptr(ptr) + xOffset)) >= y {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intGeqConstLeft(x int, yOffset uintptr) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if x >= *(*int)(unsafe.Pointer(uintptr(ptr) + yOffset)) {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intEqlConstRight(xOffset uintptr, y int) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if *(*int)(unsafe.Pointer(uintptr(ptr) + xOffset)) == y {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intEqlConstLeft(x int, yOffset uintptr) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if x == *(*int)(unsafe.Pointer(uintptr(ptr) + yOffset)) {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intNeqConstRight(xOffset uintptr, y int) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if *(*int)(unsafe.Pointer(uintptr(ptr) + xOffset)) != y {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intNeqConstLeft(x int, yOffset uintptr) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if x != *(*int)(unsafe.Pointer(uintptr(ptr) + yOffset)) {
			return trueValuePointer
		}
		return falseValuePointer
	}
}

func intAndConstRight(xOffset uintptr, y int, zOffset uintptr) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		z := unsafe.Pointer(uintptr(ptr) + zOffset)
		*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + xOffset)) & y
		return z
	}
}

func intRemConstRight(xOffset uintptr, y int, zOffset uintptr) internal.Compute {
	return func(ptr unsafe.Pointer) unsafe.Pointer {
		z := unsafe.Pointer(uintptr(ptr) + zOffset)
		*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + xOffset)) % y
		return z
	}
}

func newGeq(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	if opType.Kind() == reflect.Int {
		if xOffset, y, ok := intDirectRightConst(operands); ok {
			return intGeqConstRight(xOffset, y), nil
		}
		if yOffset, x, ok := intDirectLeftConst(operands); ok {
			return intGeqConstLeft(x, yOffset), nil
		}
	}
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intGeq, nil
		case reflect.Float64:
			return expr.float64Geq, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intGeq, nil
		case reflect.Float64:
			return expr.float64Geq, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v * %v", opType.Name(), opType.Name()))
}

func newLeq(opType reflect.Type, operands Operands, control *Control) (internal.Compute, error) {
	if opType.Kind() == reflect.Int {
		if xOffset, y, ok := intDirectRightConst(operands); ok {
			return intLeqConstRight(xOffset, y), nil
		}
		if yOffset, x, ok := intDirectLeftConst(operands); ok {
			return intLeqConstLeft(x, yOffset), nil
		}
	}
	switch operands.pathway() {
	case exec.PathwayDirect:
		expr := operands.directBinaryExpr()
		switch opType.Kind() {
		case reflect.Int:
			return expr.intLeq, nil
		case reflect.Float64:
			return expr.float64Leq, nil
		}
	default:
		expr, err := operands.binaryExpr(control)
		if err != nil {
			return nil, err
		}
		switch opType.Kind() {
		case reflect.Int:
			return expr.intLeq, nil
		case reflect.Float64:
			return expr.float64Leq, nil
		}
	}
	return Unsupported(fmt.Sprintf("unsupported  %v * %v", opType.Name(), opType.Name()))
}

type comparison struct {
	x *exec.Operand
	y *exec.Operand
}

func (e *comparison) NilEql(ptr unsafe.Pointer) unsafe.Pointer {
	x := e.x.Compute(ptr)
	if x == nil || *(*unsafe.Pointer)(x) == nil {
		return trueValuePointer
	}
	return falseValuePointer
}

func (e *comparison) NilNeq(ptr unsafe.Pointer) unsafe.Pointer {
	x := e.x.Compute(ptr)
	if x != nil && *(*unsafe.Pointer)(x) != nil {
		return trueValuePointer
	}
	return falseValuePointer
}

func (e *comparison) DirectEql(ptr unsafe.Pointer) unsafe.Pointer {
	x := e.x.Value(e.x.Compute(ptr))
	y := e.y.Value(e.y.Compute(ptr))
	result := falseValuePointer
	if x == y {
		result = trueValuePointer
	}
	return result
}

func (e *comparison) DirectNeq(ptr unsafe.Pointer) unsafe.Pointer {
	x := e.x.Value(e.x.Compute(ptr))
	y := e.y.Value(e.y.Compute(ptr))
	result := falseValuePointer
	if x != y {
		result = trueValuePointer
	}
	return result
}

func (e *comparison) Eql(ptr unsafe.Pointer) unsafe.Pointer {
	x := e.x.Value(e.x.Compute(ptr))
	y := e.y.Value(e.y.Compute(ptr))
	result := falseValuePointer
	if x == y {
		result = trueValuePointer
	}
	return result
}

func (e *comparison) Neq(ptr unsafe.Pointer) unsafe.Pointer {
	x := e.x.Value(e.x.Compute(ptr))
	y := e.y.Value(e.y.Compute(ptr))
	result := falseValuePointer
	if x != y {
		result = trueValuePointer
	}
	return result
}

func newComparison(operands Operands, control *Control) (*comparison, error) {
	result := &comparison{}
	var err error
	if operands.HasNil() {
		x := operands.NonNilOperand()
		if result.x, err = x.NewOperand(control); err != nil {
			return nil, err
		}
		return result, nil
	}
	x := operands[xOp]
	y := operands[yOp]
	if result.x, err = x.NewOperand(control); err != nil {
		return nil, err
	}
	if result.y, err = y.NewOperand(control); err != nil {
		return nil, err
	}
	return result, err
}

type binaryExpr struct {
	xNode   *exec.Operand
	yNode   *exec.Operand
	zOffset uintptr
}

func (e *binaryExpr) z(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + e.zOffset)
}

func (e *binaryExpr) x(ptr unsafe.Pointer) unsafe.Pointer {
	return e.xNode.Compute(ptr)
}

func (e *binaryExpr) y(ptr unsafe.Pointer) unsafe.Pointer {
	return e.yNode.Compute(ptr)
}

func (e *binaryExpr) stringAdd(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*string)(z) = *(*string)(e.x(ptr)) + *(*string)(e.y(ptr))
	return z
}

func (e *binaryExpr) intRem(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) % *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intAdd(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) + *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intAddAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x += *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) intSub(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) - *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intSubAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x -= *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) intMul(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	x := e.x(ptr)
	y := e.y(ptr)
	*(*int)(z) = *(*int)(x) * *(*int)(y)
	return z
}

func (e *binaryExpr) intMulAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x *= *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) intQuo(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) / *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intQuoAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x /= *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) intGtr(ptr unsafe.Pointer) unsafe.Pointer {
	if *(*int)(e.x(ptr)) > *(*int)(e.y(ptr)) {
		return trueValuePointer
	}
	return falseValuePointer
}

func (e *binaryExpr) intLss(ptr unsafe.Pointer) unsafe.Pointer {
	if *(*int)(e.x(ptr)) < *(*int)(e.y(ptr)) {
		return trueValuePointer
	}
	return falseValuePointer
}

func (e *binaryExpr) intGeq(ptr unsafe.Pointer) unsafe.Pointer {
	if *(*int)(e.x(ptr)) >= *(*int)(e.y(ptr)) {
		return trueValuePointer
	}
	return falseValuePointer
}

func (e *binaryExpr) intLeq(ptr unsafe.Pointer) unsafe.Pointer {
	if *(*int)(e.x(ptr)) <= *(*int)(e.y(ptr)) {
		return trueValuePointer
	}
	return falseValuePointer
}

func (e *binaryExpr) intShr(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) >> *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intShl(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) << *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intAnd(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) & *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intAndAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x &= *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) intOr(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) | *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intOrAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x |= *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) intXor(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) ^ *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intXorAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x ^= *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) intAndNot(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) &^ *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) intAndNotAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x &^= *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) float64Add(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*float64)(z) = *(*float64)(e.x(ptr)) + *(*float64)(e.y(ptr))
	return z
}

func (e *binaryExpr) float64AddAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*float64)(e.x(ptr))
	*x += *(*float64)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) stringAddAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*string)(e.x(ptr))
	*x += *(*string)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) float64Sub(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*float64)(z) = *(*float64)(e.x(ptr)) - *(*float64)(e.y(ptr))
	return z
}

func (e *binaryExpr) float64SubAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*float64)(e.x(ptr))
	*x -= *(*float64)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) float64Mul(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*float64)(z) = *(*float64)(e.x(ptr)) * *(*float64)(e.y(ptr))
	return z
}

func (e *binaryExpr) float64MulAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x *= *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) float64Quo(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(e.x(ptr)) / *(*int)(e.y(ptr))
	return z
}

func (e *binaryExpr) float64QuoAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(e.x(ptr))
	*x /= *(*int)(e.y(ptr))
	return unsafe.Pointer(x)
}

func (e *binaryExpr) float64Gtr(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*int)(e.x(ptr)) > *(*int)(e.y(ptr)) {
		result = trueValuePointer
	}
	return result
}

func (e *binaryExpr) float64Lss(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*int)(e.x(ptr)) < *(*int)(e.y(ptr)) {
		result = trueValuePointer
	}
	return result
}

func (e *binaryExpr) float64Geq(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*int)(e.x(ptr)) >= *(*int)(e.y(ptr)) {
		result = trueValuePointer
	}
	return result
}

func (e *binaryExpr) float64Leq(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*int)(e.x(ptr)) <= *(*int)(e.y(ptr)) {
		result = falseValuePointer
	}
	return result
}

func (e *binaryExpr) and(ptr unsafe.Pointer) unsafe.Pointer {
	if !*(*bool)(e.x(ptr)) {
		return falseValuePointer
	}
	if *(*bool)(e.y(ptr)) {
		return trueValuePointer
	}
	return falseValuePointer
}

func (e *binaryExpr) or(ptr unsafe.Pointer) unsafe.Pointer {
	if *(*bool)(e.x(ptr)) {
		return trueValuePointer
	}
	if *(*bool)(e.y(ptr)) {
		return trueValuePointer
	}
	return falseValuePointer
}

type directBinaryExpr struct {
	xOffset uintptr
	yOffset uintptr
	zOffset uintptr
}

func (e *directBinaryExpr) z(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + e.zOffset)
}

func (e *directBinaryExpr) stringAdd(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*string)(z) = *(*string)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) + *(*string)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intRem(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) % *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intAdd(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) + *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intSub(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) - *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intMul(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) * *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intQuo(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) / *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) stringAddAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*string)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x += *(*string)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) intAddAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	y := *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	*x += y
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) intSubAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x -= *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) intMulAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x *= *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) intQuoAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x /= *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) intGtr(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) > *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		result = trueValuePointer
	}
	return result
}

func (e *directBinaryExpr) intLss(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) < *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		result = trueValuePointer
	}
	return result
}

func (e *directBinaryExpr) intGeq(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) >= *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		result = trueValuePointer
	}
	return result
}

func (e *directBinaryExpr) intLeq(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) <= *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		result = trueValuePointer
	}
	return result
}

func (e *directBinaryExpr) intShr(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) >> *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intShl(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) << *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intAnd(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) & *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intAndAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x &= *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) intOr(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) | *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intOrAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x |= *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) intXor(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) ^ *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intXorAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x ^= *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) intAndNot(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*int)(z) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) &^ *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) intAndNotAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x &^= *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) float64Add(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*float64)(z) = *(*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) + *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) float64AddAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x += *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) float64Sub(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*float64)(z) = *(*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) - *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) float64SubAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x -= *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) float64Mul(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*float64)(z) = *(*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) * *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) float64MulAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x *= *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) float64Quo(ptr unsafe.Pointer) unsafe.Pointer {
	z := e.z(ptr)
	*(*float64)(z) = *(*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) / *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return z
}

func (e *directBinaryExpr) float64QuoAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := (*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	*x /= *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return unsafe.Pointer(x)
}

func (e *directBinaryExpr) float64Gtr(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) > *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		result = trueValuePointer
	}
	return result
}

func (e *directBinaryExpr) float64Lss(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) < *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		result = trueValuePointer
	}
	return result
}

func (e *directBinaryExpr) float64Geq(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) >= *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		result = trueValuePointer
	}
	return result
}

func (e *directBinaryExpr) float64Leq(ptr unsafe.Pointer) unsafe.Pointer {
	result := falseValuePointer
	if *(*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) <= *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		result = trueValuePointer
	}
	return result
}

func (e *directBinaryExpr) lAnd(ptr unsafe.Pointer) unsafe.Pointer {
	if !*(*bool)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) {
		return falseValuePointer
	}
	if *(*bool)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		return trueValuePointer
	}
	return falseValuePointer
}

func (e *directBinaryExpr) lOr(ptr unsafe.Pointer) unsafe.Pointer {
	if *(*bool)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) {
		return trueValuePointer
	}
	if *(*bool)(unsafe.Pointer(uintptr(ptr) + e.yOffset)) {
		return trueValuePointer
	}
	return falseValuePointer
}
