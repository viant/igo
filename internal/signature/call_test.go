package signature

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/igo/exec"
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
		{
			description: "func(interface{}) int64",
			fn: func(v interface{}) int64 {
				switch actual := v.(type) {
				case int:
					return int64(actual)
				case int64:
					return actual
				}
				return 0
			},
		},
		{
			description: "func(int,int) bool",
			fn: func(x, y int) bool {
				return x > y
			},
		},
		{
			description: "func(string,string) int",
			fn: func(x, y string) int {
				if x == y {
					return 0
				}
				if x < y {
					return -1
				}
				return 1
			},
		},
		{
			description: "func(int64,int64) int64",
			fn: func(x, y int64) int64 {
				return x + y
			},
		},
		{
			description: "func(int64,int64) bool",
			fn: func(x, y int64) bool {
				return x >= y
			},
		},
	}

	for _, testCase := range testCases {
		caller := exec.AsCaller(testCase.fn)
		assert.NotNil(t, caller, testCase.description)
	}

}
