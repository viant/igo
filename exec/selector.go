package exec

import (
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

// Selector represent data selector
type Selector struct {
	ID string
	*xunsafe.Field
	Pathway
	Slice       *xunsafe.Slice
	Index       *Operand
	IsErrorType bool
	//and Ancestors lookup need to take place to resolve data address
	Ancestors []*Selector
	Pos       uint16
	xPos      []uint16
}

// XPos returns selector position including with ancestor
func (s *Selector) XPos() []uint16 {
	if len(s.xPos) > 0 {
		return s.xPos
	}
	if len(s.Ancestors) > 0 {
		for _, item := range s.Ancestors {
			s.xPos = append(s.xPos, item.Pos)
		}
	}
	s.xPos = append(s.xPos, s.Pos)
	return s.xPos
}

// IndexPointer returns slice item pointer
func (s *Selector) IndexPointer(ptr unsafe.Pointer, index int) unsafe.Pointer {
	slicePtr := s.Upstream(ptr)
	return s.Slice.PointerAt(slicePtr, uintptr(index))
}

// Upstream returns Ancestors pointer, excluding current Selector ogf
func (s *Selector) Upstream(ptr unsafe.Pointer) unsafe.Pointer {
	transientPtr := ptr
	i := 0
	l := len(s.Ancestors)
begin:
	if i >= l {
		return transientPtr
	}
	sel := s.Ancestors[i]
	if idxSel := sel.Index; idxSel != nil {
		idx := xunsafe.AsInt(idxSel.Compute(ptr))
		transientPtr = sel.Slice.PointerAt(transientPtr, uintptr(idx))
		if sel.Kind() == reflect.Ptr {
			transientPtr = xunsafe.DerefPointer(transientPtr)
		}
	} else {
		transientPtr = sel.ValuePointer(transientPtr)
	}
	i++
	goto begin
}

// Addr returns Selector address
func (s *Selector) Addr(pointer unsafe.Pointer) unsafe.Pointer {
	var ret unsafe.Pointer
	if s.Pathway.IsDirect() {
		if len(s.Ancestors) == 0 {
			return s.Field.Pointer(pointer)
		}
		return s.Field.Pointer(s.Upstream(pointer))
	} else if idxSel := s.Index; idxSel != nil {
		idx := xunsafe.AsInt(idxSel.Compute(pointer))
		ptr := s.Upstream(pointer)
		ret = s.Slice.PointerAt(ptr, uintptr(idx))
	} else {
		ret = unsafe.Pointer(uintptr(s.Upstream(pointer)) + s.Field.Offset)
		//mey break some func ,
		//if s.Kind() == reflect.Ptr && s.Elem().Kind() == reflect.Struct {
		//	ret = xunsafe.DerefPointer(ret)
		//}
	}
	return ret
}

// UpstreamOffset returns Selector Offset
func (s *Selector) UpstreamOffset() uintptr {
	if s == nil {
		return 0
	}
	result := uintptr(0)
	for _, u := range s.Ancestors {
		result += u.Field.Offset
	}
	return result
}

// Interface return an interface,
func (s *Selector) Interface(ptr unsafe.Pointer) interface{} {
	adr := s.Upstream(ptr)
	if s.IsErrorType {
		return xunsafe.AsError(adr)
	}
	x := s.Field.Interface(adr)
	return x
}

// SetInterface sets interface, it takes a state pointer
func (s *Selector) SetInterface(ptr unsafe.Pointer, value interface{}) {
	addr := s.Upstream(ptr)
	s.Field.SetValue(addr, value)
}

func (s *Selector) SetInterfaces(ptr unsafe.Pointer, value []interface{}) {
	addr := s.Upstream(ptr)
	s.Field.SetValue(addr, value)
}

// SetValue sets value, tt takes an upstream pointer
func (s *Selector) SetValue(ptr unsafe.Pointer, value interface{}) {
	if s.IsErrorType {
		err, _ := value.(error)
		fPtr := s.Field.Pointer(ptr)
		*xunsafe.AsErrorPtr(fPtr) = err
	} else {
		s.Field.SetValue(ptr, value)
	}
}

// Offset returns Selector Offset
func (s *Selector) Offset() uintptr {
	if s == nil {
		return 0
	}
	result := s.Field.Offset
	for _, u := range s.Ancestors {
		result += u.Field.Offset
	}
	return result
}
