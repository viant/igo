package et

import (
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/exec"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

//NewMapper creates a mapper
func NewMapper(x *exec.Selector, params []*exec.Selector, results []*exec.Selector, body New) (New, reflect.Type, error) {
	destType := results[0].Type
	destSliceType := reflect.SliceOf(destType)
	return func(control *Control) (exec.Compute, error) {
		var err error
		result := &mapper{x: x, srcSlice: xunsafe.NewSlice(x.Type), destSlice: xunsafe.NewSlice(destSliceType)}
		if result.body, err = body(control); err != nil {
			return nil, err
		}
		result.value = params[0]
		if len(params) > 1 {
			result.indexOffset = params[1].Offset()
			result.hasIndex = true
		}
		result.xElemType = xunsafe.NewType(x.Type.Elem())
		result.derefElem = params[0].Type.Kind() != reflect.Ptr
		result.isComponentPtr = destType.Kind() == reflect.Ptr
		return result.compute, nil
	}, destSliceType, nil
}

type mapper struct {
	srcSlice       *xunsafe.Slice
	destSlice      *xunsafe.Slice
	body           exec.Compute
	isComponentPtr bool
	xElemType      *xunsafe.Type
	derefElem      bool
	x              *exec.Selector
	value          *exec.Selector
	hasIndex       bool
	indexOffset    uintptr
}

func (s *mapper) compute(ptr unsafe.Pointer) unsafe.Pointer {
	slicePtr := s.x.Addr(ptr)
	sliceLen := s.srcSlice.Len(slicePtr)
	destSlice := reflect.MakeSlice(s.destSlice.Type, sliceLen, sliceLen)
	destSlicePtr := xunsafe.ValuePointer(&destSlice)
	for i := 0; i < sliceLen; i++ {
		item := s.srcSlice.ValuePointerAt(slicePtr, i)
		if s.derefElem {
			item = s.xElemType.Deref(item)
		}
		s.value.SetValue(s.value.Upstream(ptr), item)
		if s.hasIndex {
			*(*int)(unsafe.Pointer(uintptr(ptr) + s.indexOffset)) = i
		}
		mappedPtr := s.body(ptr)
		destItem := s.destSlice.PointerAt(destSlicePtr, uintptr(i))
		if s.isComponentPtr {
			*(*unsafe.Pointer)(destItem) = *(*unsafe.Pointer)(mappedPtr)
		} else {
			xunsafe.Copy(destItem, mappedPtr, int(s.destSlice.Type.Elem().Size()))
		}
	}
	return destSlicePtr
}
