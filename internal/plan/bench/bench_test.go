package bench

import (
	"github.com/maja42/goval"
	"github.com/stretchr/testify/assert"
	"github.com/traefik/yaegi/interp"
	"github.com/viant/igo/exec"
	"github.com/viant/igo/exec/expr"
	"github.com/viant/igo/internal/plan"
	"github.com/xtaci/goeval"
	"reflect"
	"strings"
	"testing"
)

var benchExecution *exec.Executor
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

type complexState struct {
	Sum    int
	Count  int
	Mode   int
	Output int
	Flags  map[string]bool
	Scores map[int]int
}

type filterRow struct {
	Name  string
	Score int
}

type adapterInRow struct {
	Name  string
	Score int
}

type adapterOutRow struct {
	Label string
	Value int
}

var mergeSortSrc = buildMergeSortInput(256)
var mergeSortExpectedSum = sumIntSlice(mergeSortSrc)

func buildMergeSortInput(size int) []int {
	result := make([]int, size)
	x := 7919
	for i := 0; i < size; i++ {
		x = (x*1103515245 + 12345) & 0x7fffffff
		result[i] = x % 100000
	}
	return result
}

func sumIntSlice(data []int) int {
	total := 0
	for i := 0; i < len(data); i++ {
		total += data[i]
	}
	return total
}

func mergeSortNativeRun(src, data, tmp []int) (int, int) {
	n := len(src)
	for i := 0; i < n; i++ {
		data[i] = src[i]
	}
	width := 1
	for width < n {
		left := 0
		for left < n {
			mid := left + width
			if mid > n {
				mid = n
			}
			right := mid + width
			if right > n {
				right = n
			}
			i := left
			j := mid
			k := left
			for i < mid && j < right {
				if data[i] <= data[j] {
					tmp[k] = data[i]
					i++
				} else {
					tmp[k] = data[j]
					j++
				}
				k++
			}
			for i < mid {
				tmp[k] = data[i]
				i++
				k++
			}
			for j < right {
				tmp[k] = data[j]
				j++
				k++
			}
			for p := left; p < right; p++ {
				data[p] = tmp[p]
			}
			left = right
		}
		width += width
	}

	ok := 1
	sum := 0
	for i := 0; i < n; i++ {
		sum += data[i]
		if i > 0 && data[i-1] > data[i] {
			ok = 0
		}
	}
	return ok, sum
}

func complexStateNative(st *complexState) {
	st.Sum = 0
	st.Count = 0
	st.Mode = 0
	st.Output = 0
	for i := 0; i < 200; i++ {
		if i&1 == 0 {
			st.Sum += i
		} else {
			st.Sum -= i / 2
		}
		switch i % 5 {
		case 0:
			st.Mode = 1
		case 1:
			st.Mode = 2
		case 2:
			st.Mode = 3
		case 3:
			st.Mode = 4
		default:
			st.Mode = 5
		}
		if st.Flags["hot"] {
			st.Output += st.Scores[i%4]
		}
		if st.Sum > 3000 {
			break
		}
		if i%7 == 0 {
			continue
		}
		st.Count++
	}
}

