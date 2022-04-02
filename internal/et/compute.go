package et

import (
	"github.com/viant/igo/internal"
)

//New function creating a compute function
type New func(control *Control) (internal.Compute, error)
