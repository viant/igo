package et

import (
	"fmt"
	"github.com/viant/igo/internal/exec"
)

//Unsupported returns unsupported fn
func Unsupported(message string) (exec.Compute, error) {
	return nil, fmt.Errorf(message)
}
