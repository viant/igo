package exec

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

//New function to create variable
type New func() *Variables

//Variables represents execution plan variables
type Variables struct {
	pointer   unsafe.Pointer
	value     interface{}
	index     map[string]uint16
	selectors []*Selector
}

//Pointer returns variables pointer
func (v *Variables) Pointer() unsafe.Pointer {
	return v.pointer
}

//SetValue sets value for
func (v *Variables) SetValue(name string, value interface{}) error {
	sel, err := v.Selector(name)
	if err != nil {
		return err
	}
	sel.SetValue(sel.Upstream(v.pointer), value)
	return nil
}

//SetValueAt set values at position
func (v *Variables) SetValueAt(idx int, value interface{}) {
	sel := v.selectors[idx]
	sel.SetValue(sel.Upstream(v.pointer), value)
}

//ValueAt returns value for variable index
func (v *Variables) ValueAt(idx int) interface{} {
	sel := v.selectors[idx]
	return sel.Interface(sel.Upstream(v.pointer))
}

//Value returns value for supplied name or error
func (v *Variables) Value(name string) (interface{}, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return nil, err
	}
	return sel.Interface(sel.Upstream(v.pointer)), nil
}

//SetInt set int value
func (v *Variables) SetInt(name string, value int) error {
	sel, err := v.Selector(name)
	if err != nil {
		return err
	}
	if sel.Pathway == PathwayDirect {
		sel.SetInt(v.pointer, value)
		return nil
	}
	sel.SetInt(sel.Upstream(v.pointer), value)
	return nil
}

//Index returns variable index for specified name
func (v *Variables) Index(name string) (int, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return 0, err
	}
	return int(sel.Pos), nil
}

//Int returns int or error for supplied variable name
func (v *Variables) Int(name string) (int, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return 0, err
	}
	return sel.Int(sel.Upstream(v.pointer)), nil
}

//SetString set value for supplied variable name or error if variable name is invalid
func (v *Variables) SetString(name string, value string) error {
	sel, err := v.Selector(name)
	if err != nil {
		return err
	}
	sel.SetString(sel.Upstream(v.pointer), value)
	return nil
}

//SetStringAt set value for supplied variable index
func (v *Variables) SetStringAt(index int, value string) error {
	sel := v.selectors[index]
	sel.SetString(sel.Upstream(v.pointer), value)
	return nil
}

//String set string value for supplied variable name or error if variable name is invalid
func (v *Variables) String(name string) (string, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return "", err
	}
	return sel.String(sel.Upstream(v.pointer)), nil
}

//SetBool set boolean value for supplied variable name or error if variable name is invalid
func (v *Variables) SetBool(name string, value bool) error {
	sel, err := v.Selector(name)
	if err != nil {
		return err
	}
	sel.SetBool(sel.Upstream(v.pointer), value)
	return nil
}

//Bool returns boolean value or error for supplied variable name
func (v *Variables) Bool(name string) (bool, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return false, err
	}
	return sel.Bool(sel.Upstream(v.pointer)), nil
}

//SetBoolAt set boolean value at selor position
func (v *Variables) SetBoolAt(idx int, value bool) {
	sel := v.selectors[idx]
	sel.SetBool(sel.Upstream(v.pointer), value)
}

//SetFloat64 set float64 value for supplied variable name or error if variable name is invalid
func (v *Variables) SetFloat64(name string, value float64) error {
	sel, err := v.Selector(name)
	if err != nil {
		return err
	}
	sel.SetFloat64(sel.Upstream(v.pointer), value)
	return nil
}

//Float64 returns float64 value or error for supplied variable name
func (v *Variables) Float64(name string) (float64, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return 0.0, err
	}
	return sel.Float64(sel.Upstream(v.pointer)), nil
}

//Selector returns selector for a supploed name or error
func (v *Variables) Selector(name string) (*Selector, error) {
	idx, ok := v.index[name]
	if !ok {
		return nil, fmt.Errorf("undefined %v", name)
	}
	return v.selectors[idx], nil
}

//VariablesNew creates variables
func VariablesNew(dest reflect.Type, selectors []*Selector) New {
	var index = make(map[string]uint16)
	if len(selectors) == 0 {
		fields := xunsafe.NewStruct(dest).Fields
		for i := range fields {
			sel := &Selector{Field: &fields[i], ID: fields[i].Name}
			sel.Pos = uint16(len(selectors))
			sel.Pathway = SelectorPathway(sel)
			selectors = append(selectors, sel)
		}
	}
	trueSelPos := -1
	for i := range selectors {
		sel := selectors[i]
		index[sel.ID] = sel.Pos
		if sel.ID == "true" {
			trueSelPos = int(sel.Pos)
		}
	}
	return func() *Variables {
		value := reflect.New(dest).Interface()
		ret := &Variables{
			pointer:   xunsafe.AsPointer(value),
			value:     value,
			index:     index,
			selectors: selectors,
		}
		if trueSelPos != -1 {
			ret.SetBoolAt(trueSelPos, true)
		}
		return ret
	}
}
