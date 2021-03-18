package client

import (
	"github.com/enorith/supports/reflection"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"reflect"
	"strings"
)

var typeResponse reflect.Type

type responseInjector struct {
	resp *http.Response
}

func (r *responseInjector) Injection(abs interface{}, last reflect.Value) (reflect.Value, error) {
	if typeResponse == reflection.StructType(abs) {
		return reflect.ValueOf(r.resp), nil
	}

	unmarshalResponse(r.resp, abs, last)

	return last, nil
}

func (r *responseInjector) When(abs interface{}) bool {
	return typeResponse == reflection.StructType(abs) || reflection.SubStructOf(abs, typeResponse) > -1
}

func unmarshalResponse(resp *http.Response, v interface{}, last reflect.Value) {
	i := reflection.SubStructOf(v, typeResponse)
	if !last.IsValid() {
		last = reflection.ValueOf(v)
	}

	if i > -1 {
		last.Elem().Field(i).Set(reflect.ValueOf(resp))
		if strings.Contains(resp.Header.Get("Content-Type"), "json") {
			_ = jsoniter.NewDecoder(resp.Body).Decode(last.Interface())
		}
	}

}

func init() {
	typeResponse = reflect.TypeOf(http.Response{})
}
