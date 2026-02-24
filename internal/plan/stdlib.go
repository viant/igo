package plan

import (
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/traefik/yaegi/stdlib"
)

var (
	stdlibIndexOnce sync.Once
	stdlibPkgKey    map[string]string
)

var stdlibPkgOverrides = map[string]string{
	"rand":     "math/rand/rand",
	"template": "text/template/template",
}

func initStdlibIndex() {
	stdlibPkgKey = map[string]string{}
	keys := make([]string, 0, len(stdlib.Symbols))
	for key := range stdlib.Symbols {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		pos := strings.LastIndex(key, "/")
		if pos == -1 || pos == len(key)-1 {
			continue
		}
		pkg := key[pos+1:]
		if _, exists := stdlibPkgKey[pkg]; exists {
			continue
		}
		stdlibPkgKey[pkg] = key
	}
	for pkg, key := range stdlibPkgOverrides {
		if _, ok := stdlib.Symbols[key]; ok {
			stdlibPkgKey[pkg] = key
		}
	}
}

func splitPkgSymbol(id string) (string, string, bool) {
	pos := strings.LastIndex(id, ".")
	if pos <= 0 || pos == len(id)-1 {
		return "", "", false
	}
	return id[:pos], id[pos+1:], true
}

func stdlibSymbol(id string) (reflect.Value, bool) {
	pkg, symbol, ok := splitPkgSymbol(id)
	if !ok {
		return reflect.Value{}, false
	}
	stdlibIndexOnce.Do(initStdlibIndex)
	key, ok := stdlibPkgKey[pkg]
	if !ok {
		return reflect.Value{}, false
	}
	pkgSymbols, ok := stdlib.Symbols[key]
	if !ok {
		return reflect.Value{}, false
	}
	value, ok := pkgSymbols[symbol]
	return value, ok
}

func (s *Scope) stdlibFunction(id string) (interface{}, bool) {
	value, ok := stdlibSymbol(id)
	if !ok || !value.IsValid() || value.Kind() != reflect.Func {
		return nil, false
	}
	fn := value.Interface()
	s.funcs[id] = fn
	return fn, true
}

func (s *Scope) stdlibType(name string) (reflect.Type, bool) {
	value, ok := stdlibSymbol(name)
	if !ok || !value.IsValid() {
		return nil, false
	}
	if value.Kind() == reflect.Ptr {
		t := value.Type().Elem()
		s.types[name] = t
		return t, true
	}
	return nil, false
}
