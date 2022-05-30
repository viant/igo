package plan

import (
	"fmt"
	"github.com/viant/igo/exec"
	"go/ast"
	"reflect"
	"strings"
	"time"
)

var errType = reflect.TypeOf(new(error)).Elem()
var intType = reflect.TypeOf(0)
var uint8Type = reflect.TypeOf(uint8(0))
var boolType = reflect.TypeOf(true)
var stringType = reflect.TypeOf("")
var float64Type = reflect.TypeOf(0.0)
var float32Type = reflect.TypeOf(float32(0.0))
var timeType = reflect.TypeOf(time.Time{})

func (s *Scope) ensureType(rawType reflect.Type) {
	if rawType.Name() != "" {
		if _, ok := s.types[rawType.Name()]; !ok {
			s.types[rawType.Name()] = rawType
		}
	}
}

//EmbedType embeds supplied type
func (s *Scope) EmbedType(name string, t reflect.Type) {
	xField := s.mem.embedField(name, t)
	sel := &exec.Selector{Field: xField, ID: xField.Name}
	s.regsterSelector(sel)
}

//RegisterNamedType register named type
func (s *Scope) RegisterNamedType(name string, t reflect.Type) {
	s.types[name] = t
}

//RegisterType register type
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
	isArray := isArrayIdentifier(exprType)
	typeName := stringifyExpr(exprType, 0)
	isPointer := false
	if strings.HasPrefix(typeName, "*") {
		typeName = typeName[1:]
		isPointer = true
	}
	var litType reflect.Type
	if isArray {
		elemType := s.lookupType(typeName)
		if elemType == nil {
			return nil, fmt.Errorf("undefined type: %s", typeName)
		}
		if isPointer {
			elemType = reflect.PtrTo(elemType)
		}
		return reflect.SliceOf(elemType), nil
	}
	litType = s.lookupType(typeName)
	if litType == nil {
		return nil, fmt.Errorf("undefined type: %s", typeName)
	}
	if isPointer {
		litType = reflect.PtrTo(litType)
	}
	return litType, nil
}

func (s *Scope) lookupType(name string) reflect.Type {
	if result := baseType(name); result != nil {
		return result
	}
	return s.types[name]
}

func baseType(name string) reflect.Type {
	switch name {
	case "int":
		return intType
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
