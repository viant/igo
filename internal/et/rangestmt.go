package et

import (
	"fmt"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type rangeStmt struct {
	slice          *xunsafe.Slice
	body           internal.Compute
	isComponentPtr bool
	x              *exec.Operand
	key            *exec.Selector
	value          *exec.Selector
}

func (s *rangeStmt) computeRange(ptr unsafe.Pointer) unsafe.Pointer {
	slicePtr := s.x.Compute(ptr)
	s.slice.Range(slicePtr, func(index int, item interface{}) bool {
		fmt.Printf("%T %v\n", item, item, s.isComponentPtr)
		*(*int)(s.key.Addr(ptr)) = index
		if s.value != nil {
			itemPtr := xunsafe.AsPointer(item)
			//if s.isComponentPtr {
			//itemPtr = xunsafe.DerefPointer(itemPtr)
			//}
			*(*unsafe.Pointer)(s.value.Addr(ptr)) = itemPtr
		}
		s.body(ptr)
		return true
	})
	return nil
}

//NewRange creates a range stmt
func NewRange(x *Operand, key, value *exec.Selector, body New) (New, error) {
	return func(control *Control) (internal.Compute, error) {
		var err error
		sliceType := x.Selector.Type
		if sliceType.Kind() == reflect.Ptr {
			sliceType = sliceType.Elem()
		}
		stmt := &rangeStmt{slice: xunsafe.NewSlice(sliceType), key: key, value: value}
		stmt.isComponentPtr = sliceType.Elem().Kind() == reflect.Ptr
		if stmt.x, err = x.NewOperand(control); err != nil {
			return nil, err
		}
		if stmt.body, err = body(control); err != nil {
			return nil, err
		}
		return stmt.computeRange, nil

	}, nil
}
