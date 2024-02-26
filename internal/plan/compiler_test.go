package plan

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/notify"
	"github.com/viant/igo/option"
	"log"
	"reflect"
	"strconv"
	"testing"
	"unsafe"
)

func TestScope_Function(t *testing.T) {
	{
		scope := NewScope()

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
	type Foo struct {
		Z int
	}
	{
		scope := NewScope()
		scope.RegisterType(reflect.TypeOf(Foo{}))
		fn, err := scope.Function(`func(x, y int, foo Foo) int {
		return (x+y)/foo.Z
	}`)
		if err != nil {
			log.Fatalln(err)
		}
		typeFn, ok := fn.(func(int, int, Foo) int)
		if !ok {
			log.Fatalf("expected: %T, but had: %T", typeFn, fn)
		}
		actual := typeFn(9, 3, Foo{3})
		assert.Equal(t, 4, actual)
	}

	{
		scope := NewScope()
		fn, err := scope.Function(`func(x, y float64) int {
		return int(x+y)
	}`)
		assert.Nil(t, err)
		actual, ok := fn.(func(float64, float64) int)
		if !assert.True(t, ok) {
			return
		}
		assert.Equal(t, 11, actual(6.4, 5.1))
	}
}

func TestScope_Compile(t *testing.T) {
	type Bar struct {
		ID1   int
		Name1 string
	}
	type Account struct {
		Total float64
	}

	type Foo struct {
		ID   int
		Name string
		Bar  *Bar
	}

	var testCases = []struct {
		description  string
		scope        *Scope
		code         string
		expect       interface{}
		fnName       string
		fn           interface{}
		resultPath   string
		variableName string
		variable     interface{}
		customTypes  []reflect.Type
	}{

		{
			description: "basic assigment",
			scope:       NewScope(),
			code: `
			 i := 202
			 z := i
			`,
			resultPath: "z",
			expect:     202,
		},
		{
			description: "basic assigment",
			scope:       NewScope(),
			code: `
			 i := 30
			 i++
			 z := i
			`,
			resultPath: "z",
			expect:     31,
		},
		{
			description: "basic arithmetic expr with scopes",
			scope:       NewScope(),
			code: ` r := 0
			{
			i := 101
			{
				z := 230;
				{
					k := 23;
					r = i + (z * k) / 6
				}
			}
	}
`,
			resultPath: "r",
			expect:     101 + (230*23)/6,
		},

		{
			description: "basic arithmetic expr",
			scope:       NewScope(),
			code: ` i := 101
z := 230;
k := 23;
r := i + (z * k) / 6
`,
			resultPath: "r",
			expect:     101 + (230*23)/6,
		},

		{
			description: "basic arithmetic expr 2",
			scope:       NewScope(),
			code: ` a := 5
b := 7
c := 9
d := 1			
r := a + b + 10 * c + d
`,
			resultPath: "r",
			expect:     5 + 7 + 10*9 + 1,
		},
		{
			description: "basic self assigment",
			scope:       NewScope(),
			code: `
			 i := 202
			 z := 101
			 z += i
			`,
			resultPath: "z",
			expect:     303,
		},
		{
			description: "function call - max",
			code:        `r := max(12,3)`,
			fnName:      "max",
			scope:       NewScope(),
			fn: func(x, y int) int {
				if x > y {
					return x
				}
				return y
			},
			resultPath: "r",
			expect:     12,
		},
		{
			description: "function call - parse",
			code:        `r, e := parse("true")`,
			fnName:      "parse",
			fn:          testFn(strconv.ParseBool),
			resultPath:  "r",
			expect:      true,
		},
		{
			description: "struct assigment",
			scope:       NewScope(),
			code: `
			 n := "abc"
			 foo := &Foo{ID:123, Name: n}
			 z := foo.Name
			`,
			resultPath:  "z",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      "abc",
		},
		{
			description: "nested struct assigment",
			scope:       NewScope(),
			code: `
			 n := "abc"
			 foo := &Foo{ID:123, Name: n, Bar: &Bar{ID1:456}}
			 z := foo.Bar.ID1
			`,
			resultPath:  "z",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      456,
		},

		{
			description: "slice assigment []*Foo{}",
			code: `
			 n := "abc"
			 foos := []*Foo{&Foo{ID:123, Name: n}}
			 f0 := foos[0]
			 z := f0.ID
			`,
			resultPath:  "z",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      123,
		},
		{
			description: "slice assigment []Foo{}",
			code: `
			 n := "abc"
			 foos := []Foo{Foo{ID:123, Name: n}}
			 f0 := foos[0]
			 z := f0.ID
			`,
			resultPath:  "z",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      123,
		},

		{
			description: "if-else statement",
			code: `a := 10
b := 11
r := 0
			if b < a {
				r = 101
			} else {
				r = 222
			}
			`,
			resultPath: "r",
			expect:     222,
		},

		{
			description: "if-else nested statement",
			code: `a := 10
r := 0
			if a < 4 {
				r = 111
			} else if a < 8 {
				r = 222
			} else if a < 12 {
				r = 333
			}
			`,
			resultPath: "r",
			expect:     333,
		},

		{
			description: "if-else statement with return",
			code: `a := 24
b := 11
r := 0
			if  a >  b {
				r = 101
				return 
			} 
			r = 200
			
			`,
			resultPath: "r",
			expect:     101,
		},

		{
			description: "for statement",
			code: `a := 10
			r := 0
			for i :=0;i<a;i++{
				r +=  i * i
			}
			`,
			resultPath: "r",
			expect:     285,
		},
		{
			description: "slice append",
			code: `
			var foos = make([]*Foo, 0)
			for i:=0;i<10;i++ {
				foos = append(foos, &Foo{ID:i, Name:"nxc"})
			}
			z := len(foos)
			`,
			resultPath:  "z",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      10,
		},
		{
			description: "for statement with break",
			code: `i := 0
			for z :=0;z<20; {
				i++										
				if i > 10 {
					break
				}
				z++
			}
			`,
			resultPath: "i",
			expect:     11,
		},

		{
			description: "range []Foo{}",
			code: `
			 n := "abc"
			 foos := []Foo{Foo{ID:10, Name: n}, Foo{ID:20, Name: n}}
	  		 z := 0
			 for i, foo := range foos {
				z += foo.ID + i
			 }
			`,
			resultPath:  "z",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      31,
		},
		{
			description: "range []*Foo{}",
			code: `
			 n := "abc"
			 foos := []*Foo{&Foo{ID:15, Name: n}, &Foo{ID:20, Name: n}}
	  		 z := 0
			 for i, foo := range foos {
				z += foo.ID + i
			 }
			`,
			resultPath:  "z",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      36,
		},

		{
			description:  "set custom",
			variableName: "c",
			variable:     &ctest{},
			code: `
			c.SetActive(true)
			z := c.Active
			`,
			resultPath: "z",
			expect:     true,
		},

		{
			description: "if-else statement with return  - explicit variable result",
			code: `func() (x int) {
			a := 24
			b := 12
			if  a >  b {
				return 12 
			} 
			return 32
}`,
			resultPath: "x",
			expect:     12,
		},

		{
			description: "if-else statement with return",
			code: `func() int {
			a := 24
			b := 12
			if  a >  b {
				return 12 
			} 
			return 32
}`,
			resultPath: "Result0",
			expect:     12,
		},
		{
			description: "type def",
			code: `type Foo struct {
		ID int
		Name string
	}
	var foos = make([]*Foo, 0)
	for i:=0;i<10;i++ {
		foos = append(foos, &Foo{ID:i, Name:"nxc"})
	}
	s := 0
	for i, foo := range foos {
		if i %2  == 0 {
			s += foo.ID		
		}
	}`,
			resultPath: "s",
			expect:     20,
		},

		{
			description: "define variable",
			code: `
	x := 0.0
	for _, account := range accounts {
		x += account.Total
	}
	`,
			variableName: "accounts",
			variable: []Account{
				{Total: 1.3},
				{Total: 3.7},
			},
			expect:     5.0,
			resultPath: "x",
		},

		{
			description: "slice.Reduce",
			code: `
			 n := "abc"
			 foos := []*Foo{&Foo{ID:15, Name: n}, &Foo{ID:20, Name: n}}
			 z := foos.Reduce(func(acc int, foo *Foo) int {
				return acc + foo.ID
			}, 1)
			`,
			resultPath:  "z",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      36,
		},
		{
			description: "slice.Map",
			code: `
			 ints := []int{1,2,3,4}
			 r := ints.Map(func(i int, index int) int {
				return index + i * 2
			})
			`,
			resultPath:  "r",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      []int{2, 5, 8, 11},
		},

		{
			description: "range []*Foo{}",
			code: `
			 n := "abc"
			 foos := []*Foo{&Foo{ID:15, Name: n}, &Foo{ID:20, Name: n}}
	  		 bars := foos.Map(func(foo *Foo) *Bar {
				return &Bar{ID1:foo.ID, Name1: foo.Name}
			 })
			`,
			resultPath:  "bars",
			customTypes: []reflect.Type{reflect.TypeOf(Foo{}), reflect.TypeOf(Bar{})},
			expect:      []*Bar{{ID1: 15, Name1: "abc"}, {ID1: 20, Name1: "abc"}},
		},
	}

	for i, testCase := range testCases {
		fmt.Printf("[%v]: %v\n", i, testCase.code)
		if testCase.scope == nil {
			testCase.scope = NewScope()
		}
		if testCase.fn != "" {
			testCase.scope.RegisterFunc(testCase.fnName, testCase.fn)
		}
		if len(testCase.customTypes) > 0 {
			for _, tp := range testCase.customTypes {
				err := testCase.scope.RegisterType(tp)
				assert.Nil(t, err, testCase.description)
			}
		}
		if testCase.variableName != "" {
			testCase.scope.DefineVariable(testCase.variableName, reflect.TypeOf(testCase.variable))
		}
		execution, err := testCase.scope.Compile(testCase.code)
		if !assert.Nil(t, err, testCase.description) {
			fmt.Println(err)
			continue
		}
		state := execution.NewState()
		if testCase.variableName != "" {
			state.SetValue(testCase.variableName, testCase.variable)
		}
		execution.Exec(state)
		actual, err := state.Value(testCase.resultPath)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		if !assert.Equal(t, testCase.expect, actual, testCase.description) {

			continue
		}
	}

}

func TestDefineEmbedVariable(t *testing.T) {
	type B struct {
		Acc int
	}
	type A struct {
		Count       int
		ActiveCount int
		B           B
	}

	//	scope := NewScope(option.NewTracker("State", reflect.TypeOf(A{}), true))
	scope := NewScope()
	scope.DefineEmbedVariable("State", reflect.TypeOf(A{}))
	scope.DefineVariable("t", reflect.TypeOf(true))
	exec, err := scope.Compile(`
	r := 0
	if t  {
		Count = 10
		B.Acc = 31
		r = Count * B.Acc
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := exec.NewState()
	state.SetBool("t", true)

	//tracker := state.Tracker()
	//test mutation
	{
		state.SetBool("t", true)
		exec.Exec(state)
		actual, _ := state.Int("r")
		assert.Equal(t, 310, actual)
	}
}

func TestNewTrackedScope(t *testing.T) {
	type B struct {
		Acc int
	}
	type A struct {
		Count       int
		ActiveCount int
		B           B
	}

	scope := NewScope(option.WithTracker(notify.NewTracker("a", reflect.TypeOf(A{}), false)))
	scope.DefineVariable("t", reflect.TypeOf(true))
	exec, err := scope.Compile(`
	r := 0
	if t  {
		a.Count = 10
		a.B.Acc = 31
		r = a.Count * a.B.Acc
	}
	`)
	if !assert.Nil(t, err) {
		return
	}
	state := exec.NewState()

	count, _ := state.Selector("a.Count")
	acc, _ := state.Selector("a.B.Acc")

	tracker := state.Tracker()
	//test mutation
	{
		state.SetBool("t", true)
		exec.Exec(state)
		actual, _ := state.Int("r")
		assert.Equal(t, 310, actual)
		assert.True(t, tracker.Has(count.Pos), "count tracker")
		assert.True(t, tracker.Has(acc.Pos), "acc tracker")

	}
	//test no  mutation
	{
		tracker.Reset()
		state.SetBool("t", false)
		exec.Exec(state)
		actual, _ := state.Int("r")
		assert.Equal(t, 0, actual)
		assert.False(t, tracker.Has(count.Pos), "count tracker")
		assert.False(t, tracker.Has(acc.Pos), "acc tracker")
	}

}

type ctest struct {
	Active bool
}

func (c *ctest) SetActive(b bool) {
	c.Active = b
}

type testFn func(s string) (bool, error)

func (f testFn) Call(ptr unsafe.Pointer, args []*exec.Operand) unsafe.Pointer {
	s := args[0].Compute(ptr)
	r, e := f(*(*string)(s))
	return unsafe.Pointer(&[2]unsafe.Pointer{
		unsafe.Pointer(&r), unsafe.Pointer(&e),
	})
}