func complexStateNativeGuard(st *complexState) {
	st.Sum = 0
	st.Count = 0
	st.Mode = 0
	st.Output = 0
	for i := 0; i < 200; i++ {
		if i&1 == 0 {
			st.Sum += i
		} else {
			st.Sum -= i / 2
		}
		switch i % 5 {
		case 0:
			st.Mode = 1
		case 1:
			st.Mode = 2
		case 2:
			st.Mode = 3
		case 3:
			st.Mode = 4
		default:
			st.Mode = 5
		}
		if st.Flags["hot"] {
			st.Output += st.Scores[i%4]
		}
		if st.Sum >= 3000 {
			break
		}
		if i%7 != 0 {
			st.Count++
		}
	}
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
	var err error
	benchExecution, err = scope.Compile(benchLoopCode)
	if err != nil {
		panic(err)
	}
	benchVars = benchExecution.NewState()

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

var benchLoopIgo *exec.Executor
var benchLoopIgoVars *exec.State
var benchSwitchIgo *exec.Executor
var benchSwitchIgoVars *exec.State
var benchDeferIgo *exec.Executor
var benchDeferIgoVars *exec.State
var benchMapIgo *exec.Executor
var benchMapIgoVars *exec.State
var benchMapStringBoolIgo *exec.Executor
var benchMapStringBoolIgoVars *exec.State
var benchMapIntFloatIgo *exec.Executor
var benchMapIntFloatIgoVars *exec.State
var benchMapInt64FloatIgo *exec.Executor
var benchMapInt64FloatIgoVars *exec.State
var benchMapUint64FloatIgo *exec.Executor
var benchMapUint64FloatIgoVars *exec.State
var benchMapRangeIgo *exec.Executor
var benchMapRangeIgoVars *exec.State
var benchVariadicStringVoidIgo *exec.Executor
var benchVariadicStringVoidIgoVars *exec.State
var benchVariadicInterfaceIgo *exec.Executor
var benchVariadicInterfaceIgoVars *exec.State
var benchComplexStateIgo *exec.Executor
var benchComplexStateIgoVars *exec.State
var benchComplexStateGuardIgo *exec.Executor
var benchComplexStateGuardIgoVars *exec.State
var benchMergeSortIgo *exec.Executor
var benchMergeSortIgoVars *exec.State
var benchConditionalCallIgo *exec.Executor
var benchConditionalCallIgoVars *exec.State
var benchConditionalCallSink int
var benchRowFilterIgo *exec.Executor
var benchRowFilterIgoVars *exec.State
var benchRowFilterData []*filterRow
var benchInlineHelperAdapterIgo *exec.Executor
var benchInlineHelperAdapterIgoVars *exec.State
var benchInlineHelperAdapterData []*adapterInRow

func initLoopIgo() {
	scope := plan.NewScope()
	scope.RegisterFunc("print", testPrintln)
	var err error
	benchLoopIgo, err = scope.Compile(`count :=0
for i :=0;i<100;i++ {
	count += i
}
print(count)	
	`)
	if err != nil {
		panic(err)
	}
	benchLoopIgoVars = benchLoopIgo.NewState()
	initSwitchIgo()
	initDeferIgo()
	initMapIgo()
	initMapStringBoolIgo()
	initMapIntFloatIgo()
	initMapInt64FloatIgo()
	initMapUint64FloatIgo()
	initMapRangeIgo()
	initVariadicStringVoidIgo()
	initVariadicInterfaceIgo()
	initComplexStateIgo()
	initComplexStateGuardIgo()
	initMergeSortIgo()
	initConditionalCallIgo()
	initRowFilterIgo()
	initInlineHelperAdapterIgo()
}

func initSwitchIgo() {
	scope := plan.NewScope()
	var err error
	benchSwitchIgo, err = scope.Compile(`
		r := 0
		for i := 0; i < 100; i++ {
			switch i % 3 {
			case 0:
				r += 1
			case 1:
				r += 2
			default:
				r += 3
			}
		}
	`)
	if err != nil {
		panic(err)
	}
	benchSwitchIgoVars = benchSwitchIgo.NewState()
}

func initDeferIgo() {
	scope := plan.NewScope()
	scope.RegisterFunc("mark", func(v int) int {
		return v
	})
	var err error
	benchDeferIgo, err = scope.Compile(`
		defer mark(1)
		defer mark(2)
		x := 1
		_ = x
	`)
	if err != nil {
		panic(err)
	}
	benchDeferIgoVars = benchDeferIgo.NewState()
}

func initMapIgo() {
	scope := plan.NewScope()
	var err error
	benchMapIgo, err = scope.Compile(`
		m := map[string]int{"a": 1, "b": 2}
		sum := 0
		for i := 0; i < 100; i++ {
			v, ok := m["a"]
			if ok {
				sum += v
			}
			m["b"] = sum
		}
	`)
	if err != nil {
		panic(err)
	}
	benchMapIgoVars = benchMapIgo.NewState()
}

func initMapStringBoolIgo() {
	scope := plan.NewScope()
	var err error
	benchMapStringBoolIgo, err = scope.Compile(`
		m := map[string]bool{"on": true, "off": false}
		count := 0
		for i := 0; i < 100; i++ {
			v, ok := m["on"]
			if ok && v {
				count++
			}
			m["off"] = (count % 2) == 0
		}
	`)
	if err != nil {
		panic(err)
	}
	benchMapStringBoolIgoVars = benchMapStringBoolIgo.NewState()
}

func initMapIntFloatIgo() {
	scope := plan.NewScope()
	var err error
	benchMapIntFloatIgo, err = scope.Compile(`
		m := map[int]float64{1: 1.5, 2: 3.0}
		sum := 0.0
		for i := 0; i < 100; i++ {
			v, ok := m[1]
			if ok {
				sum += v
			}
			m[2] = sum
		}
	`)
	if err != nil {
		panic(err)
	}
	benchMapIntFloatIgoVars = benchMapIntFloatIgo.NewState()
}

func initMapInt64FloatIgo() {
	scope := plan.NewScope()
	var err error
	benchMapInt64FloatIgo, err = scope.Compile(`
		m := make(map[int64]float64, 2)
		var k1 int64 = 1
		var k2 int64 = 2
		m[k1] = 1.5
		m[k2] = 3.0
		sum := 0.0
		for i := 0; i < 100; i++ {
			v, ok := m[k1]
			if ok {
				sum += v
			}
			m[k2] = sum
		}
	`)
	if err != nil {
		panic(err)
	}
	benchMapInt64FloatIgoVars = benchMapInt64FloatIgo.NewState()
}

func initMapUint64FloatIgo() {
	scope := plan.NewScope()
	var err error
	benchMapUint64FloatIgo, err = scope.Compile(`
		m := make(map[uint64]float64, 2)
		var k1 uint64 = 1
		var k2 uint64 = 2
		m[k1] = 1.5
		m[k2] = 3.0
		sum := 0.0
		for i := 0; i < 100; i++ {
			v, ok := m[k1]
			if ok {
				sum += v
			}
			m[k2] = sum
		}
	`)
	if err != nil {
		panic(err)
	}
	benchMapUint64FloatIgoVars = benchMapUint64FloatIgo.NewState()
}

func initMapRangeIgo() {
	scope := plan.NewScope()
	var err error
	benchMapRangeIgo, err = scope.Compile(`
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		sum := 0
		for k, v := range m {
			sum += len(k) + v
		}
	`)
	if err != nil {
		panic(err)
	}
	benchMapRangeIgoVars = benchMapRangeIgo.NewState()
}

func initVariadicStringVoidIgo() {
	scope := plan.NewScope()
	scope.RegisterFunc("sinkV", func(prefix string, v ...interface{}) {
		_ = prefix
		_ = v
	})
	var err error
	benchVariadicStringVoidIgo, err = scope.Compile(`
		for i := 0; i < 100; i++ {
			sinkV("x", i, i + 1, i + 2)
		}
	`)
	if err != nil {
		panic(err)
	}
	benchVariadicStringVoidIgoVars = benchVariadicStringVoidIgo.NewState()
}

func initVariadicInterfaceIgo() {
	scope := plan.NewScope()
	scope.DefineVariable("out", reflect.TypeOf((*interface{})(nil)).Elem())
	scope.RegisterFunc("mixV", func(base interface{}, v ...interface{}) interface{} {
		if len(v) == 0 {
			return base
		}
		return v[0]
	})
	var err error
	benchVariadicInterfaceIgo, err = scope.Compile(`
		for i := 0; i < 100; i++ {
			out = mixV(i, i + 1, i + 2, i + 3)
		}
	`)
	if err != nil {
		panic(err)
	}
	benchVariadicInterfaceIgoVars = benchVariadicInterfaceIgo.NewState()
}

func initComplexStateIgo() {
	scope := plan.NewScope()
	stateType := reflect.TypeOf(complexState{})
	_, _ = scope.DefineVariable("st", stateType)
	var err error
	benchComplexStateIgo, err = scope.Compile(`
		st.Sum = 0
		st.Count = 0
		st.Mode = 0
		st.Output = 0
		for i := 0; i < 200; i++ {
			if (i & 1) == 0 {
				st.Sum += i
			} else {
				st.Sum -= i / 2
			}
			switch i % 5 {
			case 0:
				st.Mode = 1
			case 1:
				st.Mode = 2
			case 2:
				st.Mode = 3
			case 3:
				st.Mode = 4
			default:
				st.Mode = 5
			}
			if st.Flags["hot"] {
				st.Output += st.Scores[i % 4]
			}
			if st.Sum > 3000 {
				break
			}
			if i % 7 == 0 {
				continue
			}
			st.Count++
		}
	`)
	if err != nil {
		panic(err)
	}
	benchComplexStateIgoVars = benchComplexStateIgo.NewState()
	err = benchComplexStateIgoVars.SetValue("st", complexState{
		Flags:  map[string]bool{"hot": true},
		Scores: map[int]int{0: 3, 1: 5, 2: 7, 3: 11},
	})
	if err != nil {
		panic(err)
	}
}

func initComplexStateGuardIgo() {
	scope := plan.NewScope()
	stateType := reflect.TypeOf(complexState{})
	_, _ = scope.DefineVariable("st", stateType)
	var err error
	benchComplexStateGuardIgo, err = scope.Compile(`
		st.Sum = 0
		st.Count = 0
		st.Mode = 0
		st.Output = 0
		for i := 0; i < 200; i++ {
			if (i & 1) == 0 {
				st.Sum += i
			} else {
				st.Sum -= i / 2
			}
			switch i % 5 {
			case 0:
				st.Mode = 1
			case 1:
				st.Mode = 2
			case 2:
				st.Mode = 3
			case 3:
				st.Mode = 4
			default:
				st.Mode = 5
			}
			if st.Flags["hot"] {
				st.Output += st.Scores[i % 4]
			}
			if st.Sum >= 3000 {
				break
			}
			if i % 7 != 0 {
				st.Count++
			}
		}
	`)
	if err != nil {
		panic(err)
	}
	benchComplexStateGuardIgoVars = benchComplexStateGuardIgo.NewState()
	err = benchComplexStateGuardIgoVars.SetValue("st", complexState{
		Flags:  map[string]bool{"hot": true},
		Scores: map[int]int{0: 3, 1: 5, 2: 7, 3: 11},
	})
	if err != nil {
		panic(err)
	}
}

func initMergeSortIgo() {
	scope := plan.NewScope()
	intSliceType := reflect.TypeOf([]int{})
	_, _ = scope.DefineVariable("src", intSliceType)
	_, _ = scope.DefineVariable("data", intSliceType)
	_, _ = scope.DefineVariable("tmp", intSliceType)
	_, _ = scope.DefineVariable("mergeOK", reflect.TypeOf(0))
	_, _ = scope.DefineVariable("mergeSum", reflect.TypeOf(0))
	var err error
	benchMergeSortIgo, err = scope.Compile(`
		n := len(src)
		i := 0
		for i < n {
			data[i] = src[i]
			i++
		}
		width := 1
		left := 0
		mid := 0
		right := 0
		j := 0
		k := 0
		p := 0
		prev := 0
		for width < n {
			left = 0
			for left < n {
				mid = left + width
				if mid > n {
					mid = n
				}
				right = mid + width
				if right > n {
					right = n
				}
				i = left
				j = mid
				k = left
				for i < mid && j < right {
					if data[i] <= data[j] {
						tmp[k] = data[i]
						i++
					} else {
						tmp[k] = data[j]
						j++
					}
					k++
				}
				for i < mid {
					tmp[k] = data[i]
					i++
					k++
				}
				for j < right {
					tmp[k] = data[j]
					j++
					k++
				}
				p = left
				for p < right {
					data[p] = tmp[p]
					p++
				}
				left = right
			}
			width += width
		}
		mergeOK = 1
		mergeSum = 0
		if n > 0 {
			prev = data[0]
		}
		i = 0
		for i < n {
			mergeSum += data[i]
			if i > 0 && prev > data[i] {
				mergeOK = 0
			}
			prev = data[i]
			i++
		}
	`)
	if err != nil {
		panic(err)
	}
	benchMergeSortIgoVars = benchMergeSortIgo.NewState()
	data := make([]int, len(mergeSortSrc))
	tmp := make([]int, len(mergeSortSrc))
	if err = benchMergeSortIgoVars.SetValue("src", append([]int(nil), mergeSortSrc...)); err != nil {
		panic(err)
	}
	if err = benchMergeSortIgoVars.SetValue("data", data); err != nil {
		panic(err)
	}
	if err = benchMergeSortIgoVars.SetValue("tmp", tmp); err != nil {
		panic(err)
	}
}

func initConditionalCallIgo() {
	scope := plan.NewScope()
	scope.RegisterFunc("sinkInt", func(v int) {
		benchConditionalCallSink += v
	})
	var err error
	benchConditionalCallIgo, err = scope.Compile(`
		for i := 0; i < 100; i++ {
			if i % 3 != 0 {
				sinkInt(i)
			}
		}
	`)
	if err != nil {
		panic(err)
	}
	benchConditionalCallIgoVars = benchConditionalCallIgo.NewState()
}

func initRowFilterIgo() {
	scope := plan.NewScope()
	_ = scope.RegisterType(reflect.TypeOf(filterRow{}))
	_, _ = scope.DefineVariable("rows", reflect.TypeOf([]*filterRow{}))
	_, _ = scope.DefineVariable("needle", reflect.TypeOf(""))
	_, _ = scope.DefineVariable("minScore", reflect.TypeOf(0))
	var err error
	benchRowFilterIgo, err = scope.Compile(`
		filtered := make([]*filterRow, 0)
		needleLower := strings.ToLower(needle)
		for _, row := range rows {
			if row == nil {
				continue
			}
			if row.Name == "" {
				continue
			}
			if contains(strings.ToLower(row.Name), needleLower) && row.Score >= minScore {
				filtered = append(filtered, row)
			}
		}
		filterCount := len(filtered)
		filterSum := 0
		for _, row := range filtered {
			filterSum += row.Score
		}
	`)
	if err != nil {
		panic(err)
	}
	benchRowFilterIgoVars = benchRowFilterIgo.NewState()
	benchRowFilterData = []*filterRow{
		{Name: "alpha", Score: 10},
		{Name: "beta", Score: 2},
		{Name: "galaxy", Score: 8},
		nil,
		{Name: "", Score: 100},
		{Name: "ALMANAC", Score: 7},
		{Name: "zeta", Score: 1},
		{Name: "palace", Score: 6},
	}
	if err = benchRowFilterIgoVars.SetValue("rows", benchRowFilterData); err != nil {
		panic(err)
	}
	if err = benchRowFilterIgoVars.SetString("needle", "al"); err != nil {
		panic(err)
	}
	if err = benchRowFilterIgoVars.SetInt("minScore", 6); err != nil {
		panic(err)
	}
}

func initInlineHelperAdapterIgo() {
	scope := plan.NewScope()
	_ = scope.RegisterType(reflect.TypeOf(adapterInRow{}))
	_ = scope.RegisterType(reflect.TypeOf(adapterOutRow{}))
	_, _ = scope.DefineVariable("rows", reflect.TypeOf([]*adapterInRow{}))
	_, _ = scope.DefineVariable("needle", reflect.TypeOf(""))
	_, _ = scope.DefineVariable("minScore", reflect.TypeOf(0))
	_, _ = scope.DefineVariable("outRows", reflect.TypeOf([]adapterOutRow{}))
	_, _ = scope.DefineVariable("outSum", reflect.TypeOf(0))
	var err error
	benchInlineHelperAdapterIgo, err = scope.Compile(`
		adapt := func(row *adapterInRow, needle string, minScore int) (adapterOutRow, bool) {
			nameLower := strings.ToLower(row.Name)
			needleLower := strings.ToLower(needle)
			if contains(nameLower, needleLower) && row.Score >= minScore {
				return adapterOutRow{Label: nameLower, Value: row.Score + 1}, true
			}
			return adapterOutRow{}, false
		}

		outRows = make([]adapterOutRow, 0)
		outSum = 0
		for _, row := range rows {
			item, ok := adapt(row, needle, minScore)
			if !ok {
				continue
			}
			outRows = append(outRows, item)
			outSum += item.Value
		}
	`)
	if err != nil {
		panic(err)
	}
	benchInlineHelperAdapterIgoVars = benchInlineHelperAdapterIgo.NewState()
	benchInlineHelperAdapterData = []*adapterInRow{
		{Name: "ALPHA", Score: 10},
		{Name: "beta", Score: 3},
		{Name: "galaxy", Score: 8},
		{Name: "palace", Score: 9},
		{Name: "zeta", Score: 1},
	}
	if err = benchInlineHelperAdapterIgoVars.SetValue("rows", benchInlineHelperAdapterData); err != nil {
		panic(err)
	}
	if err = benchInlineHelperAdapterIgoVars.SetString("needle", "al"); err != nil {
		panic(err)
	}
	if err = benchInlineHelperAdapterIgoVars.SetInt("minScore", 8); err != nil {
		panic(err)
	}
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

func Benchmark_Switch_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchSwitchIgo.Exec(benchSwitchIgoVars)
	}
}

