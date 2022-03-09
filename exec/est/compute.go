package est

import "github.com/viant/igo/exec"

//New function creating a compute function
type New func(control *Control) (exec.Compute, error)
