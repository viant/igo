package plan

import (
	"github.com/viant/igo/internal/exec/et"
	"github.com/viant/igo/exec"
	"reflect"
	"strconv"
)

//Scope represents compilation scope
type Scope struct {
	*et.Control
	execType  interface{}
	count     *int
	upstream  []string
	prefix    string
	selectors *[]*exec.Selector
	index     map[string]uint16
	types     map[string]reflect.Type
	funcs     map[string]interface{}
	out       *[]*exec.Selector
	in        *[]*exec.Selector
	transients *int
	mem       *memType
}


func (s *Scope) subScope() *Scope {
	*s.count++
	var types = make(map[string]reflect.Type)
	for k, v := range s.types {
		types[k] = v
	}
	return &Scope{
		Control:   s.Control,
		upstream:  append(s.upstream, s.prefix),
		prefix:    "s" + strconv.Itoa(*s.count),
		selectors: s.selectors,
		index:     s.index,
		funcs:     s.funcs,
		types:     types,
		count:     s.count,
		mem:       s.mem,
		in:        s.in,
		out:       s.out,
		transients:s.transients,
	}
}



//NewScope creates compilation scope
func NewScope() *Scope {
	mem := newMemType()
	mem.addField("flow", reflect.TypeOf(uint64(0)))
	var selectors = make([]*exec.Selector, 0, 3)
	count := 0
	transients := 0
	var results []*exec.Selector
	var params []*exec.Selector
	control := et.Control{}
	s := &Scope{
		Control:   &control,
		prefix:    "",
		count:     &count,
		transients: &transients,
		mem:       mem,
		selectors: &selectors,
		index:     map[string]uint16{},
		types:     map[string]reflect.Type{},
		funcs:     map[string]interface{}{},
		out:       &results,
		in:        &params,
	}
	_, _ = s.DefineVariable("true", reflect.TypeOf(true))
	return s
}
