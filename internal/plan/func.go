package plan

import (
	"fmt"
	"github.com/viant/igo/internal/exec/et"
	"github.com/viant/xunsafe"
	"go/ast"
	"reflect"
	"strconv"
)

//RegisterFunc register func
func (s *Scope) RegisterFunc(name string, fn interface{}) {
	s.funcs[name] = fn
}

func (s *Scope) lookupFunction(id string) (interface{}, bool) {
	switch id {
	case "len":
		return length, true
	case "append":
		return sliceAppend, true
	case "float64":
		return asFlot64, true
	case "float32":
		return asFlot32, true
	case "int":
		return asInt, true
	case "int64":
		return asInt64, true
	case "string":
		return asString, true
	}
	fn, ok := s.funcs[id]
	return fn, ok
}

func (s *Scope) compileBuildIn(id string, expr *ast.CallExpr) (et.New, reflect.Type, bool, error) {
	holder, method, err := s.selectorFun(id)
	if err != nil {
		return nil, nil, false, err
	}
	switch method {
	case "Reduce":
		newFn, fType, err := s.compileReducer(holder, expr)
		return newFn, fType, true, err
	case "Map":
		newFn, fType, err := s.compileMapper(holder, expr)
		return newFn, fType, true, err
	case "make":
		newFn, fType, err := s.compileMake(expr)
		return newFn, fType, true, err
	}
	return nil, nil, false, nil
}

func (s *Scope) compileMake(expr *ast.CallExpr) (et.New, reflect.Type, error) {
	if len(expr.Args) < 2 {
		return nil, nil, fmt.Errorf("invalid make args call")
	}
	sliceType, err := s.discoverType(expr.Args[0])
	if err != nil {
		return nil, nil, err
	}
	args, err := s.assembleOperands(expr.Args[1:], false)
	if err != nil {
		return nil, nil, err
	}
	return et.NewMake(sliceType, args)
}

func asFlot64(v interface{}) float64 {
	switch actual := v.(type) {
	case float64:
		return actual
	case float32:
		return float64(actual)
	case int:
		return float64(actual)
	case uint:
		return float64(actual)
	case int64:
		return float64(actual)
	case uint64:
		return float64(actual)
	case string:
		f, err := strconv.ParseFloat(actual, 64)
		if err != nil {
			panic(err)
		}
		return f
	}
	panic(fmt.Errorf("unsupported float64 cast: %T", v))
}

func asFlot32(v interface{}) float32 {
	switch actual := v.(type) {
	case float64:
		return float32(actual)
	case float32:
		return actual
	case int:
		return float32(actual)
	case uint:
		return float32(actual)
	case int64:
		return float32(actual)
	case uint64:
		return float32(actual)
	case string:
		f, err := strconv.ParseFloat(actual, 64)
		if err != nil {
			panic(err)
		}
		return float32(f)
	}
	panic(fmt.Errorf("unsupported float32 cast: %T", v))
}

func asInt(v interface{}) int {
	switch actual := v.(type) {
	case float64:
		return int(actual)
	case float32:
		return int(actual)
	case int:
		return actual
	case uint:
		return int(actual)
	case int64:
		return int(actual)
	case uint64:
		return int(actual)
	case string:
		i, err := strconv.Atoi(actual)
		if err != nil {
			panic(err)
		}
		return i
	}
	panic(fmt.Errorf("unsupported int cast: %T", v))
}

func asInt64(v interface{}) int64 {
	switch actual := v.(type) {
	case float64:
		return int64(actual)
	case float32:
		return int64(actual)
	case int:
		return int64(actual)
	case uint:
		return int64(actual)
	case int64:
		return actual
	case uint64:
		return int64(actual)
	case string:
		i, err := strconv.Atoi(actual)
		if err != nil {
			panic(err)
		}
		return int64(i)
	}
	panic(fmt.Errorf("unsupported int64 cast: %T", v))
}

func asString(v interface{}) string {
	switch actual := v.(type) {
	case float64:
		return strconv.FormatFloat(actual, 'f', 10, 64)
	case float32:
		return strconv.FormatFloat(float64(actual), 'f', 10, 32)
	case int:
		return strconv.Itoa(actual)
	case uint:
		return strconv.Itoa(int(actual))
	case int64:
		return strconv.Itoa(int(actual))
	case uint64:
		return strconv.Itoa(int(actual))
	case string:
		return actual
	case []byte:
		return string(actual)
	}
	return fmt.Sprintf("%v", v)
}

func length(v interface{}) int {
	if s, ok := v.(string); ok {
		return len(s)
	}
	ptr := xunsafe.AsPointer(v)
	header := *(*reflect.SliceHeader)(ptr)
	return header.Len
}

func sliceAppend(source interface{}, args ...interface{}) interface{} {
	aSlice := reflect.ValueOf(source)
	switch len(args) {
	case 1:
		return reflect.Append(aSlice, reflect.ValueOf(args[0])).Interface()
	case 2:
		return reflect.Append(aSlice, reflect.ValueOf(args[0]), reflect.ValueOf(args[1])).Interface()
	case 3:
		return reflect.Append(aSlice, reflect.ValueOf(args[0]), reflect.ValueOf(args[1]), reflect.ValueOf(args[2])).Interface()
	}
	var values = make([]reflect.Value, len(args))
	for i, arg := range args {
		values[i] = reflect.ValueOf(arg)
	}
	return reflect.Append(aSlice, values...).Interface()
}
