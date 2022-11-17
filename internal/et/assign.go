package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"reflect"
	"unsafe"
)

//NewCallExprAssign create an assign expression
func NewCallExprAssign(caller exec.Caller, args []*Operand, dest []*Operand) (New, error) {
	expr := &callAssign{
		callExpr: &callExpr{caller: caller},
	}
	argOperands := Operands(args)
	destOperands := Operands(dest)
	isDirect := destOperands.pathway() == exec.PathwayDirect

	return func(control *Control) (internal.Compute, error) {
		var err error
		if expr.dest, err = destOperands.operands(control); err != nil {
			return nil, err
		}
		if len(expr.dest) == 1 {
			expr.offset = expr.dest[0].Offset
		}

		if expr.callExpr.operands, err = argOperands.operands(control); err != nil {
			return nil, err
		}

		switch len(dest) {
		case 1:
			if isDirect {
				switch dest[0].Kind() {
				case reflect.Int:
					return expr.computeDirectInt, nil
				case reflect.String:
					return expr.computeDirectString, nil
				case reflect.Bool:
					return expr.computeDirectBool, nil
				case reflect.Float64:
					return expr.computeDirectFloat64, nil
				}
			}
			switch dest[0].Kind() {
			case reflect.Int:
				return expr.computeInt, nil
			case reflect.String:
				return expr.computeString, nil
			case reflect.Bool:
				return expr.computeBool, nil
			case reflect.Float64:
				return expr.computeFloat64, nil
			}
			return expr.compute, nil
		case 2:
			return expr.computeR2, nil
		case 3:
			return expr.computeR3, nil
		case 4:
			return expr.computeR4, nil
		case 5:
			return expr.computeR5, nil
		case 6:
			return expr.computeR6, nil
		}
		return nil, fmt.Errorf("too many returns variables: %v, not supported yet", len(dest))
	}, nil
}

//NewAssignExpr create an assign expression
func NewAssignExpr(trackPos []uint16, ops ...*Operand) New {
	operands := Operands(ops)
	opType := operands[1].Type.Type()
	isDirect := operands.pathway() == exec.PathwayDirect
	trackable := len(trackPos) > 0
	return func(exec *Control) (internal.Compute, error) {
		assignExpr, err := operands.assignExpr(exec)
		assignExpr.trackPos = trackPos
		if err != nil {
			return nil, err
		}
		if isDirect {
			switch opType.Kind() {
			case reflect.Int:
				if trackable {
					return assignExpr.directIntTackedAssign, nil
				}
				return assignExpr.directIntAssign, nil
			case reflect.Float64:
				if trackable {
					return assignExpr.directFloat64TackedAssign, nil
				}
				return assignExpr.directFloat64Assign, nil
			case reflect.String:
				if trackable {
					return assignExpr.directStringTackedAssign, nil
				}
				return assignExpr.directStringAssign, nil
			case reflect.Bool:

				if trackable {
					return assignExpr.directBoolTackedAssign, nil
				}
				return assignExpr.directBoolAssign, nil
			default:
				if opType.ConvertibleTo(errType) {
					return assignExpr.directErrorAssign, nil
				}
				if trackable {
					return assignExpr.directTrackedAssign, nil
				}
				return assignExpr.directAssign, nil
			}
		}
		switch opType.Kind() {
		case reflect.Int:
			if trackable {
				return assignExpr.intTrackedAssign, nil
			}
			return assignExpr.intAssign, nil
		case reflect.Float64:
			if trackable {
				return assignExpr.float64TrackedAssign, nil
			}
			return assignExpr.float64Assign, nil
		case reflect.String:
			if trackable {
				return assignExpr.stringTrackedAssign, nil
			}
			return assignExpr.stringAssign, nil
		case reflect.Bool:
			if trackable {
				return assignExpr.boolTrackedAssign, nil
			}
			return assignExpr.boolAssign, nil
		default:
			if opType.ConvertibleTo(errType) {
				return assignExpr.errorAssign, nil
			}
			if trackable {
				return assignExpr.trackedAssign, nil
			}
			return assignExpr.assign, nil
		}
	}
}

type callAssign struct {
	dest     []*exec.Operand
	trackPos []int
	offset   uintptr
	callExpr *callExpr
}

func (a *callAssign) computeInt(ptr unsafe.Pointer) unsafe.Pointer {
	result := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	destPtr := a.dest[0].Compute(ptr)
	*(*int)(destPtr) = *(*int)(result)
	return result
}

func (a *callAssign) computeString(ptr unsafe.Pointer) unsafe.Pointer {
	result := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	destPtr := a.dest[0].Compute(ptr)
	*(*string)(destPtr) = *(*string)(result)
	return result
}

