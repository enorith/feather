package feather

import "os"

var DefaultClient *Client

func Get(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return DefaultClient.Get(url, opts...)
}

func Post(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return DefaultClient.Post(url, opts...)
}

func Put(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return DefaultClient.Put(url, opts...)
}

func Patch(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return DefaultClient.Patch(url, opts...)
}

func Delete(url string, opts ...RequestOptions) (*PendingRequest, error) {
	return DefaultClient.Delete(url, opts...)
}

func Request(method, url string, opts ...RequestOptions) (*PendingRequest, error) {
	return DefaultClient.Request(method, url, opts...)
}

func Download(url string, filename string, opts ...RequestOptions) (*PendingRequest, error) {
	var opt RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	file, e := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0775)

	if e != nil {
		return nil, e
	}

	opt.Sink = file

	return Get(url, opt)
}

func init() {
	DefaultClient = NewClient()
}
