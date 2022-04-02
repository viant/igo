package et

import (
	"fmt"
	"github.com/viant/igo/state"
	"github.com/viant/xunsafe"
)

const (
	xOp = 0
	yOp = 1
	zOp = 2
)

//Operands represents opertands
type Operands []*Operand

func (o Operands) pathway() state.Pathway {
	result := state.PathwayUndefined
	for _, candidate := range o {
		if candidate.Selector == nil {
			return state.PathwayUndefined
		}
		if candidate.Pathway >= result {
			result = candidate.Pathway
			continue
		}
		result = state.PathwayRef
	}
	return result
}

func (o Operands) directBinaryExpr() *directBinaryExpr {
	expr := &directBinaryExpr{xOffset: o[xOp].Offset(), yOffset: o[yOp].Offset(), zOffset: o[zOp].Offset()}
	return expr
}

func (o Operands) binaryExpr(control *Control) (*binaryExpr, error) {
	result := &binaryExpr{}
	var err error
	if result.xNode, err = o[xOp].NewOperand(control); err != nil {
		return nil, err
	}
	if result.yNode, err = o[yOp].NewOperand(control); err != nil {
		return nil, err
	}
	result.zOffset = o[zOp].Offset()
	return result, err
}

func (o Operands) assignExpr(control *Control) (*assignExpr, error) {
	result := &assignExpr{}
	var err error
	x := o[xOp]
	result.x = x.Selector
	if result.x == nil {
		return nil, fmt.Errorf("destSlice selector was nil")
	}
	result.xOffset = result.x.Offset()

	y := o[yOp]
	if result.y, err = y.NewOperand(control); err != nil {
		return nil, err
	}
	if y.Selector != nil {
		result.yOffset = y.Offset()
	}
	return result, err
}

func (o Operands) operands(control *Control) ([]*state.Operand, error) {
	var result = make([]*state.Operand, len(o))
	var err error
	for i, operand := range o {
		if result[i], err = operand.NewOperand(control); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (o Operands) selectors() []*state.Selector {
	var result = make([]*state.Selector, len(o))
	for i, operand := range o {
		result[i] = operand.Selector
	}
	return result
}

func (o Operands) types() []*xunsafe.Type {
	var result = make([]*xunsafe.Type, len(o))
	for i, operand := range o {
		result[i] = operand.Type
	}
	return result
}