func Benchmark_Defer_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchDeferIgo.Exec(benchDeferIgoVars)
	}
}

func Benchmark_Map_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchMapIgo.Exec(benchMapIgoVars)
	}
}

func Benchmark_MapStringBool_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchMapStringBoolIgo.Exec(benchMapStringBoolIgoVars)
	}
}

func Benchmark_MapIntFloat64_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchMapIntFloatIgo.Exec(benchMapIntFloatIgoVars)
	}
}

func Benchmark_MapInt64Float64_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchMapInt64FloatIgo.Exec(benchMapInt64FloatIgoVars)
	}
}

func Benchmark_MapUint64Float64_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchMapUint64FloatIgo.Exec(benchMapUint64FloatIgoVars)
	}
}

func Benchmark_MapRange_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchMapRangeIgo.Exec(benchMapRangeIgoVars)
	}
}

func Benchmark_ComplexState_Native(b *testing.B) {
	state := complexState{
		Flags:  map[string]bool{"hot": true},
		Scores: map[int]int{0: 3, 1: 5, 2: 7, 3: 11},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		complexStateNative(&state)
	}
}

func Benchmark_ComplexState_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchComplexStateIgo.Exec(benchComplexStateIgoVars)
	}
}

func Benchmark_ComplexState_Guard_Native(b *testing.B) {
	state := complexState{
		Flags:  map[string]bool{"hot": true},
		Scores: map[int]int{0: 3, 1: 5, 2: 7, 3: 11},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		complexStateNativeGuard(&state)
	}
}