func (a *callAssign) computeBool(ptr unsafe.Pointer) unsafe.Pointer {
	result := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	destPtr := a.dest[0].Compute(ptr)
	*(*bool)(destPtr) = *(*bool)(result)
	return result
}

func (a *callAssign) computeFloat64(ptr unsafe.Pointer) unsafe.Pointer {
	result := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	destPtr := a.dest[0].Compute(ptr)
	*(*float64)(destPtr) = *(*float64)(result)
	return result
}

func (a *callAssign) computeDirectInt(ptr unsafe.Pointer) unsafe.Pointer {
	result := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	destPtr := unsafe.Pointer(uintptr(ptr) + a.dest[0].Offset)
	*(*int)(destPtr) = *(*int)(result)
	return result
}

func (a *callAssign) computeDirectString(ptr unsafe.Pointer) unsafe.Pointer {
	result := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	destPtr := unsafe.Pointer(uintptr(ptr) + a.dest[0].Offset)
	*(*string)(destPtr) = *(*string)(result)
	return result
}

func (a *callAssign) computeDirectBool(ptr unsafe.Pointer) unsafe.Pointer {
	result := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	destPtr := unsafe.Pointer(uintptr(ptr) + a.dest[0].Offset)
	*(*bool)(destPtr) = *(*bool)(result)
	return result
}

func (a *callAssign) computeDirectFloat64(ptr unsafe.Pointer) unsafe.Pointer {
	result := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	destPtr := unsafe.Pointer(uintptr(ptr) + a.dest[0].Offset)
	*(*float64)(destPtr) = *(*float64)(result)
	return result
}

func (a *callAssign) compute(ptr unsafe.Pointer) unsafe.Pointer {
	result := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	a.dest[0].SetValuePtr(ptr, result)
	return nil
}

func (a *callAssign) computeR2(ptr unsafe.Pointer) unsafe.Pointer {
	ret := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	results := *(*[2]unsafe.Pointer)(ret)
	a.dest[0].SetValuePtr(ptr, results[0])
	a.dest[1].SetValuePtr(ptr, results[1])
	return nil
}

func (a *callAssign) computeR3(ptr unsafe.Pointer) unsafe.Pointer {
	ret := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	results := *(*[3]unsafe.Pointer)(ret)
	a.dest[0].SetValuePtr(ptr, results[0])
	a.dest[1].SetValuePtr(ptr, results[1])
	a.dest[2].SetValuePtr(ptr, results[2])
	return nil
}

func (a *callAssign) computeR4(ptr unsafe.Pointer) unsafe.Pointer {
	ret := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	results := *(*[4]unsafe.Pointer)(ret)
	a.dest[0].SetValuePtr(ptr, results[0])
	a.dest[1].SetValuePtr(ptr, results[1])
	a.dest[2].SetValuePtr(ptr, results[2])
	a.dest[3].SetValuePtr(ptr, results[3])
	return nil
}

func (a *callAssign) computeR5(ptr unsafe.Pointer) unsafe.Pointer {
	ret := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	results := *(*[5]unsafe.Pointer)(ret)
	a.dest[0].SetValuePtr(ptr, results[0])
	a.dest[1].SetValuePtr(ptr, results[1])
	a.dest[2].SetValuePtr(ptr, results[2])
	a.dest[3].SetValuePtr(ptr, results[3])
	a.dest[4].SetValuePtr(ptr, results[4])
	return nil
}

func (a *callAssign) computeR6(ptr unsafe.Pointer) unsafe.Pointer {
	ret := a.callExpr.caller.Call(ptr, a.callExpr.operands)
	results := *(*[6]unsafe.Pointer)(ret)
	a.dest[0].SetValuePtr(ptr, results[0])
	a.dest[0].SetValuePtr(ptr, results[0])
	a.dest[1].SetValuePtr(ptr, results[1])
	a.dest[2].SetValuePtr(ptr, results[2])
	a.dest[3].SetValuePtr(ptr, results[3])
	a.dest[4].SetValuePtr(ptr, results[4])
	a.dest[5].SetValuePtr(ptr, results[5])
	return nil
}

type assignExpr struct {
	trackPos []uint16
	x        *exec.Selector
	xOffset  uintptr
	y        *exec.Operand
	yOffset  uintptr
}

func (e *assignExpr) setMutation(ptr unsafe.Pointer) {
	tracker := *(**exec.Tracker)(unsafe.Pointer(uintptr(ptr) + 8))
	tracker.Set(e.trackPos)
}

func (e *assignExpr) directIntTackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	dest := (*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	src := *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	if *dest == src {
		return unsafe.Pointer(dest)
	}
	*dest = src
	e.setMutation(ptr)
	return nil
}

func (e *assignExpr) directIntAssign(ptr unsafe.Pointer) unsafe.Pointer {
	*(*int)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) = *(*int)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return nil
}

