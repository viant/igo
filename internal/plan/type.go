package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"go/ast"
	"reflect"
	"time"
)

var errType = reflect.TypeOf(new(error)).Elem()
var intType = reflect.TypeOf(0)
var int64Type = reflect.TypeOf(int64(0))
var uintType = reflect.TypeOf(uint(0))
var uint64Type = reflect.TypeOf(uint64(0))
var uint8Type = reflect.TypeOf(uint8(0))
var boolType = reflect.TypeOf(true)
var stringType = reflect.TypeOf("")
var float64Type = reflect.TypeOf(0.0)
var float32Type = reflect.TypeOf(float32(0.0))
var timeType = reflect.TypeOf(time.Time{})
var anyType = reflect.TypeOf((*interface{})(nil)).Elem()

func (s *Scope) ensureType(rawType reflect.Type) {
	if rawType.Name() != "" {
		if _, ok := s.types[rawType.Name()]; !ok {
			s.types[rawType.Name()] = rawType
		}
	}
}

// EmbedType embeds supplied type
func (s *Scope) EmbedType(name string, t reflect.Type) {
	xField := s.mem.embedField(name, t)
	sel := &exec.Selector{Field: xField, ID: xField.Name}
	s.registerSelector(sel)
}

// RegisterNamedType register named type
func (s *Scope) RegisterNamedType(name string, t reflect.Type) {
	s.types[name] = t
}

// RegisterType register type
func (s *Scope) RegisterType(t reflect.Type) error {
	if t.Name() == "" {
		return fmt.Errorf("type name was empty")
	}
	if _, ok := s.types[t.Name()]; ok {
		return fmt.Errorf("type %v already defined", t.Name())
	}
	s.types[t.Name()] = t
	return nil
}

func (s *Scope) discoverType(exprType ast.Expr) (reflect.Type, error) {
	switch actual := exprType.(type) {
	case *ast.ArrayType:
		elemType, err := s.discoverType(actual.Elt)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(elemType), nil
	case *ast.Ellipsis:
		elemType, err := s.discoverType(actual.Elt)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(elemType), nil
	case *ast.MapType:
		keyType, err := s.discoverType(actual.Key)
		if err != nil {
			return nil, err
		}
		valueType, err := s.discoverType(actual.Value)
		if err != nil {
			return nil, err
		}
		return reflect.MapOf(keyType, valueType), nil
	case *ast.StarExpr:
		elemType, err := s.discoverType(actual.X)
		if err != nil {
			return nil, err
		}
		return reflect.PtrTo(elemType), nil
	case *ast.StructType:
		return s.discoverStructType(actual)
	case *ast.InterfaceType:
		// Non-empty interface construction is not representable here without a full method type system.
		// Treat interface declarations as interface{} for evaluator compatibility.
		return anyType, nil
	case *ast.FuncType:
		return s.discoverFuncType(actual)
	case *ast.ChanType:
		elemType, err := s.discoverType(actual.Value)
		if err != nil {
			return nil, err
		}
		dir := reflect.BothDir
		switch actual.Dir {
		case ast.SEND | ast.RECV:
			dir = reflect.BothDir
		case ast.RECV:
			dir = reflect.RecvDir
		case ast.SEND:
			dir = reflect.SendDir
		}
		return reflect.ChanOf(dir, elemType), nil
	}

	typeName := stringifyExpr(exprType, 0)
	litType := s.lookupType(typeName)
	if litType == nil {
		return nil, fmt.Errorf("undefined type: %s", typeName)
	}
	return litType, nil
}

func (s *Scope) discoverStructType(actual *ast.StructType) (reflect.Type, error) {
	fields := make([]reflect.StructField, 0, len(actual.Fields.List))
	for _, field := range actual.Fields.List {
		fieldType, err := s.discoverType(field.Type)
		if err != nil {
			return nil, err
		}
		if len(field.Names) == 0 {
			return nil, fmt.Errorf("embedded struct fields are not supported yet")
		}
		for _, n := range field.Names {
			fields = append(fields, reflect.StructField{
				Name: n.Name,
				Type: fieldType,
			})
		}
	}
	return reflect.StructOf(fields), nil
}

func (s *Scope) discoverFuncType(actual *ast.FuncType) (reflect.Type, error) {
	in, variadic, err := s.discoverFieldListTypes(actual.Params, true)
	if err != nil {
		return nil, err
	}
	out, _, err := s.discoverFieldListTypes(actual.Results, false)
	if err != nil {
		return nil, err
	}
	return reflect.FuncOf(in, out, variadic), nil
}

func (s *Scope) discoverFieldListTypes(list *ast.FieldList, allowVariadic bool) ([]reflect.Type, bool, error) {
	if list == nil || len(list.List) == 0 {
		return nil, false, nil
	}
	result := make([]reflect.Type, 0, len(list.List))
	var variadic bool
	for i, field := range list.List {
		fieldTypeExpr := field.Type
		if ellipsis, ok := fieldTypeExpr.(*ast.Ellipsis); ok {
			if !allowVariadic {
				return nil, false, fmt.Errorf("variadic result is not supported")
			}
			if i != len(list.List)-1 {
				return nil, false, fmt.Errorf("variadic parameter must be last")
			}
			fieldTypeExpr = &ast.ArrayType{Elt: ellipsis.Elt}
			variadic = true
		}
		fieldType, err := s.discoverType(fieldTypeExpr)
		if err != nil {
			return nil, false, err
		}
		namesCount := len(field.Names)
		if namesCount == 0 {
			namesCount = 1
		}
		for j := 0; j < namesCount; j++ {
			result = append(result, fieldType)
		}
	}
	return result, variadic, nil
}

func (s *Scope) lookupType(name string) reflect.Type {
	if result := baseType(name); result != nil {
		return result
	}
	if t, ok := s.stdlibType(name); ok {
		return t
	}
	return s.types[name]
}

func baseType(name string) reflect.Type {
	switch name {
	case "int":
		return intType
	case "int64":
		return int64Type
	case "uint":
		return uintType
	case "uint64":
		return uint64Type
	case "bool":
		return boolType
	case "string":
		return stringType
	case "float32":
		return float32Type
	case "float64":
		return float64Type
	case "time.Time":
		return timeType
	}
	return nil
}

func derefType(p reflect.Type) reflect.Type {
	if p.Kind() == reflect.Ptr {
		return derefType(p.Elem())
	}
	return p
}

func structType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Interface:
		return structType(t.Elem())
	case reflect.Struct:
		return t
	case reflect.Ptr:
		return structType(t.Elem())
	}
	return nil
}