func Benchmark_ComplexState_Guard_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchComplexStateGuardIgo.Exec(benchComplexStateGuardIgoVars)
	}
}

func Test_MergeSort_Igo(t *testing.T) {
	benchMergeSortIgo.Exec(benchMergeSortIgoVars)
	ok, err := benchMergeSortIgoVars.Int("mergeOK")
	assert.NoError(t, err)
	assert.Equal(t, 1, ok)
	sum, err := benchMergeSortIgoVars.Int("mergeSum")
	assert.NoError(t, err)
	assert.Equal(t, mergeSortExpectedSum, sum)
}

func Benchmark_MergeSort_Native(b *testing.B) {
	src := append([]int(nil), mergeSortSrc...)
	data := make([]int, len(src))
	tmp := make([]int, len(src))
	lastOK := 0
	lastSum := 0
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lastOK, lastSum = mergeSortNativeRun(src, data, tmp)
	}
	assert.Equal(b, 1, lastOK)
	assert.Equal(b, mergeSortExpectedSum, lastSum)
}

func Benchmark_MergeSort_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchMergeSortIgo.Exec(benchMergeSortIgoVars)
	}
	ok, err := benchMergeSortIgoVars.Int("mergeOK")
	assert.NoError(b, err)
	assert.Equal(b, 1, ok)
	sum, err := benchMergeSortIgoVars.Int("mergeSum")
	assert.NoError(b, err)
	assert.Equal(b, mergeSortExpectedSum, sum)
}

