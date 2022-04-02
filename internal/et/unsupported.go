package et

import (
	"fmt"
	"github.com/viant/igo/internal"
)

//Unsupported returns unsupported fn
func Unsupported(message string) (internal.Compute, error) {
	return nil, fmt.Errorf(message)
}
