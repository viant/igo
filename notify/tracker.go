package notify

import "reflect"

// Tracker represents a tracker option
type Tracker struct {
	Name   string
	Target reflect.Type
	Embed  bool
}

// NewTracker returns a tracker
func NewTracker(name string, p reflect.Type, embed bool) *Tracker {
	return &Tracker{Name: name, Target: p, Embed: embed}
}
