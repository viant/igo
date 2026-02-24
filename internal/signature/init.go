package signature

import (
	"github.com/viant/igo/exec"
	"reflect"
)

var buildIn = []reflect.Type{
	reflect.TypeOf(new(fnV)).Elem(),
	reflect.TypeOf(new(iiiFn)).Elem(),
	reflect.TypeOf(new(iiFn)).Elem(),
	reflect.TypeOf(new(ivFn)).Elem(),
	reflect.TypeOf(new(iisbFn)).Elem(),
	reflect.TypeOf(new(sssbFn)).Elem(),
	reflect.TypeOf(new(ssbFn)).Elem(),
	reflect.TypeOf(new(f64f64f64Fn)).Elem(),
	reflect.TypeOf(new(sssFn)).Elem(),
	reflect.TypeOf(new(svrFn)).Elem(),
	reflect.TypeOf(new(svrs)).Elem(),
	reflect.TypeOf(new(viFn)).Elem(),
	reflect.TypeOf(new(vf64Fn)).Elem(),
	reflect.TypeOf(new(vf32Fn)).Elem(),
	reflect.TypeOf(new(vbFn)).Elem(),
	reflect.TypeOf(new(vsFn)).Elem(),
	reflect.TypeOf(new(vi64Fn)).Elem(),
	reflect.TypeOf(new(iibFn)).Elem(),
	reflect.TypeOf(new(ssiFn)).Elem(),
	reflect.TypeOf(new(i64i64i64Fn)).Elem(),
	reflect.TypeOf(new(i64i64bFn)).Elem(),
	reflect.TypeOf(new(svrieFn)).Elem(),
	reflect.TypeOf(new(vvsvFn)).Elem(),
	reflect.TypeOf(new(vrieFn)).Elem(),
}

func init() {
	for i := range buildIn {
		exec.RegisterSignature(buildIn[i])
	}

}
