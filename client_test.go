package feather_test

import (
	"testing"

	"github.com/enorith/feather"
)

func TestClientRequest(t *testing.T) {
	c := feather.NewClient()

	pr, e := c.Get("http://www.baidu.com")
	if e != nil {
		t.Fatalf("request error %v", e)
	}

	t.Logf("response: %s", pr.Wait().Content())
}

func TestClientRequestProxy(t *testing.T) {
	c := feather.NewClient(feather.Options{ProxyUrl: "socks5://127.0.0.1:7890"})

	pr, e := c.Get("https://www.google.com")
	if e != nil {
		t.Fatalf("request error %v", e)
	}

	result := pr.Wait()
	if result.Err != nil {
		t.Fatalf("request error %v", result.Err)
	}

	t.Logf("response: %s", result.Content())
}

type Result struct {
	Channels []string `json:"channels"`
}

func TestClientUnmarshalBody(t *testing.T) {
	c := feather.NewClient(feather.Options{ProxyUrl: "none"})

	pr, e := c.Get("http://ubuntu:4161/lookup", feather.RequestOptions{
		Query: map[string][]string{
			"topic": {"topic"},
		}})
	if e != nil {
		t.Fatalf("request error %v", e)
	}

	result := pr.Wait()
	if result.Err != nil {
		t.Fatalf("request error %v", result.Err)
	}

	var re Result
	result.Unmarshal(&re)
	t.Logf("response: %v", re)
}
