package plan

import (
	"github.com/viant/igo/exec"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type appendOneCaller struct {
	sliceType reflect.Type
	slice     *xunsafe.Slice
	elemSize  int
}

func newAppendOneCaller(sliceType reflect.Type) *appendOneCaller {
	return &appendOneCaller{
		sliceType: sliceType,
		slice:     xunsafe.NewSlice(sliceType),
		elemSize:  int(sliceType.Elem().Size()),
	}
}

func (c *appendOneCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	if len(args) != 2 {
		panic("appendOneCaller expects 2 arguments")
	}
	slicePtr := args[0].Compute(ptr)
	header := (*reflect.SliceHeader)(slicePtr)
	if header.Cap < header.Len+1 {
		oldLen := header.Len
		newCap := oldLen + 1
		if header.Cap > 0 {
			newCap = header.Cap * 2
		}
		newSlice := reflect.MakeSlice(c.sliceType, newCap, newCap)
		if oldLen > 0 {
			oldSlice := reflect.NewAt(c.sliceType, slicePtr).Elem()
			reflect.Copy(newSlice, oldSlice)
		}
		newHeader := (*reflect.SliceHeader)(xunsafe.ValuePointer(&newSlice))
		header.Data = newHeader.Data
		header.Cap = newHeader.Cap
		header.Len = oldLen
	}
	dest := c.slice.PointerAt(slicePtr, uintptr(header.Len))
	src := args[1].Compute(ptr)
	xunsafe.Copy(dest, src, c.elemSize)
	header.Len++
	return slicePtr
}
