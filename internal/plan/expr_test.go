package plan

import (
	"github.com/stretchr/testify/assert"
	"io"
	"reflect"
	"testing"
)

func TestScope_IntExpression(t *testing.T) {
	type Foo struct {
		X int
		Y int
	}
	var testCase = []struct {
		description string
		expr        string
		state       map[string]interface{}
		expect      int
	}{
		{
			description: "basic addition",
			expr:        "x+y",
			state:       map[string]interface{}{"x": 10, "y": 20},
			expect:      30,
		},

		{
			description: "basic subtraction",
			expr:        "x-y",
			state:       map[string]interface{}{"x": 10, "y": 20},
			expect:      -10,
		},
		{
			description: "simple expr +-",
			expr:        "x+y-z",
			state:       map[string]interface{}{"x": 10, "y": 20, "z": 5},
			expect:      25,
		},
		{
			description: "simple expr */ literal",
			expr:        "5*x/y*z",
			state:       map[string]interface{}{"x": 10, "y": 25, "z": 5},
			expect:      10,
		},
		{
			description: "simple expr */ literal",
			expr:        "5*x/y*z",
			state:       map[string]interface{}{"x": 10, "y": 25, "z": 5},
			expect:      10,
		},
		{
			description: "<<",
			expr:        "x<<y",
			state:       map[string]interface{}{"x": 1, "y": 3},
			expect:      8,
		},
		{
			description: ">>",
			expr:        "x>>y",
			state:       map[string]interface{}{"x": 16, "y": 2},
			expect:      4,
		},
		{
			description: "struct x+y",
			expr:        "foo.X + foo.Y",
			state:       map[string]interface{}{"foo": Foo{X: 2, Y: 3}},
			expect:      5,
		},
	}

	for _, testCase := range testCase {
		scope := NewScope()
		for k, v := range testCase.state {
			_, err := scope.DefineVariable(k, reflect.TypeOf(v))
			assert.Nil(t, err, "failed to define: "+k)
		}
		anExpr, err := scope.IntExpression(testCase.expr)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		for k, v := range testCase.state {
			if iV, ok := v.(int); ok {
				err = anExpr.State.SetInt(k, iV)
				assert.Nil(t, err)
			} else {
				err = anExpr.State.SetValue(k, v)
				assert.Nil(t, err)
			}
		}
		actual := anExpr.Compute()
		assert.Equal(t, testCase.expect, actual)
	}
}

func TestScope_BoolExpression(t *testing.T) {

	type Foo struct {
		X int
		Y int
	}
	var fooPtr *Foo
	var err error

	var testCase = []struct {
		description string
		expr        string
		state       map[string]interface{}
		expect      bool
	}{
		{
			description: "err != nil",
			expr:        "err != nil",
			state:       map[string]interface{}{"err": io.EOF},
			expect:      true,
		},
		{
			description: "x > y",
			expr:        "x>y",
			state:       map[string]interface{}{"x": 10, "y": 20},
			expect:      false,
		},

		{
			description: "x == y",
			expr:        "x == y",
			state:       map[string]interface{}{"x": 10.0, "y": 10.0},
			expect:      true,
		},
		{
			description: "x >= y",
			expr:        "x >= y",
			state:       map[string]interface{}{"x": 9.0, "y": 10.0},
			expect:      false,
		},
		{
			description: "foo != nil",
			expr:        "foo != nil",
			state:       map[string]interface{}{"foo": &Foo{X: 1}},
			expect:      true,
		},
		{
			description: "foo == nil",
			expr:        "foo == nil",
			state:       map[string]interface{}{"foo": fooPtr},
			expect:      true,
		},
		{
			description: "err == nil",
			expr:        "err == nil",
			state:       map[string]interface{}{"err": err},
			expect:      true,
		},
		{
			description: `name == ""`,
			expr:        `name == ""`,
			state:       map[string]interface{}{"name": ""},
			expect:      true,
		},
		{
			description: `name != ""`,
			expr:        `name != ""`,
			state:       map[string]interface{}{"name": "abc"},
			expect:      true,
		},
		{
			description: "bool addition",
			expr:        "(x < y)",
			state:       map[string]interface{}{"x": 10, "y": 20},
			expect:      true,
		},
		{
			description: "bool addition",
			expr:        "! (x < y)",
			state:       map[string]interface{}{"x": 10, "y": 20},
			expect:      false,
		},
	}

	for _, testCase := range testCase {
		scope := NewScope()
		for k, v := range testCase.state {
			vType := reflect.TypeOf(v)
			if k == "err" {
				vType = errType
			}
			_, err := scope.DefineVariable(k, vType)
			assert.Nil(t, err, "failed to define: "+k)
		}
		anExpr, err := scope.BoolExpression(testCase.expr)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		for k, v := range testCase.state {
			if iV, ok := v.(string); ok {
				err = anExpr.State.SetString(k, iV)
				assert.Nil(t, err)
			} else {
				err = anExpr.State.SetValue(k, v)
				assert.Nil(t, err)
			}
		}
		actual := anExpr.Compute()
		assert.Equal(t, testCase.expect, actual, testCase.description)
	}
}

