package plan

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIiiCaller_Call(t *testing.T) {
	var testCases = []struct {
		description string
		fn          interface{}
	}{
		{
			description: "func(int,int) int",
			fn: func(x, y int) int {
				if x > y {
					return x
				}
				return y
			},
		},
		{
			description: "func(string,...interface{}) ",
			fn:          fmt.Printf,
		},
	}

	for _, testCase := range testCases {
		caller := asCaller(testCase.fn)
		assert.NotNil(t, caller, testCase.description)
	}

}
