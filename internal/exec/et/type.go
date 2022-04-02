package et

import "reflect"

var errType = reflect.TypeOf(new(error)).Elem()
var boolType = reflect.TypeOf(true)
