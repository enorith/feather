package feather

import (
	"errors"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type Result struct {
	Response    *http.Response
	Err         error
	content     []byte
	contentRead bool
}

func (r *Result) Content() []byte {
	if r.contentRead {
		return r.content
	}
	if r.Err != nil {
		return nil
	}

	b, _ := ioutil.ReadAll(r.Response.Body)
	defer r.Response.Body.Close()
	return b
}

func (r *Result) ContentString() string {
	return string(r.Content())
}

func (r *Result) Unmarshal(v interface{}) error {
	if r.Err != nil {
		return errors.New("unmarshal error response")
	}

	return jsoniter.Unmarshal(r.Content(), v)
}
