package est

import (
	"github.com/viant/igo/exec"
)

//Control controls tree generation based on flow features
type Control struct {
	exec.Flow
}

//Concat concat
func (c *Control) Concat(flag exec.Flow) {
	c.Flow |= flag
}
