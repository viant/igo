package plan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScope_Compile_UnsupportedGoStmt(t *testing.T) {
	scope := NewScope()
	_, err := scope.Compile(`go println("x")`)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "go statement is not supported")
	}
}

func TestScope_Compile_UnsupportedSelectStmt(t *testing.T) {
	scope := NewScope()
	_, err := scope.Compile(`
		ch := make(chan int, 1)
		select {
		case v := <-ch:
			_ = v
		default:
		}
	`)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "select statement is not supported")
	}
}

func TestScope_Compile_UnsupportedSendStmt(t *testing.T) {
	scope := NewScope()
	_, err := scope.Compile(`
		ch := make(chan int, 1)
		ch <- 1
	`)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "channel send is not supported")
	}
}
