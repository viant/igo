package option

import "reflect"

//Tracker represents a tracker option
type Tracker struct {
	Name string
	Target reflect.Type
}

//NewTracker returns a tracker
func NewTracker(name string, p reflect.Type) *Tracker {
	return &Tracker{Name: name, Target: p}
}

