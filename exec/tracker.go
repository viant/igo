package exec

import "reflect"

//Tracker abstraction to track mutation
type Tracker struct {
	init     []uint64
	Mutation []uint64
	Nested   []*Tracker
}

//Reset reset modification status
func (t *Tracker) Reset() {
	copy(t.Mutation, t.init)
	if len(t.Nested) == 0 {
		return
	}
	for i := range t.Nested {
		t.Nested[i].Reset()
	}
}

func (t *Tracker) Clone() *Tracker {
	var result = &Tracker{
		init:     t.init,
		Mutation: make([]uint64, len(t.init)),
		Nested:   make([]*Tracker, len(t.Nested)),
	}
	for i, item := range t.Nested {
		if item == nil {
			continue
		}
		result.Nested[i] = t.Nested[i].Clone()
	}
	return result
}

//NewTracker creates a tracker
func NewTracker(target reflect.Type) *Tracker {
	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}
	if target.Kind() != reflect.Struct {
		return nil
	}
	fieldCount := target.NumField()
	var result = &Tracker{
		init:     make([]uint64, 1+fieldCount/64),
		Mutation: make([]uint64, 1+fieldCount/64),
		Nested:   make([]*Tracker, fieldCount),
	}
	for i := 0; i < fieldCount; i++ {
		filed := target.Field(i)
		fType := filed.Type
		if fType.Kind() == reflect.Ptr {
			fType = fType.Elem()
		}
		if fType.Kind() != reflect.Struct {
			continue
		}
		result.Nested[i] = NewTracker(fType)
	}
	return result
}
