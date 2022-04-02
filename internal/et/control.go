package et

import (
	"github.com/viant/igo/internal"
)

//Control controls tree generation based on flow features
type Control struct {
	internal.Flow
}

//Concat concat
func (c *Control) Concat(flag internal.Flow) {
	c.Flow |= flag
}
