package plan

import (
	"github.com/viant/xunsafe"
	"reflect"
	"strings"
)

var defaultPkg = reflect.TypeOf(Scope{}).PkgPath()

//memType represents
type memType struct {
	reflect.Type
	fields []reflect.StructField
}

func (t *memType) embedField(name string, fType reflect.Type) *xunsafe.Field {
	pkg := ""
	idx := len(t.fields)
	t.fields = append(t.fields, reflect.StructField{Name: strings.Title(name), Anonymous: true, Type: fType, PkgPath: pkg})
	t.Type = reflect.StructOf(t.fields)
	result := xunsafe.NewField(t.Type.Field(idx))
	return result
}

func (t *memType) updateField(name string, fType reflect.Type) {
	for i, f := range t.fields {
		if f.Name == name {
			t.fields[i].Type = fType
			break
		}
	}
	t.Type = reflect.StructOf(t.fields)
}

func (t *memType) addField(name string, fType reflect.Type) *xunsafe.Field {
	pkg := ""
	if name[0] > 'Z' {
		pkg = defaultPkg
	}
	idx := len(t.fields)
	t.fields = append(t.fields, reflect.StructField{Name: name, Type: fType, PkgPath: pkg})
	t.Type = reflect.StructOf(t.fields)
	result := xunsafe.NewField(t.Type.Field(idx))
	return result
}

func newMemType() *memType {
	return &memType{Type: reflect.StructOf([]reflect.StructField{})}
}
