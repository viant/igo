package plan

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal/et"
	"github.com/viant/igo/metric"
	"github.com/viant/igo/option"
	"reflect"
	"strconv"
)

//Scope represents compilation scope
type Scope struct {
	*et.Control
	Metric     metric.Stmt
	execType   interface{}
	count      *int
	upstream   []string
	prefix     string
	selectors  *[]*exec.Selector
	index      map[string]uint16
	types      map[string]reflect.Type
	funcs      map[string]interface{}
	out        *[]*exec.Selector
	in         *[]*exec.Selector
	transients *int
	mem        *memType
	trackType  reflect.Type
	trackRoot  string
}

func (s *Scope) subScope() *Scope {
	*s.count++
	var types = make(map[string]reflect.Type)
	for k, v := range s.types {
		types[k] = v
	}
	return &Scope{
		Control:    s.Control,
		upstream:   append(s.upstream, s.prefix),
		prefix:     "s" + strconv.Itoa(*s.count),
		selectors:  s.selectors,
		index:      s.index,
		funcs:      s.funcs,
		types:      types,
		count:      s.count,
		mem:        s.mem,
		in:         s.in,
		out:        s.out,
		transients: s.transients,
		trackType:  s.trackType,
	}
}

func (s *Scope) setTracker(tracker *option.Tracker) {
	if tracker == nil {
		return
	}
	_, _ = s.DefineVariable(tracker.Name, tracker.Target)
	s.trackType = tracker.Target
	s.trackRoot = tracker.Name
}

//NewScope creates compilation scope
func NewScope(options ...option.Option) *Scope {
	mem := newMemType()
	mem.addField("_flow", reflect.TypeOf(uint64(0)))
	tracker := option.Options(options).Tracker()
	if tracker != nil {
		mem.addField("_trk", reflect.TypeOf(&exec.Tracker{}))
	}
	ret := newScope(mem)
	ret.setTracker(tracker)
	return ret
}

func newScope(mem *memType) *Scope {
	var selectors = make([]*exec.Selector, 0, 3)
	count := 0
	transients := 0
	var results []*exec.Selector
	var params []*exec.Selector
	control := et.Control{}
	s := &Scope{
		Control:    &control,
		prefix:     "",
		count:      &count,
		transients: &transients,
		mem:        mem,
		selectors:  &selectors,
		index:      map[string]uint16{},
		types:      map[string]reflect.Type{},
		funcs:      map[string]interface{}{},
		out:        &results,
		in:         &params,
	}
	_, _ = s.DefineVariable("true", reflect.TypeOf(true))
	return s
}