func Benchmark_VariadicStringVoid_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchVariadicStringVoidIgo.Exec(benchVariadicStringVoidIgoVars)
	}
}

func Benchmark_VariadicInterfaceReturn_Igo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchVariadicInterfaceIgo.Exec(benchVariadicInterfaceIgoVars)
	}
}

func Benchmark_ConditionalCall_Native(b *testing.B) {
	sum := 0
	sink := func(v int) {
		sum += v
	}
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			if i%3 != 0 {
				sink(i)
			}
		}
	}
	assert.NotZero(b, sum)
}

func Benchmark_ConditionalCall_Igo(b *testing.B) {
	benchConditionalCallSink = 0
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchConditionalCallIgo.Exec(benchConditionalCallIgoVars)
	}
	assert.NotZero(b, benchConditionalCallSink)
}

func Benchmark_RowFilter_Native(b *testing.B) {
	rows := benchRowFilterData
	needle := "al"
	minScore := 6
	filterCount := 0
	filterSum := 0
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		filtered := make([]*filterRow, 0)
		needleLower := strings.ToLower(needle)
		for _, row := range rows {
			if row == nil || row.Name == "" {
				continue
			}
			if strings.Contains(strings.ToLower(row.Name), needleLower) && row.Score >= minScore {
				filtered = append(filtered, row)
			}
		}
		filterCount = len(filtered)
		filterSum = 0
		for _, row := range filtered {
			filterSum += row.Score
		}
	}
	assert.Equal(b, 4, filterCount)
	assert.Equal(b, 31, filterSum)
}

