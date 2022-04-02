package et

import (
	"fmt"
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/exec"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

//NewComposite returns composite literal evaluator
func NewComposite(cType reflect.Type, operands Operands) New {
	rType := cType
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}
	switch rType.Kind() {
	case reflect.Struct:
		return func(control *Control) (exec.Compute, error) {
			var fields = make([]*xunsafe.Field, len(operands))
			var baseTypes = make([]bool, len(fields))
			for i, op := range operands {
				fields[i] = xunsafe.FieldByName(rType, op.Key)
				if fields[i] == nil {
					return nil, fmt.Errorf("udefined %v.%v", rType.Name(), op.Key)
				}
				baseTypes[i] = isBaseType(fields[i].Kind())
			}
			operands, err := operands.operands(control)
			if err != nil {
				return nil, err
			}
			composite := newCompositeStruct(operands, fields, cType.Kind() == reflect.Ptr, cType, baseTypes)
			return composite.compute, nil
		}
	case reflect.Slice:
		return func(control *Control) (exec.Compute, error) {
			operands, err := operands.operands(control)
			if err != nil {
				return nil, err
			}
			composite := newCompositeSlice(cType, operands)
			return composite.compute, nil
		}
	}

	return func(exec *Control) (exec.Compute, error) {
		return nil, fmt.Errorf("composite not yet supported")
	}
}

type compositeStruct struct {
	operands  []*exec.Operand
	fields    []*xunsafe.Field
	baseTypes []bool
	isPtr     bool
	*xunsafe.Type
}

func (c *compositeStruct) compute(ptr unsafe.Pointer) unsafe.Pointer {
	var reflectPtr reflect.Value
	if c.isPtr {
		reflectPtr = reflect.New(c.Type.Type().Elem())
	} else {
		reflectPtr = reflect.New(c.Type.Type())
	}
	structPtr := xunsafe.ValuePointer(&reflectPtr)

	for i, op := range c.operands {
		value := op.Compute(ptr)
		if c.baseTypes[i] {
			c.fields[i].Set(structPtr, op.Interface(value))
			continue
		}
		c.fields[i].SetValue(structPtr, op.Interface(value))
	}
	if c.isPtr {
		return xunsafe.RefPointer(structPtr)
	}
	return structPtr
}

func newCompositeStruct(operands []*exec.Operand, fields []*xunsafe.Field, isPtr bool, cType reflect.Type, baseTypes []bool) *compositeStruct {
	return &compositeStruct{
		operands:  operands,
		fields:    fields,
		isPtr:     isPtr,
		Type:      xunsafe.NewType(cType),
		baseTypes: baseTypes,
	}
}

func isBaseType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.String, reflect.Bool, reflect.Float64:
		return true
	}
	return false
}

type compositeSlice struct {
	operands            []*exec.Operand
	isBaseComponentType bool
	isComponentPtr      bool
	isPtr               bool
	compType            reflect.Type
	*xunsafe.Type
}

func newCompositeSlice(sType reflect.Type, operands []*exec.Operand) *compositeSlice {
	result := &compositeSlice{operands: operands, Type: xunsafe.NewType(sType)}
	rType := sType
	if sType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		result.isPtr = true
		result.compType = rType.Elem()
	} else {
		result.compType = rType.Elem()
	}
	result.isComponentPtr = result.compType.Kind() == reflect.Ptr
	result.isBaseComponentType = isBaseType(result.compType.Kind())
	return result
}

func (c *compositeSlice) compute(ptr unsafe.Pointer) unsafe.Pointer {
	sliceLen := len(c.operands)
	var slice reflect.Value
	sliceType := c.Type.Type()
	var xSlice *xunsafe.Slice
	if c.isPtr {
		slice = reflect.MakeSlice(sliceType.Elem(), sliceLen, sliceLen)
		xSlice = xunsafe.NewSlice(sliceType.Elem())
	} else {
		slice = reflect.MakeSlice(sliceType, sliceLen, sliceLen)
		xSlice = xunsafe.NewSlice(sliceType.Elem())
	}
	slicePtr := xunsafe.ValuePointer(&slice)
	for i, op := range c.operands {
		value := op.Compute(ptr)
		itemPtr := xSlice.PointerAt(slicePtr, uintptr(i))
		if c.isComponentPtr {
			*(*unsafe.Pointer)(itemPtr) = value
			continue
		}
		xunsafe.Copy(itemPtr, value, int(c.compType.Size()))
	}
	if c.isPtr {
		return xunsafe.RefPointer(slicePtr)
	}
	return slicePtr
}