func TestScope_StringExpression(t *testing.T) {
	type Foo struct {
		X string
		Y string
	}
	var testCase = []struct {
		description string
		expr        string
		state       map[string]interface{}
		expect      string
	}{
		{
			description: "basic concatenation",
			expr:        "x+y",
			state:       map[string]interface{}{"x": "abc", "y": "def"},
			expect:      "abcdef",
		},
		{
			description: "basic concatenation",
			expr:        "foo.X+y",
			state:       map[string]interface{}{"foo": Foo{X: "this is "}, "y": "test"},
			expect:      "this is test",
		},
	}

	for _, testCase := range testCase {
		scope := NewScope()
		for k, v := range testCase.state {
			vType := reflect.TypeOf(v)
			if k == "err" {
				vType = errType
			}
			_, err := scope.DefineVariable(k, vType)
			assert.Nil(t, err, "failed to define: "+k)
		}
		anExpr, err := scope.StringExpression(testCase.expr)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		for k, v := range testCase.state {
			if iV, ok := v.(int); ok {
				err = anExpr.State.SetInt(k, iV)
				assert.Nil(t, err)
			} else {
				err = anExpr.State.SetValue(k, v)
				assert.Nil(t, err)
			}
		}
		actual := anExpr.Compute()
		assert.Equal(t, testCase.expect, actual, testCase.description)
	}
}

func TestScope_Float64Expression(t *testing.T) {
	type Foo struct {
		X float64
		Y string
	}
	var testCase = []struct {
		description string
		expr        string
		state       map[string]interface{}
		expect      float64
	}{
		{
			description: "basic +",
			expr:        "x+y",
			state:       map[string]interface{}{"x": 1.4, "y": 3.4},
			expect:      4.8,
		},
		{
			description: "struct +",
			expr:        "foo.X+y",
			state:       map[string]interface{}{"foo": Foo{X: 2.2}, "y": 2.2},
			expect:      4.4,
		},
	}

	for _, testCase := range testCase {
		scope := NewScope()
		for k, v := range testCase.state {
			vType := reflect.TypeOf(v)
			if k == "err" {
				vType = errType
			}
			_, err := scope.DefineVariable(k, vType)
			assert.Nil(t, err, "failed to define: "+k)
		}
		anExpr, err := scope.Float64Expression(testCase.expr)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		for k, v := range testCase.state {
			if iV, ok := v.(float64); ok {
				err = anExpr.State.SetFloat64(k, iV)
				assert.Nil(t, err)
			} else {
				err = anExpr.State.SetValue(k, v)
				assert.Nil(t, err)
			}
		}
		actual := anExpr.Compute()
		assert.Equal(t, testCase.expect, actual, testCase.description)
	}
}