func Benchmark_RowFilter_Igo(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		benchRowFilterIgo.Exec(benchRowFilterIgoVars)
	}
	filterCount, err := benchRowFilterIgoVars.Int("filterCount")
	assert.NoError(b, err)
	assert.Equal(b, 4, filterCount)
	filterSum, err := benchRowFilterIgoVars.Int("filterSum")
	assert.NoError(b, err)
	assert.Equal(b, 31, filterSum)
}

func Benchmark_InlineHelperAdapter_Native(b *testing.B) {
	rows := benchInlineHelperAdapterData
	needle := "al"
	minScore := 8
	outRows := make([]adapterOutRow, 0, len(rows))
	outSum := 0
	adapt := func(row *adapterInRow, needle string, minScore int) (adapterOutRow, bool) {
		nameLower := strings.ToLower(row.Name)
		needleLower := strings.ToLower(needle)
		if strings.Contains(nameLower, needleLower) && row.Score >= minScore {
			return adapterOutRow{Label: nameLower, Value: row.Score + 1}, true
		}
		return adapterOutRow{}, false
	}
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		outRows = outRows[:0]
		outSum = 0
		for _, row := range rows {
			item, ok := adapt(row, needle, minScore)
			if !ok {
				continue
			}
			outRows = append(outRows, item)
			outSum += item.Value
		}
	}
	assert.Equal(b, 3, len(outRows))
	assert.Equal(b, 30, outSum)
}

func Benchmark_InlineHelperAdapter_Igo(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		benchInlineHelperAdapterIgo.Exec(benchInlineHelperAdapterIgoVars)
	}
	outSum, err := benchInlineHelperAdapterIgoVars.Int("outSum")
	assert.NoError(b, err)
	assert.Equal(b, 30, outSum)
	value, err := benchInlineHelperAdapterIgoVars.Value("outRows")
	assert.NoError(b, err)
	rows, ok := value.([]adapterOutRow)
	assert.True(b, ok)
	assert.Equal(b, 3, len(rows))
}
