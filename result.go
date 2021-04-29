package feather

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type Result struct {
	*http.Response
	Err         error
	content     []byte
	contentRead bool
}

func (r *Result) Content() []byte {
	if r.contentRead {
		return r.content
	}
	r.contentRead = true
	if r.Response == nil {
		return nil
	}

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(b))
	r.content = b

	return b
}

func (r *Result) ContentString() string {
	return string(r.Content())
}

func (r *Result) Unmarshal(v interface{}) error {
	if r.Response == nil {
		return errors.New("unmarshal nil response")
	}

	return jsoniter.Unmarshal(r.Content(), v)
}