func (e *assignExpr) directStringTackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	dest := (*string)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	src := *(*string)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	if *dest == src {
		return unsafe.Pointer(dest)
	}
	*dest = src
	e.setMutation(ptr)
	return nil
}

func (e *assignExpr) directStringAssign(ptr unsafe.Pointer) unsafe.Pointer {
	*(*string)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) = *(*string)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return nil
}

func (e *assignExpr) directBoolTackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	dest := (*bool)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	src := *(*bool)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	if *dest == src {
		return unsafe.Pointer(dest)
	}
	*dest = src
	e.setMutation(ptr)
	return nil
}

func (e *assignExpr) directBoolAssign(ptr unsafe.Pointer) unsafe.Pointer {
	*(*bool)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) = *(*bool)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return nil
}

func (e *assignExpr) directFloat64TackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	dest := (*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset))
	src := *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	if *dest == src {
		return unsafe.Pointer(dest)
	}
	*dest = src
	e.setMutation(ptr)
	return nil
}

func (e *assignExpr) directFloat64Assign(ptr unsafe.Pointer) unsafe.Pointer {
	*(*float64)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) = *(*float64)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return nil
}

func (e *assignExpr) directErrorAssign(ptr unsafe.Pointer) unsafe.Pointer {
	*(*error)(unsafe.Pointer(uintptr(ptr) + e.xOffset)) = *(*error)(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	return nil
}

func (e *assignExpr) directTrackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	value := e.y.Interface(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	e.x.SetValue(unsafe.Pointer(uintptr(ptr)+e.xOffset), value)
	e.setMutation(ptr)
	return nil
}

func (e *assignExpr) directAssign(ptr unsafe.Pointer) unsafe.Pointer {
	value := e.y.Interface(unsafe.Pointer(uintptr(ptr) + e.yOffset))
	e.x.SetValue(unsafe.Pointer(uintptr(ptr)+e.xOffset), value)
	return nil
}

func (e *assignExpr) intTrackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := *(*int)(e.y.Compute(ptr))
	upstream := e.x.Upstream(ptr)
	if x == (e.x.Int(upstream)) {
		return nil
	}
	e.x.SetInt(upstream, x)
	e.setMutation(ptr)
	return nil
}

func (e *assignExpr) intAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := *(*int)(e.y.Compute(ptr))
	e.x.SetInt(e.x.Upstream(ptr), x)
	return nil
}

func (e *assignExpr) stringTrackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := *(*string)(e.y.Compute(ptr))
	upstream := e.x.Upstream(ptr)
	if x == (e.x.String(upstream)) {
		return nil
	}
	e.x.SetString(upstream, x)
	e.setMutation(ptr)
	return nil
}

func (e *assignExpr) stringAssign(ptr unsafe.Pointer) unsafe.Pointer {
	e.x.SetString(e.x.Upstream(ptr), *(*string)(e.y.Compute(ptr)))
	return nil
}

func (e *assignExpr) boolTrackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := *(*bool)(e.y.Compute(ptr))
	upstream := e.x.Upstream(ptr)
	if x == (e.x.Bool(upstream)) {
		return nil
	}
	e.x.SetBool(upstream, x)
	e.setMutation(ptr)
	return nil
}

func (e *assignExpr) boolAssign(ptr unsafe.Pointer) unsafe.Pointer {
	e.x.SetBool(e.x.Upstream(ptr), *(*bool)(e.y.Compute(ptr)))
	return nil
}

func (e *assignExpr) float64TrackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	x := *(*float64)(e.y.Compute(ptr))
	upstream := e.x.Upstream(ptr)
	if x == (e.x.Float64(upstream)) {
		return nil
	}
	e.x.SetFloat64(upstream, x)
	e.setMutation(ptr)
	return nil
}

func (e *assignExpr) float64Assign(ptr unsafe.Pointer) unsafe.Pointer {
	e.x.SetFloat64(e.x.Upstream(ptr), *(*float64)(e.y.Compute(ptr)))
	return nil
}

func (e *assignExpr) errorAssign(ptr unsafe.Pointer) unsafe.Pointer {
	e.x.SetError(e.x.Upstream(ptr), *(*error)(e.y.Compute(ptr)))
	return nil
}

func (e *assignExpr) assign(ptr unsafe.Pointer) unsafe.Pointer {
	yPtr := e.y.Compute(ptr)
	value := e.y.Interface(yPtr)
	e.x.SetValue(e.x.Upstream(ptr), value)
	return nil
}

func (e *assignExpr) trackedAssign(ptr unsafe.Pointer) unsafe.Pointer {
	yPtr := e.y.Compute(ptr)
	value := e.y.Interface(yPtr)
	e.x.SetValue(e.x.Upstream(ptr), value)
	e.setMutation(ptr)
	return nil
}
