package exec

import (
	"reflect"
)

var funcTypes = []reflect.Type{}
var callerTypes = []reflect.Type{}

var c *Caller
var callerType = reflect.TypeOf(c).Elem()

var f *Func
var funcType = reflect.TypeOf(f).Elem()

//RegisterSignature register caller or func signature type adapter
func RegisterSignature(p reflect.Type) {
	if p.AssignableTo(callerType) {
		callerTypes = append(callerTypes, p)
	}
	if p.AssignableTo(funcType) {
		funcTypes = append(funcTypes, p)
	}

}

//AsFunc returns func or nil
func AsFunc(fn interface{}) Func {
	if aFunc, ok := fn.(Func); ok {
		return aFunc
	}
	fnValue := reflect.ValueOf(fn)
	for _, candidate := range funcTypes {
		if fnValue.CanConvert(candidate) {
			res := fnValue.Convert(candidate).Interface()
			return res.(Func)
		}
	}
	return nil
}

//AsCaller returns a caller or nil
func AsCaller(fn interface{}) Caller {
	if caller, ok := fn.(Caller); ok {
		return caller
	}
	fnValue := reflect.ValueOf(fn)
	for _, candidate := range callerTypes {
		if fnValue.CanConvert(candidate) {
			res := fnValue.Convert(candidate).Interface()
			return res.(Caller)
		}
	}
	return nil
}
