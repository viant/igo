package et

import "github.com/viant/igo/internal/exec"

//New function creating a compute function
type New func(control *Control) (exec.Compute, error)
