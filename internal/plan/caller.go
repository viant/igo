package plan

import (
	"github.com/viant/igo/exec"
	_ "github.com/viant/igo/internal/signature"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

func toReflectValue(v interface{}, target reflect.Type) reflect.Value {
	if v == nil {
		return reflect.Zero(target)
	}
	value := reflect.ValueOf(v)
	if value.Type() == target {
		return value
	}
	if value.Type().AssignableTo(target) {
		return value
	}
	if value.Type().ConvertibleTo(target) {
		return value.Convert(target)
	}
	return value
}

type fixedReflectCaller struct {
	fn           reflect.Value
	fnType       reflect.Type
	signatureLen int
	params       []reflect.Value
	ret2         [2]unsafe.Pointer
	ret3         [3]unsafe.Pointer
	ret4         [4]unsafe.Pointer
	ret5         [5]unsafe.Pointer
}

func (c *fixedReflectCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	if len(args) != c.signatureLen {
		panic("invalid number of function arguments")
	}
	for i := 0; i < c.signatureLen; i++ {
		vPtr := args[i].Compute(ptr)
		v := args[i].Interface(vPtr)
		c.params[i] = toReflectValue(v, c.fnType.In(i))
	}
	out := c.fn.Call(c.params)
	return packResults(out, &c.ret2, &c.ret3, &c.ret4, &c.ret5)
}

type variadicReflectCaller struct {
	fn            reflect.Value
	fnType        reflect.Type
	signatureLen  int
	sliceArgIdx   int
	params        []reflect.Value
	varElemType   reflect.Type
	varSliceType  reflect.Type
	varInterfaces []interface{}
	ret2          [2]unsafe.Pointer
	ret3          [3]unsafe.Pointer
	ret4          [4]unsafe.Pointer
	ret5          [5]unsafe.Pointer
}

func (c *variadicReflectCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	if len(args) < c.sliceArgIdx {
		panic("invalid number of function arguments")
	}

	for i := 0; i < c.sliceArgIdx; i++ {
		vPtr := args[i].Compute(ptr)
		v := args[i].Interface(vPtr)
		c.params[i] = toReflectValue(v, c.fnType.In(i))
	}

	varArgCount := len(args) - c.sliceArgIdx
	if c.varElemType.Kind() == reflect.Interface && c.varElemType.NumMethod() == 0 {
		if cap(c.varInterfaces) < varArgCount {
			c.varInterfaces = make([]interface{}, varArgCount)
		}
		callArgs := c.varInterfaces[:varArgCount]
		for i := 0; i < varArgCount; i++ {
			vPtr := args[c.sliceArgIdx+i].Compute(ptr)
			callArgs[i] = args[c.sliceArgIdx+i].Interface(vPtr)
		}
		c.params[c.sliceArgIdx] = reflect.ValueOf(callArgs)
	} else {
		callArgs := reflect.MakeSlice(c.varSliceType, varArgCount, varArgCount)
		for i := 0; i < varArgCount; i++ {
			vPtr := args[c.sliceArgIdx+i].Compute(ptr)
			v := args[c.sliceArgIdx+i].Interface(vPtr)
			callArgs.Index(i).Set(toReflectValue(v, c.varElemType))
		}
		c.params[c.sliceArgIdx] = callArgs
	}

	out := c.fn.Call(c.params)
	return packResults(out, &c.ret2, &c.ret3, &c.ret4, &c.ret5)
}

func packResults(out []reflect.Value, ret2 *[2]unsafe.Pointer, ret3 *[3]unsafe.Pointer, ret4 *[4]unsafe.Pointer, ret5 *[5]unsafe.Pointer) unsafe.Pointer {
	switch len(out) {
	case 0:
		return nil
	case 1:
		return xunsafe.ValuePointer(&out[0])
	case 2:
		ret2[0] = xunsafe.ValuePointer(&out[0])
		ret2[1] = xunsafe.ValuePointer(&out[1])
		return unsafe.Pointer(ret2)
	case 3:
		ret3[0] = xunsafe.ValuePointer(&out[0])
		ret3[1] = xunsafe.ValuePointer(&out[1])
		ret3[2] = xunsafe.ValuePointer(&out[2])
		return unsafe.Pointer(ret3)
	case 4:
		ret4[0] = xunsafe.ValuePointer(&out[0])
		ret4[1] = xunsafe.ValuePointer(&out[1])
		ret4[2] = xunsafe.ValuePointer(&out[2])
		ret4[3] = xunsafe.ValuePointer(&out[3])
		return unsafe.Pointer(ret4)
	case 5:
		ret5[0] = xunsafe.ValuePointer(&out[0])
		ret5[1] = xunsafe.ValuePointer(&out[1])
		ret5[2] = xunsafe.ValuePointer(&out[2])
		ret5[3] = xunsafe.ValuePointer(&out[3])
		ret5[4] = xunsafe.ValuePointer(&out[4])
		return unsafe.Pointer(ret5)
	}
	panic("too many return parameters")
}

type fnVError func() error

func (f fnVError) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	err := f()
	return unsafe.Pointer(&err)
}

