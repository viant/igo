package exec

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

//New function to create variable
type New func() *State

//State represents execution plan variables
type State struct {
	pointer   unsafe.Pointer
	tracker   *Tracker
	value     interface{}
	index     map[string]uint16
	selectors []*Selector
}

//Pointer returns variables pointer
func (v *State) Pointer() unsafe.Pointer {
	return v.pointer
}


//SetValue sets value for
func (v *State) SetValue(name string, value interface{}) error {
	sel, err := v.Selector(name)
	if err != nil {
		return err
	}
	sel.SetValue(sel.Upstream(v.pointer), value)
	return nil
}

//SetValueAt set values at position
func (v *State) SetValueAt(idx int, value interface{}) {
	sel := v.selectors[idx]
	sel.SetValue(sel.Upstream(v.pointer), value)
}

//ValueAt returns value for variable index
func (v *State) ValueAt(idx int) interface{} {
	sel := v.selectors[idx]
	return sel.Interface(sel.Upstream(v.pointer))
}

//Value returns value for supplied name or error
func (v *State) Value(name string) (interface{}, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return nil, err
	}
	return sel.Interface(sel.Upstream(v.pointer)), nil
}

//SetInt set int value
func (v *State) SetInt(name string, value int) error {
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
func (v *State) Index(name string) (int, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return 0, err
	}
	return int(sel.Pos), nil
}

//Int returns int or error for supplied variable name
func (v *State) Int(name string) (int, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return 0, err
	}
	return sel.Int(sel.Upstream(v.pointer)), nil
}

//SetString set value for supplied variable name or error if variable name is invalid
func (v *State) SetString(name string, value string) error {
	sel, err := v.Selector(name)
	if err != nil {
		return err
	}
	sel.SetString(sel.Upstream(v.pointer), value)
	return nil
}

//SetStringAt set value for supplied variable index
func (v *State) SetStringAt(index int, value string) error {
	sel := v.selectors[index]
	sel.SetString(sel.Upstream(v.pointer), value)
	return nil
}

//String set string value for supplied variable name or error if variable name is invalid
func (v *State) String(name string) (string, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return "", err
	}
	return sel.String(sel.Upstream(v.pointer)), nil
}

//SetBool set boolean value for supplied variable name or error if variable name is invalid
func (v *State) SetBool(name string, value bool) error {
	sel, err := v.Selector(name)
	if err != nil {
		return err
	}
	sel.SetBool(sel.Upstream(v.pointer), value)
	return nil
}

//Bool returns boolean value or error for supplied variable name
func (v *State) Bool(name string) (bool, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return false, err
	}
	return sel.Bool(sel.Upstream(v.pointer)), nil
}

//SetBoolAt set boolean value at selor position
func (v *State) SetBoolAt(idx int, value bool) {
	sel := v.selectors[idx]
	sel.SetBool(sel.Upstream(v.pointer), value)
}

//SetFloat64 set float64 value for supplied variable name or error if variable name is invalid
func (v *State) SetFloat64(name string, value float64) error {
	sel, err := v.Selector(name)
	if err != nil {
		return err
	}
	sel.SetFloat64(sel.Upstream(v.pointer), value)
	return nil
}

//Float64 returns float64 value or error for supplied variable name
func (v *State) Float64(name string) (float64, error) {
	sel, err := v.Selector(name)
	if err != nil {
		return 0.0, err
	}
	return sel.Float64(sel.Upstream(v.pointer)), nil
}

//Selector returns selector for a supploed name or error
func (v *State) Selector(name string) (*Selector, error) {
	idx, ok := v.index[name]
	if !ok {
		return nil, fmt.Errorf("undefined %v", name)
	}
	return v.selectors[idx], nil
}

func (v *State) Tracker() *Tracker {
	return v.tracker
}

//StateNew creates variables
func StateNew(dest reflect.Type, selectors []*Selector, tracker *Tracker) New {
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

	return func() *State {
		value := reflect.New(dest).Interface()
		ret := &State{
			pointer:   xunsafe.AsPointer(value),
			value:     value,
			index:     index,
			selectors: selectors,
		}
		if tracker != nil {
			ret.tracker = tracker.Clone()
			*(**Tracker)(unsafe.Pointer(uintptr(ret.pointer) +  8)) = ret.tracker
		}
		if trueSelPos != -1 {
			ret.SetBoolAt(trueSelPos, true)
		}
		return ret
	}
}
