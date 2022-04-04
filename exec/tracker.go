package exec

import (
	"reflect"
)


//Tracker abstraction to track mutation
type Tracker struct {
	init     []uint64
	Mutation []uint64
	Nested   []*Tracker
}


//Set sets mutation for filed pos
func (t *Tracker) Set(pos []uint16) {
	Uint64s(t.Mutation).SetBit(int(pos[0]))
	if len(pos) > 1 {
		t.Nested[pos[0]].Set(pos[1:])
	}
}


//Changed returns true if changes
func (t *Tracker) Changed(pos ...uint16) bool {
	if len(pos) > 1 {
		return t.Changed(pos[1:]...)
	}
	return Uint64s(t.Mutation).HasBit(int(pos[0]))
}


//Reset reset modification status
func (t *Tracker) Reset() {
	copy(t.Mutation, t.init)
	if len(t.Nested) == 0 {
		return
	}
	for i := range t.Nested {
		if t.Nested[i] == nil {
			continue
		}
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


type Uint64s []uint64

//HasBit returns true if a bit at position in set
func (o Uint64s) HasBit(pos int) bool {
	return o[index(pos)] & (1 << pos % 64) != 0
}

//ClearBit clears bit at position in set
func (o Uint64s) ClearBit(pos int) {
	o[index(pos)] &= ^(1 << (pos % 64))
}

//SetBit sets bit at position in set
func (o Uint64s) SetBit(pos int) {
	o[index(pos)] |= 1 << (pos % 64)
}

func index(pos int) int {
	return pos / 64
}



