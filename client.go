package feather

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	DefaultTimeout = 5 * time.Second
)

var defaultClient *http.Client

func defaultHandler(opt Options) Handler {
	return func(r *http.Request) (*http.Response, error) {
		return clientFromOptions(opt).Do(r)
	}
}

func clientFromOptions(o Options) *http.Client {
	if defaultClient == nil {
		defaultClient = &http.Client{
			Timeout: o.Timeout,
		}
	}

	return defaultClient
}

type (
	// RequestOptions simplified http request options
	RequestOptions struct {
		Body       io.Reader
		Handler    Handler
		Json       interface{}
		Header     http.Header
		FormParams url.Values
		Query      url.Values
	}
	Client struct {
		p   Pipeline
		opt Options
	}
	Options struct {
		BaseUri string
		Timeout time.Duration
	}
)

func (c *Client) Get(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodGet, url, opts...)
}

func (c *Client) Post(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodPost, url, opts...)
}

func (c *Client) Put(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodPut, url, opts...)
}

func (c *Client) Patch(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodPatch, url, opts...)
}

func (c *Client) Delete(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodDelete, url, opts...)
}

func (c *Client) Request(method, path string, opts ...RequestOptions) (*PendingRequest, error) {
	if len(c.opt.BaseUri) > 0 {
		path = fmt.Sprintf("%s/%s", strings.TrimSuffix(c.opt.BaseUri, "/"), strings.TrimPrefix(path, "/"))
	}
	o, err := requestOptions(c.opt, opts...)
	if err != nil {
		return nil, err
	}

	req, e := NewRequestFromOptions(method, path, o)
	if e != nil {
		return nil, e
	}

	return newPendingRequest(req, func(r *http.Request) (*http.Response, error) {
		return c.p.Resolve(r, o.Handler)
	}).do(), nil
}

func NewRequestFromOptions(method string, path string, o RequestOptions) (*http.Request, error) {
	req, e := http.NewRequest(method, path, o.Body)
	if e != nil {
		return nil, e
	}

	req.Header = mergeValues(req.Header, o.Header)
	req.Form = mergeValues(req.Header, o.FormParams)
	var values url.Values
	values = mergeValues(values, o.Query)
	req.URL.RawQuery = values.Encode()

	return req, nil
}

func requestOptions(co Options, opts ...RequestOptions) (RequestOptions, error) {
	var o RequestOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Header == nil {
		o.Header = make(http.Header)
	}
	if o.Query == nil {
		o.Query = make(url.Values)
	}

	if o.Handler == nil {
		o.Handler = defaultHandler(co)
	}

	if o.FormParams != nil {
		o.Body = strings.NewReader(o.FormParams.Encode())
		o.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if o.Json != nil {
		b, e := jsoniter.Marshal(o.Json)
		if e != nil {
			return o, e
		}
		o.Body = bytes.NewReader(b)
		o.Header.Set("Content-Type", "application/json")
	}

	return o, nil
}

func (c *Client) Interceptor(pf PipeFunc) *Client {
	c.p = c.p.Push(pf)

	return c
}

func NewClient(opts ...Options) *Client {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &Client{opt: opt}
}

func mergeValues(val, val1 map[string][]string) map[string][]string {
	if val == nil {
		val = make(map[string][]string)
	}
	for k, v := range val1 {
		val[k] = v
	}

	return val
}
