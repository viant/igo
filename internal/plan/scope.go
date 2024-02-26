package plan

import (
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal/et"
	"github.com/viant/igo/metric"
	"github.com/viant/igo/notify"
	"github.com/viant/igo/option"
	"reflect"
	"strconv"
)

// Scope represents compilation scope
type Scope struct {
	*et.Control
	options    option.Options
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
	trackLen   int
	trackRoot  string
}

func (s *Scope) SubScope(options ...option.Option) *Scope {
	*s.count++
	var types = make(map[string]reflect.Type)
	for k, v := range s.types {
		types[k] = v
	}
	result := &Scope{
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
		Metric:     s.Metric,
	}
	result.options.Merge(&s.options, options)
	return result
}

func (s *Scope) setTracker(tracker *notify.Tracker) error {
	if tracker == nil {
		return nil
	}
	if tracker.Embed {
		if err := s.DefineEmbedVariable(tracker.Name, tracker.Target); err != nil {
			return err
		}
	} else {
		if _, err := s.DefineVariable(tracker.Name, tracker.Target); err != nil {
			return err
		}
	}
	s.trackType = tracker.Target
	s.trackRoot = tracker.Name
	return nil
}

// NewScope creates compilation scope
func NewScope(options ...option.Option) *Scope {
	mem := newMemType()
	mem.addField("_flow", reflect.TypeOf(uint64(0)))
	opts := option.NewOptions(options...)

	tracker := opts.Tracker
	if tracker != nil {
		mem.addField("_trk", reflect.TypeOf(&exec.Tracker{}))
	}

	ret := newScope(mem)
	ret.options.Merge(&ret.options, options)
	err := ret.setTracker(tracker)
	if err != nil {
		panic(err)
	}
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
	_, _ = s.DefineVariable("false", reflect.TypeOf(true))

	return s
}
