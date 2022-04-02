package bench

import (
	"github.com/maja42/goval"
	"github.com/stretchr/testify/assert"
	"github.com/traefik/yaegi/interp"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/internal"
	"github.com/viant/igo/internal/exec"
	"github.com/viant/igo/internal/expr"
	"github.com/viant/igo/internal/plan"
	"github.com/xtaci/goeval"
	"reflect"
	"testing"
)

var benchExecution *internal.Executor
var benchVars *exec.State
var benchNative = func() int {
	z := 0
	a := 100000000
	r := 1
	for i := 1; i <= a; i++ {
		r += i
	}
	z = r
	return z
}

var benchLoopCode = `
	z := 0
	a := 100000000
	r := 1
	for i := 1; i <= a; i++ {
		r += i
	}
	z = r
	`

func initForBench() {

	scope := plan.NewScope()
	var newVars exec.New
	var err error
	benchExecution, newVars, err = scope.Compile(benchLoopCode)
	if err != nil {
		panic(err)
	}
	benchVars = newVars()

}

func BenchmarkLongLoop_Igo(b *testing.B) {
	initForBench()
	b.ResetTimer()
	for k := 0; k < b.N; k++ {
		b.ReportAllocs()
		benchExecution.Exec(benchVars)
	}
	z, err := benchVars.Int("z")
	assert.Nil(b, err)
	assert.Equal(b, 5000000050000001, z)
}

func BenchmarkLongLoop_Native(b *testing.B) {
	z := 0

	for k := 0; k < b.N; k++ {
		b.ReportAllocs()
		z = benchNative()
	}
	assert.Equal(b, 5000000050000001, z)
}

func BenchmarkLoop_Yaegi(b *testing.B) {
	b.ReportAllocs()
	for k := 0; k < b.N; k++ {
		i := interp.New(interp.Options{})
		_, err := i.Eval(benchLoopCode)
		assert.Nil(b, err)
	}
}

func testPrintln(a ...interface{}) (n int, err error) {
	return 0, nil
}

var loopNative func()

func Benchmark_Loop_Native(b *testing.B) {
	loopNative = func() {
		count := 0
		for i := 0; i < 100; i++ {
			count = count + i
		}
		testPrintln(count)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		loopNative()
	}
}

func Benchmark_Loop_GoEval(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := goeval.NewScope()
		s.Set("print", testPrintln)
		s.Eval(`count := 0`)
		s.Eval(`for i:=0;i<100;i=i+1 {
			count=count+i
		}`)
		s.Eval(`print(count)`)
	}
}

func init() {
	initLoopIgo()
}
func Benchmark_Loop_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchLoopIgo.Exec(benchLoopIgoVars)
	}
}

var benchLoopIgo *internal.Executor
var benchLoopIgoVars *exec.State

func initLoopIgo() {
	scope := plan.NewScope()
	scope.RegisterFunc("print", testPrintln)
	var stateNew exec.New
	var err error
	benchLoopIgo, stateNew, err = scope.Compile(`count :=0
for i :=0;i<100;i++ {
	count += i
}
print(count)	
	`)
	if err != nil {
		panic(err)
	}
	benchLoopIgoVars = stateNew()
}

var benchIntNativeExpr func(x, y, z int) int
var benchIntExpr *expr.Int
var intType = reflect.TypeOf(0)

func initIntExpressionBench() {

	var err error
	scope := plan.NewScope()
	scope.DefineVariable("x", intType)
	scope.DefineVariable("y", intType)
	scope.DefineVariable("z", intType)
	benchIntExpr, err = scope.IntExpression("10 + (5 * x / y * (z - 7))")
	if err != nil {
		panic(err)
	}

	benchIntNativeExpr = func(x, y, z int) int {
		return 10 + (5 * x / y * (z - 7))
	}

}

func BenchmarkScope_IntExpression(b *testing.B) {
	initIntExpressionBench()
	k := 0
	benchIntExpr.State.SetInt("x", 10)
	benchIntExpr.State.SetInt("y", 20)
	benchIntExpr.State.SetInt("z", 30)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		k = benchIntExpr.Compute()
	}
	assert.Equal(b, 56, k)
}

func BenchmarkScope_IntExpression_Native(b *testing.B) {
	initIntExpressionBench()
	k := 0
	x := 10
	y := 20
	z := 30
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		k = benchIntNativeExpr(x, y, z)
	}
	assert.Equal(b, 56, k)
}

func BenchmarkScope_IntExpression_GoVal(b *testing.B) {
	k := 0
	x := 10
	y := 20
	z := 30
	eval := goval.NewEvaluator()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result, _ := eval.Evaluate("10 + (5 * x / y * (z - 7))", map[string]interface{}{
			"x": x,
			"y": y,
			"z": z,
		}, nil)
		k = result.(int)
	}
	assert.Equal(b, 56, k)
}
