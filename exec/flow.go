package exec

import "unsafe"

//Flow controls execution flow
type Flow uint64

const (
	//RtReturn return flag
	RtReturn = Flow(0x1)
	//RtBreak break flag
	RtBreak = Flow(0x2)
	//RtContinue continue flag
	RtContinue   = Flow(0x4)
	rtBlockBreak = RtReturn | RtBreak | RtContinue
)

//HasReturn returns true if flow has return flag set
func (f Flow) HasReturn() bool {
	return f&RtReturn != 0
}

//HasBreak returns true if flow has break falg set
func (f Flow) HasBreak() bool {
	return f&RtBreak != 0
}

//HasContinue returns true if flow has continue flag set
func (f Flow) HasContinue() bool {
	return f&RtContinue != 0
}

//HasBlockControl return true if break, continue or return flag set
func (f Flow) HasBlockControl() bool {
	return f&rtBlockBreak != 0
}

//HasControl returns trye is any flag are set
func (f Flow) HasControl() bool {
	return f != 0
}

//Reset resets flow control
func (f *Flow) Reset() {
	*f = 0
}

//AsFlow returns as flow
func AsFlow(ptr unsafe.Pointer) *Flow {
	return (*Flow)(ptr)
}
