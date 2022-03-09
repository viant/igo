# Igo (go evaluator in go)

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
within reasonable time. Some existing alternative providing go evaluation on the fly like [GoEval](https://github.com/xtaci/goeval)
[GoVal](github.com/maja42/goval) are prohibitively slow.

See [performance](#performance) section for details.

## Introduction

In order to reduce execution time, this project first produces execution plan alongside with all variables needed to execute it.
One execution plan can be shared alongside many instances of scoped variables needed by executor. 
Scoped Variables holds both execution state  and  variables defined or used in the evaluation code.

```go
    scope := plan.NewScope()
	executor, varsNew, err := scope.Compile(code)
    vars := varsNew() //creates memory instance needed by executor

    executor.Exec(vars)
	
```

## Usage

### Expression

```go
package mypkg

import (
	"log"
	"reflect"
	"github.com/viant/igo/plan"
)

func ExampleScope_BoolExpression() {
	type Performance struct {
		Id      int
		Price   float64
		MetricX float64
	}
	scope := plan.NewScope()
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
	"github.com/viant/igo/plan"
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

	scope := plan.NewScope()
	executor, varsNew, err := scope.Compile(code)
	if err != nil {
		log.Fatalln(err)
	}

	vars := varsNew() //variables constructor, one per each concurent execution, execution can be shared
	executor.Exec(vars)
	result, _ := vars.Int("s")
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
	"github.com/viant/igo/plan"
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

	scope := plan.NewScope()
	err := scope.RegisterType(reflect.TypeOf(Account{})) //Register all non-primitive types used in code 
	if err != nil {
		log.Fatalln(err)
	}
	executor, varsNew, err := scope.Compile(code)
	if err != nil {
		log.Fatalln(err)
	}
	vars := varsNew()
	err = vars.SetValue("accounts", []Account{
		{Total: 1.3},
		{Total: 3.7},
	})
	if err != nil {
		log.Fatalln(err)
	}
	executor.Exec(vars)
	result, _ := vars.Float64("x")
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
	"github.com/viant/igo/plan"
)

func ExampleScope_Function() {
    type Foo struct {
        Z int
    }
    scope := plan.NewScope()
    _ = scope.RegisterType(reflect.TypeOf(Foo{}))
    fn, err := scope.Function(`func(x, y int, foo Foo) int {
            return (x+y)/foo.Z
        }`)
    if err != nil {
        log.Fatalln(err)
    }
    typeFn, ok := fn.(func(int, int, Foo) int)
    if ! ok {
     log.Fatalf("expected: %T, but had: %T", typeFn,fn)
    }
    r := typeFn(1, 2, Foo{3})
    fmt.Printf("%v\n", r)
}
```
## Registering types

To use data types defined outside the code, register type with `(Scope).RegisterType(type)` function or
`(Scope).RegisterNamedType(name, type)`

```go
    scope := plan.NewScope()
    _ = scope.RegisterType(reflect.TypeOf(Foo{}))

```


DefineVariable
## Registering function

To use function defined outside the code, register type with `(Scope).RegisterFunc(name, function)` function

```go
    scope := plan.NewScope()
    scope.RegisterFunc(testCase.fnName, testCase.fn)

```



## Performance

### Expression evaluation

See benchmark for the following expression evaluation:

```10 + (5 * x / y * (z - 7))```

```text
BenchmarkScope_IntExpression_Native-16          134789700                8.627 ns/op           0 B/op          0 allocs/op
BenchmarkScope_IntExpression-16                 19722770                57.06 ns/op            8 B/op          1 allocs/op
BenchmarkScope_IntExpression_GoVal-16             625620                2040 ns/op          2328 B/op        11 allocs/op
```

[GoVal](github.com/maja42/goval) is ~255 slower for the presented expression comparing to natively compiled code, 
while Igo is only ~7 times slower


### Code execution

See benchmark for the following code:
```go
count :=0
for i :=0;i<100;i++ {
    count += i
}
print(count)
```

[GoEval](https://github.com/xtaci/goeval) evaluation takes almost ~24K time longer than natively compiled code,
whereas this project is only around ~35 slower. As point of reference using native go reflection adds on average
around 100x time execution overhead.

```text
Benchmark_Loop_Native-16                        35385890               30.36 ns/op             0 B/op          0 allocs/op
Benchmark_Loop_Igo-16                            1000000              1081  ns/op               0 B/op          0 allocs/op
Benchmark_Loop_GoEval-16                            1429            739672  ns/op          788350 B/op       3180 allocs/op
```

See the following benchmark that runs 100 000 000 loop iteration:
```text
	z := 0
	a := 100000000
	r := 1
	for i := 1; i <= a; i++ {
		r += i
	}
	z = r
```

```text
BenchmarkLongLoop_Igo-16                          2         661982642 ns/op               8 B/op          0 allocs/op
BenchmarkLongLoop_Native-16                      48          24813792 ns/op               0 B/op          0 allocs/op
```
Igo is ~26x times slower than natively compile code.

## Bugs

This project does not implement full golang spec, but just a subset.
At least following expression/types/construct are not supported
- map type
- named interface types (since pointers are used to access/mutate data)
- go routines
- select expression
- switch expression
- closures


## Contributing to igo

Igo is an open source project and contributors are welcome!

See [TODO](TODO.md) list

## Credits and Acknowledgements

**Library Author:** Adrian Witas

