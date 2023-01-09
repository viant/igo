package exec

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"sync"
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
	pool      *sync.Pool
}

//Pointer returns variables pointer
func (s *State) Pointer() unsafe.Pointer {
	return s.pointer
}

func (s *State) Release() {
	if s.pool == nil {
		return
	}
	s.pool.Put(s)

}

func (s *State) Interface() interface{} {
	return s.value
}

//SetValue sets value for
func (s *State) SetValue(name string, value interface{}) error {
	sel, err := s.Selector(name)
	if err != nil {
		return err
	}
	sel.SetValue(sel.Upstream(s.pointer), value)
	return nil
}

//SetValueAt set values at position
func (s *State) SetValueAt(idx int, value interface{}) {
	sel := s.selectors[idx]
	sel.SetValue(sel.Upstream(s.pointer), value)
}

//ValueAt returns value for variable index
func (s *State) ValueAt(idx int) interface{} {
	sel := s.selectors[idx]
	return sel.Interface(sel.Upstream(s.pointer))
}

//Value returns value for supplied name or error
func (s *State) Value(name string) (interface{}, error) {
	sel, err := s.Selector(name)
	if err != nil {
		return nil, err
	}
	return sel.Interface(sel.Upstream(s.pointer)), nil
}

//SetInt set int value
func (s *State) SetInt(name string, value int) error {
	sel, err := s.Selector(name)
	if err != nil {
		return err
	}
	if sel.Pathway == PathwayDirect {
		sel.SetInt(s.pointer, value)
		return nil
	}
	sel.SetInt(sel.Upstream(s.pointer), value)
	return nil
}

//Index returns variable index for specified name
func (s *State) Index(name string) (int, error) {
	sel, err := s.Selector(name)
	if err != nil {
		return 0, err
	}
	return int(sel.Pos), nil
}

//Int returns int or error for supplied variable name
func (s *State) Int(name string) (int, error) {
	sel, err := s.Selector(name)
	if err != nil {
		return 0, err
	}
	return sel.Int(sel.Upstream(s.pointer)), nil
}

//SetString set value for supplied variable name or error if variable name is invalid
func (s *State) SetString(name string, value string) error {
	sel, err := s.Selector(name)
	if err != nil {
		return err
	}
	sel.SetString(sel.Upstream(s.pointer), value)
	return nil
}

//SetStringAt set value for supplied variable index
func (s *State) SetStringAt(index int, value string) error {
	sel := s.selectors[index]
	sel.SetString(sel.Upstream(s.pointer), value)
	return nil
}

//String set string value for supplied variable name or error if variable name is invalid
func (s *State) String(name string) (string, error) {
	sel, err := s.Selector(name)
	if err != nil {
		return "", err
	}
	return sel.String(sel.Upstream(s.pointer)), nil
}

//SetBool set boolean value for supplied variable name or error if variable name is invalid
func (s *State) SetBool(name string, value bool) error {
	sel, err := s.Selector(name)
	if err != nil {
		return err
	}
	sel.SetBool(sel.Upstream(s.pointer), value)
	return nil
}

//Bool returns boolean value or error for supplied variable name
func (s *State) Bool(name string) (bool, error) {
	sel, err := s.Selector(name)
	if err != nil {
		return false, err
	}
	return sel.Bool(sel.Upstream(s.pointer)), nil
}

//SetBoolAt set boolean value at selor position
func (s *State) SetBoolAt(idx int, value bool) {
	sel := s.selectors[idx]
	sel.SetBool(sel.Upstream(s.pointer), value)
}

//SetFloat64 set float64 value for supplied variable name or error if variable name is invalid
func (s *State) SetFloat64(name string, value float64) error {
	sel, err := s.Selector(name)
	if err != nil {
		return err
	}
	sel.SetFloat64(sel.Upstream(s.pointer), value)
	return nil
}

//Float64 returns float64 value or error for supplied variable name
func (s *State) Float64(name string) (float64, error) {
	sel, err := s.Selector(name)
	if err != nil {
		return 0.0, err
	}
	return sel.Float64(sel.Upstream(s.pointer)), nil
}

//Selector returns selector for a supploed name or error
func (s *State) Selector(name string) (*Selector, error) {
	idx, ok := s.index[name]
	if !ok {
		return nil, fmt.Errorf("undefined %s", name)
	}
	return s.selectors[idx], nil
}

func (s *State) Tracker() *Tracker {
	return s.tracker
}

//StateNew creates variables
func StateNew(dest reflect.Type, selectors []*Selector, tracker *Tracker, pool *sync.Pool) New {
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
	falseSelPos := -1

	for i := range selectors {
		sel := selectors[i]
		index[sel.ID] = sel.Pos
		if sel.ID == "true" {
			trueSelPos = int(sel.Pos)
		}
		if sel.ID == "false" {
			falseSelPos = int(sel.Pos)
		}
	}
	return func() *State {
		value := reflect.New(dest).Interface()
		ret := &State{
			pool:      pool,
			pointer:   xunsafe.AsPointer(value),
			value:     value,
			index:     index,
			selectors: selectors,
		}
		if tracker != nil {
			ret.tracker = tracker.Clone()
			*(**Tracker)(unsafe.Pointer(uintptr(ret.pointer) + 8)) = ret.tracker
		}
		if trueSelPos != -1 {
			ret.SetBoolAt(trueSelPos, true)
		}
		if falseSelPos != -1 {
			ret.SetBoolAt(falseSelPos, false)
		}
		return ret
	}
}
