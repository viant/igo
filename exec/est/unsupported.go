package est

import (
	"fmt"
	"github.com/viant/igo/exec"
)

//Unsupported returns unsupported fn
func Unsupported(message string) (exec.Compute, error) {
	return nil, fmt.Errorf(message)
}
