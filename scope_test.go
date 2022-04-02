package igo_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/igo"
	"testing"
)

func TestNewScope(t *testing.T) {
	{
		scope := igo.NewScope()
		fn, err := scope.Function(`func(x, y int) int {
		return x+y
	}`)
		assert.Nil(t, err)
		actual, ok := fn.(func(int, int) int)
		if !assert.True(t, ok) {
			return
		}
		assert.Equal(t, 11, actual(6, 5))
	}
}


