package feather

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

func init() {
	DefaultClient = NewClient()
}
