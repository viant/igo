package internal

import "unsafe"

// Flow controls execution flow
type Flow uint64

const (
	//RtReturn return flag
	RtReturn = Flow(0x1)
	//RtBreak break flag
	RtBreak = Flow(0x2)
	//RtContinue continue flag
	RtContinue = Flow(0x4)
	// RtGoto goto flag.
	RtGoto = Flow(0x8)
	// RtFallthrough fallthrough flag.
	RtFallthrough  = Flow(0x10)
	flowLabelShift = 32
	flowLabelMask  = Flow(0xFFFFFFFF00000000)
	flowFlagMask   = ^flowLabelMask
	rtBlockBreak   = RtReturn | RtBreak | RtContinue | RtGoto | RtFallthrough
)

// HasReturn returns true if flow has return flag set
func (f Flow) HasReturn() bool {
	return f&RtReturn != 0
}

// HasBreak returns true if flow has break falg set
func (f Flow) HasBreak() bool {
	return f&RtBreak != 0
}

// HasContinue returns true if flow has continue flag set
func (f Flow) HasContinue() bool {
	return f&RtContinue != 0
}

// HasGoto returns true if flow has goto flag set.
func (f Flow) HasGoto() bool {
	return f&RtGoto != 0
}

// HasFallthrough returns true if flow has fallthrough flag set.
func (f Flow) HasFallthrough() bool {
	return f&RtFallthrough != 0
}

// Label returns control-flow label id.
func (f Flow) Label() uint32 {
	return uint32((uint64(f) & uint64(flowLabelMask)) >> flowLabelShift)
}

// HasBlockControl return true if break, continue or return flag set
func (f Flow) HasBlockControl() bool {
	return f&rtBlockBreak != 0
}

// HasControl returns trye is any flag are set
func (f Flow) HasControl() bool {
	return f != 0
}

// Reset resets flow control
func (f *Flow) Reset() {
	*f = 0
}

// SetBranch sets branch flag and optional label id.
func (f *Flow) SetBranch(flag Flow, labelID uint32) {
	*f = (*f & flowFlagMask) | flag | (Flow(labelID) << flowLabelShift)
}

// ClearBreak clears break flag.
func (f *Flow) ClearBreak() {
	f.clearFlag(RtBreak)
}

// ClearContinue clears continue flag.
func (f *Flow) ClearContinue() {
	f.clearFlag(RtContinue)
}

// ClearGoto clears goto flag.
func (f *Flow) ClearGoto() {
	f.clearFlag(RtGoto)
}

// ClearFallthrough clears fallthrough flag.
func (f *Flow) ClearFallthrough() {
	f.clearFlag(RtFallthrough)
}

func (f *Flow) clearFlag(flag Flow) {
	*f &^= flag
	if *f&flowFlagMask == 0 {
		*f &^= flowLabelMask
	}
}

// AsFlow returns as flow
func AsFlow(ptr unsafe.Pointer) *Flow {
	return (*Flow)(ptr)
}
