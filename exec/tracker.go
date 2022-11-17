package exec

//Tracker abstraction to track mutation
type (
	Tracker struct {
		init     Uint64s
		mutation Uint64s
	}
)

func (t *Tracker) Count() int {
	return t.mutation.Count()
}

//Set sets mutation for filed pos
func (t *Tracker) Set(positions []uint16) {
	if len(positions) == 0 {
		return
	}
	for _, pos := range positions {
		t.mutation.SetBit(int(pos))
	}
}

//Any returns true if mask contain any changes
func (t *Tracker) Any(mask Uint64s) bool {
	return mask.Matches(t.mutation)
}

//Has returns true if changes
func (t *Tracker) Has(pos uint16) bool {
	return t.mutation.HasBit(pos)
}

//Reset reset modification status
func (t *Tracker) Reset() {
	copy(t.mutation, t.init)
}

func (t *Tracker) Clone() *Tracker {
	var result = &Tracker{
		init:     t.init,
		mutation: make(Uint64s, len(t.init)),
	}
	return result
}

//NewTracker creates a tracker
func NewTracker(maxPos int) *Tracker {
	var result = &Tracker{
		init:     make(Uint64s, index(maxPos)+1),
		mutation: make(Uint64s, index(maxPos)+1),
	}
	return result
}