func (f fnVError) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) { _ = f() }

type fnSError func(string) error

func (f fnSError) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	s := *(*string)(args[0].Compute(ptr))
	err := f(s)
	return unsafe.Pointer(&err)
}

func (f fnSError) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	s := *(*string)(args[0].Compute(ptr))
	_ = f(s)
}

type fnVInterface func(interface{}) interface{}

func (f fnVInterface) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	v := args[0].Interface(args[0].Compute(ptr))
	result := f(v)
	return unsafe.Pointer(&result)
}

func (f fnVInterface) CallInterface(ptr unsafe.Pointer, args []*exec.Operand) interface{} {
	v := args[0].Interface(args[0].Compute(ptr))
	return f(v)
}

func (f fnVInterface) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	v := args[0].Interface(args[0].Compute(ptr))
	f(v)
}

type fnSString func(string) string

type fnSStringCaller struct {
	fn     func(string) string
	result string
}

func (f *fnSStringCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	s := *(*string)(args[0].Compute(ptr))
	f.result = f.fn(s)
	return unsafe.Pointer(&f.result)
}

func (f *fnSStringCaller) CallVoid(ptr unsafe.Pointer, args []*exec.Operand) {
	s := *(*string)(args[0].Compute(ptr))
	f.fn(s)
}

type fnSSBoolCaller struct {
	fn     func(string, string) bool
	result bool
}

func (f *fnSSBoolCaller) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	s1 := *(*string)(args[0].Compute(ptr))
	s2 := *(*string)(args[1].Compute(ptr))
	f.result = f.fn(s1, s2)
	return unsafe.Pointer(&f.result)
}

func newCaller(fn interface{}) exec.Caller {
	switch actual := fn.(type) {
	case func() error:
		return fnVError(actual)
	case func(string) error:
		return fnSError(actual)
	case func(interface{}) interface{}:
		return fnVInterface(actual)
	case func(string) string:
		return &fnSStringCaller{fn: actual}
	case func(string, string) bool:
		return &fnSSBoolCaller{fn: actual}
	}

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()
	if !fnType.IsVariadic() {
		return &fixedReflectCaller{
			fn:           fnValue,
			fnType:       fnType,
			signatureLen: fnType.NumIn(),
			params:       make([]reflect.Value, fnType.NumIn()),
		}
	}
	sliceArgIdx := fnType.NumIn() - 1
	varSliceType := fnType.In(sliceArgIdx)
	varElemType := varSliceType.Elem()
	result := &variadicReflectCaller{
		fn:           fnValue,
		fnType:       fnType,
		signatureLen: fnType.NumIn(),
		sliceArgIdx:  sliceArgIdx,
		params:       make([]reflect.Value, fnType.NumIn()),
		varElemType:  varElemType,
		varSliceType: varSliceType,
	}
	if varElemType.Kind() == reflect.Interface && varElemType.NumMethod() == 0 {
		result.varInterfaces = make([]interface{}, 0, 8)
	}
	return result
}
