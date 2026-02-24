package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/igo/internal/et"
	"go/ast"
	"go/token"
	"reflect"
	"unsafe"
)

func (s *Scope) compileUnaryExpr(unaryExpr *ast.UnaryExpr) (et.New, reflect.Type, error) {
	switch unaryExpr.Op {
	case token.AND:
		if expr, ok := unaryExpr.X.(*ast.CompositeLit); ok {
			return s.compileCompositeLiteral(unaryExpr.Op, expr)
		}
	case token.NOT:
		operand, err := s.assembleOperand(unaryExpr.X, false)
		if err != nil {
			return nil, nil, err
		}
		if operand.Type == nil || operand.Type.Kind() != reflect.Bool {
			return nil, nil, fmt.Errorf("invalid ! operand type: %v", operand.Type)
		}
		return func(control *et.Control) (internal.Compute, error) {
			op, err := operand.NewOperand(control)
			if err != nil {
				return nil, err
			}
			expr := &notExpr{operand: op}
			return expr.compute, nil
		}, boolType, nil
	case token.ADD, token.SUB, token.XOR:
		operand, err := s.assembleOperand(unaryExpr.X, false)
		if err != nil {
			return nil, nil, err
		}
		if operand.Type == nil {
			return nil, nil, fmt.Errorf("invalid unary operand type: %v", operand.Type)
		}
		k := operand.Type.Kind()
		switch unaryExpr.Op {
		case token.ADD:
			switch k {
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Uint8, reflect.Float32, reflect.Float64:
				return func(control *et.Control) (internal.Compute, error) {
					op, err := operand.NewOperand(control)
					if err != nil {
						return nil, err
					}
					expr := &plusExpr{operand: op}
					return expr.compute, nil
				}, operand.Type.Type(), nil
			}
			return nil, nil, fmt.Errorf("invalid + operand type: %v", operand.Type)
		case token.SUB:
			switch k {
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Uint8, reflect.Float32, reflect.Float64:
				return func(control *et.Control) (internal.Compute, error) {
					op, err := operand.NewOperand(control)
					if err != nil {
						return nil, err
					}
					expr := &numericUnaryExpr{operand: op, op: unaryExpr.Op, kind: k}
					return expr.compute, nil
				}, operand.Type.Type(), nil
			}
			return nil, nil, fmt.Errorf("invalid - operand type: %v", operand.Type)
		case token.XOR:
			switch k {
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Uint8:
				return func(control *et.Control) (internal.Compute, error) {
					op, err := operand.NewOperand(control)
					if err != nil {
						return nil, err
					}
					expr := &numericUnaryExpr{operand: op, op: unaryExpr.Op, kind: k}
					return expr.compute, nil
				}, operand.Type.Type(), nil
			}
			return nil, nil, fmt.Errorf("invalid ^ operand type: %v", operand.Type)
		}
	}
	return nil, nil, fmt.Errorf("unsupported type: %T", unaryExpr.X)
}

type notExpr struct {
	operand *exec.Operand
	value   bool
}

func (n *notExpr) compute(ptr unsafe.Pointer) unsafe.Pointer {
	n.value = !*(*bool)(n.operand.Compute(ptr))
	return unsafe.Pointer(&n.value)
}

type plusExpr struct {
	operand *exec.Operand
}

func (n *plusExpr) compute(ptr unsafe.Pointer) unsafe.Pointer {
	return n.operand.Compute(ptr)
}

type numericUnaryExpr struct {
	operand *exec.Operand
	op      token.Token
	kind    reflect.Kind
	i       int
	i64     int64
	u       uint
	u64     uint64
	u8      uint8
	f32     float32
	f64     float64
}

func (n *numericUnaryExpr) compute(ptr unsafe.Pointer) unsafe.Pointer {
	valuePtr := n.operand.Compute(ptr)
	switch n.kind {
	case reflect.Int:
		v := *(*int)(valuePtr)
		if n.op == token.SUB {
			n.i = -v
		} else {
			n.i = ^v
		}
		return unsafe.Pointer(&n.i)
	case reflect.Int64:
		v := *(*int64)(valuePtr)
		if n.op == token.SUB {
			n.i64 = -v
		} else {
			n.i64 = ^v
		}
		return unsafe.Pointer(&n.i64)
	case reflect.Uint:
		v := *(*uint)(valuePtr)
		if n.op == token.SUB {
			n.u = -v
		} else {
			n.u = ^v
		}
		return unsafe.Pointer(&n.u)
	case reflect.Uint64:
		v := *(*uint64)(valuePtr)
		if n.op == token.SUB {
			n.u64 = -v
		} else {
			n.u64 = ^v
		}
		return unsafe.Pointer(&n.u64)
	case reflect.Uint8:
		v := *(*uint8)(valuePtr)
		if n.op == token.SUB {
			n.u8 = -v
		} else {
			n.u8 = ^v
		}
		return unsafe.Pointer(&n.u8)
	case reflect.Float32:
		v := *(*float32)(valuePtr)
		n.f32 = -v
		return unsafe.Pointer(&n.f32)
	case reflect.Float64:
		v := *(*float64)(valuePtr)
		n.f64 = -v
		return unsafe.Pointer(&n.f64)
	}
	return valuePtr
}
