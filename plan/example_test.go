package plan_test

import (
	"fmt"
	"github.com/viant/igo/plan"
	"log"
	"reflect"
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

	vars := varsNew() //there could be
	executor.Exec(vars)
	result, _ := vars.Int("s")
	fmt.Printf("result: %v\n", result)

}

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
	err := scope.RegisterType(reflect.TypeOf(Account{}))
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
	result, _ := vars.Float64("s")
	fmt.Printf("result: %v\n", result)

}

func ExampleScope_Function() {
	type Foo struct {
		Z int
	}
	scope := plan.NewScope()
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
	r := typeFn(1, 2, Foo{3})
	fmt.Printf("%v\n", r)
}
