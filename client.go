package feather

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	DefaultTimeout = 5 * time.Second
	NoneProxy      = "none"
)

func defaultHandler(opt Options) Handler {
	return func(r *http.Request) *Result {
		resp, err := clientFromOptions(opt).Do(r)
		return &Result{
			Response: resp,
			Err:      err,
		}
	}
}

func clientFromOptions(o Options) *http.Client {
	var rt http.RoundTripper

	if len(o.ProxyURL) > 0 {
		if o.ProxyURL == NoneProxy {
			rt = &http.Transport{}
		} else {
			pu, e := url.Parse(o.ProxyURL)
			if e == nil {
				rt = &http.Transport{
					Proxy: http.ProxyURL(pu),
				}
			}
		}

	}

	return &http.Client{
		Timeout:   o.Timeout,
		Transport: rt,
	}
}

type (
	ProgressHandler func(now, total int64)

	// RequestOptions simplified http request options
	RequestOptions struct {
		// Body is raw request body
		Body io.Reader
		// Handler is request handler
		Handler Handler
		// Json body of request
		Json interface{}
		// Header is request header
		Header http.Header
		// FormParams is request form params
		FormParams url.Values
		// Query is request query
		Query url.Values
		// Method is request method
		Method string
		// OnProgress handle content download progress
		OnProgress ProgressHandler
		// Sink response content to io.Writer, for file download
		Sink io.Writer
	}
	// Client is request clinet
	Client struct {
		p   Pipeline
		opt Options
	}
	// Options is request client options
	Options struct {
		// BaseURI is prefix of request url
		BaseURI string
		// Timeout request timeout
		Timeout time.Duration
		// ProxyURL set proxy url for request
		ProxyURL string
		// HttpErrors trigger error if response is not ok
		HttpErrors bool
	}
)

// Get send GET http request
func (c *Client) Get(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodGet, url, opts...)
}

// Post send POST http request
func (c *Client) Post(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodPost, url, opts...)
}

// Put send PUT http request
func (c *Client) Put(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodPut, url, opts...)
}

// Patch send PATCH http request
func (c *Client) Patch(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodPatch, url, opts...)
}

// Delete send DELETE http request
func (c *Client) Delete(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodDelete, url, opts...)
}

// Head send HEAD http request
func (c *Client) Head(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return c.Request(http.MethodHead, url, opts...)
}

// Request send http request
func (c *Client) Request(method, url string, opts ...RequestOptions) (*PendingRequest, error) {
	if len(c.opt.BaseURI) > 0 {
		url = fmt.Sprintf("%s/%s", strings.TrimSuffix(c.opt.BaseURI, "/"), strings.TrimPrefix(url, "/"))
	}
	o, err := requestOptions(c.opt, opts...)
	if err != nil {
		return nil, err
	}

	req, e := NewRequestFromOptions(method, url, o)
	if e != nil {
		return nil, e
	}

	if o.Sink != nil {
		var total int64

		resp, e := c.Head(url)
		if e == nil {
			resp.Then(func(res *Result) {
				ls := res.Header.Get("Content-Length")
				total, _ = strconv.ParseInt(ls, 10, 64)
			})
		}
		pw := &progressWriter{o.OnProgress, total, 0}
		c.Interceptor(sinkInterceptor(c, url, o.Sink, pw))
	}

	if c.opt.HttpErrors {
		c.Interceptor(httpErrorInterceptor)
	}
	return newPendingRequest(req, func(r *http.Request) *Result {
		return c.p.Resolve(r, o.Handler)
	}).do(), nil
}

// NewRequestFromOptions new http request from RequestOptions
func NewRequestFromOptions(method string, path string, o RequestOptions) (*http.Request, error) {
	req, e := http.NewRequest(method, path, o.Body)
	if e != nil {
		return nil, e
	}
	if o.Method != "" {
		req.Method = o.Method
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

func (c *Client) Config(opt Options) *Client {
	return NewClient(opt) // return new client
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

type HttpError struct {
	*Result
}

func (ue HttpError) Error() string {
	return http.StatusText(ue.Response.StatusCode)
}
