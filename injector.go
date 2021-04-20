package feather

import (
	"reflect"
	"strings"

	"github.com/enorith/supports/reflection"
)

var typeResult reflect.Type

type responseInjector struct {
	result *Result
}

func (r *responseInjector) Injection(abs interface{}, last reflect.Value) (reflect.Value, error) {
	if typeResult == reflection.StructType(abs) {
		return reflect.ValueOf(r.result), nil
	}

	unmarshalResponse(r.result, abs, last)

	return last, nil
}

func (r *responseInjector) When(abs interface{}) bool {
	return typeResult == reflection.StructType(abs) || reflection.SubStructOf(abs, typeResult) > -1
}

func unmarshalResponse(result *Result, v interface{}, last reflect.Value) {
	i := reflection.SubStructOf(v, typeResult)
	if !last.IsValid() {
		last = reflection.ValueOf(v)
	}

	if i > -1 {
		last.Elem().Field(i).Set(reflect.ValueOf(result))
		if strings.Contains(result.Response.Header.Get("Content-Type"), "json") {
			result.Unmarshal(last.Interface())
		}
	}

}

func init() {
	typeResult = reflect.TypeOf(Result{})
}
