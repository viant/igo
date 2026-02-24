# igo (go evaluator in go)

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/igo)](https://goreportcard.com/report/github.com/viant/igo)
[![GoDoc](https://godoc.org/github.com/viant/igo?status.svg)](https://godoc.org/github.com/viant/igo)

This library is compatible with Go 1.17+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.

- [Motivation](#motivation)
- [Introduction](#introduction)
- [Usage](#usage)
- [Performance](#performance)
- [Bugs](#bugs)
- [Contribution](#contributing-to-igo)
- [License](#license)

## Motivation

The goal of this library is to be able dynamically execute go code directly from Go/WebAssembly 
within reasonable time. Some existing alternative providing go evaluation on the fly are prohibitively slow: 
- [GoEval](https://github.com/xtaci/goeval)
- [GoVal](https://github.com/maja42/goval) 
- [Yaegi](https://github.com/traefik/yaegi) .

See [performance](#performance) section for details.

## Introduction

In order to reduce execution time, this project first produces execution plan alongside with state needed to execute it.
One execution plan can be shared alongside many instances state needed by executor. 
State holds both variables and execution state used in the evaluation code.

```go
package mypkg

import "github.com/viant/igo"

func usage() {
	scope := igo.NewScope()
	code := "go code here"
	executor, stateNew, err := scope.Compile(code)
	if err != nil {
		panic(err)
	}
	state := stateNew() //creates memory instance needed by executor 
	executor.Exec(state)
}
	
```

## Usage

### Expression

```go
package mypkg

import (
	"log"
	"reflect"
	"github.com/viant/igo"
)

func ExampleScope_BoolExpression() {
	type Performance struct {
		Id      int
		Price   float64
		MetricX float64
	}
	scope := igo.NewScope()
	_, err := scope.DefineVariable("perf", reflect.TypeOf(Performance{}))
	_, err = scope.DefineVariable("threshold", reflect.TypeOf(0.0))
	if err != nil {
		log.Fatalln(err)
	}
	//Compile bool expression
	expr, err := scope.BoolExpression("perf.MetricX > threshold && perf.Price > 1.0")
	if err != nil {
		log.Fatalln(err)
	}

	perfs := []Performance{
		{MetricX: 1.5, Price: 3.2},
		{MetricX: 1.2, Price: 1.2},
		{MetricX: 1.7, Price: 0.4},
	}
	var eval = make([]bool, len(perfs))
	for i := range perfs {
		_ = expr.Vars.SetValue("perf", perfs[i])
		_ = expr.Vars.SetFloat64("threshold", 0.7)
		eval[i] = expr.Compute()
	}
}
```

### Go evaluation

```go
package mypkg

import (
	"log"
	"fmt"
	"github.com/viant/igo"
)

func ExampleScope_Compile() {
	code := `type Foo struct {
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
}`

	scope := igo.NewScope()
	executor, stateNew, err := scope.Compile(code)
	if err != nil {
		log.Fatalln(err)
	}

	state := stateNew() //variables constructor, one per each concurent execution, execution can be shared
	executor.Exec(state)
	result, _ := state.Int("s")
	fmt.Printf("result: %v\n", result)
}
```

Setting code variables

```go
package mypkg

import (
	"log"
	"fmt"
	"reflect"
	"github.com/viant/igo"
)

func ExampleScope_DefineVariable() {
	code := `
	x := 0.0
	for _, account := range accounts {
		x += account.Total
	}
	`
	type Account struct {
		Total float64
	}

	scope := igo.NewScope()
	err := scope.RegisterType(reflect.TypeOf(Account{})) //Register all non-primitive types used in code 
	if err != nil {
		log.Fatalln(err)
	}
	executor, stateNew, err := scope.Compile(code)
	if err != nil {
		log.Fatalln(err)
	}
	state := stateNew()
	err = state.SetValue("accounts", []Account{
		{Total: 1.3},
		{Total: 3.7},
	})
	if err != nil {
		log.Fatalln(err)
	}
	executor.Exec(state)
	result, _ := state.Float64("x")
	fmt.Printf("result: %v\n", result)
}
```

### Go function

```go
package mypkg

import (
	"log"
	"fmt"
	"reflect"
	"github.com/viant/igo"
)

func ExampleScope_Function() {
	type Foo struct {
		Z int
	}
	scope := igo.NewScope()
	_ = scope.RegisterType(reflect.TypeOf(Foo{}))
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
	r := typeFn(1, 2, Foo{3})
	fmt.Printf("%v\n", r)
}
```
## Registering types

To use data types defined outside the code, register type with `(Scope).RegisterType(type)` function or
`(Scope).RegisterNamedType(name, type)`

```go
    scope := igo.NewScope()
    _ = scope.RegisterType(reflect.TypeOf(Foo{}))

```


DefineVariable
## Registering function

To use function defined outside the code, register type with `(Scope).RegisterFunc(name, function)` function

```go
    scope := igo.NewScope()
    scope.RegisterFunc(testCase.fnName, testCase.fn)

```



## Performance

### Current benchmark snapshot (Apple M1 Max, Feb 24 2026)

Command:

```bash
go test ./internal/plan/bench -run '^$' -bench '.' -benchmem -count=1
```

```text
Benchmark_Loop_Native-10                     34.09 ns/op      0 B/op      0 allocs/op
Benchmark_Loop_Igo-10                        239.0 ns/op      0 B/op      0 allocs/op

Benchmark_ConditionalCall_Native-10          80.15 ns/op      0 B/op      0 allocs/op
Benchmark_ConditionalCall_Igo-10             251.5 ns/op      0 B/op      0 allocs/op

Benchmark_ComplexState_Native-10             1846 ns/op       0 B/op      0 allocs/op
Benchmark_ComplexState_Igo-10                1898 ns/op       0 B/op      0 allocs/op

Benchmark_MergeSort_Native-10                3921 ns/op       0 B/op      0 allocs/op
Benchmark_MergeSort_Igo-10                   30703 ns/op      0 B/op      0 allocs/op

Benchmark_RowFilter_Native-10                169.5 ns/op      64 B/op     4 allocs/op
Benchmark_RowFilter_Igo-10                   1096 ns/op       160 B/op    8 allocs/op

Benchmark_InlineHelperAdapter_Native-10      104.0 ns/op      8 B/op      1 allocs/op
Benchmark_InlineHelperAdapter_Igo-10         1419 ns/op       472 B/op    18 allocs/op
```

The conditional-call benchmark reflects the fused fast path for:
- `for i := ...; ...; ... { if i <K/op { fn(i) } }`
- `for i := ...; ...; ... { if i % M <op> K { fn(i) } }`

with `<op>` in `{==, !=, <, <=, >, >=}` and `fn` registered as `func(int)` or `func(int) int`.

Additional workload notes from current snapshot:
- `ComplexState_Igo` is ~1.03x native with zero allocations.
- `MergeSort_Igo` is ~7.8x native with zero allocations.
- `RowFilter_Igo` is ~6.5x native; allocation gap is now 8 vs 4 allocs/op.
- `InlineHelperAdapter_Igo` is ~13.6x native; recent inline helper/append/toLower-path optimizations reduced this workload significantly.

### Notes on specialization strategy

Current speedups come from aggressive compile-time specialization in hot paths, not a bytecode VM:
- typed `for` loop kernels (common arithmetic and call patterns)
- typed map/slice fast paths for frequently used key/value/index shapes
- branch-heavy fused loop bodies (used by complex-state style workloads)
- stateful variadic callers for common signatures to avoid hidden allocations

## Language support status

igo supports a substantial Go subset with runtime and parser coverage tests.

### Supported (current)

- Declarations:
  - short var declarations and typed var declarations
  - type specs for `struct`, `interface`, `map`, and aliases
- Literals:
  - int, float, string, rune, imaginary
- Control flow:
  - `if`
  - `for` (including specialized int-loop variants)
  - `range` over slices/maps/strings (with key/value forms covered by tests)
  - `switch` (tag and tagless), including `fallthrough`
  - `type switch`
  - `break`, `continue`, `goto`, and labels (including labeled break/continue)
  - `defer`
  - `return`
- Types/operations:
  - type assertions (`v := x.(T)` and `v, ok := x.(T)`)
  - map read/write, `value,ok`, `delete`, range over map
  - builtins used in tests/benchmarks: `len`, `make`, `append`, `delete`
- Calls:
  - registered functions
  - variadic call paths, including optimized stateful caller paths for common signatures
- Stdlib:
  - lazy package/type/function resolution when referenced in expressions/scripts
  - tests include examples with `strings`, `fmt`, and `time`

### Not supported yet (known)

- `go` statements (goroutines in script)
- `select` statements
- channel send statements (`ch <- v`)

This is still not full Go spec coverage; unsupported constructs return compile-time errors.


## Contributing to igo

Igo is an open source project and contributors are welcome!

See [TODO](TODO.md) list

## Credits and Acknowledgements

**Library Author:** Adrian Witas
