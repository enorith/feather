package feather_test

import (
	"net/http"
	"testing"

	"github.com/enorith/feather"
)

func TestClientRequest(t *testing.T) {
	pr, e := feather.Get("https://www.baidu.com")

	if e != nil {
		t.Fatalf("request error %v", e)
	}
	result := pr.Wait()
	if result.Err != nil {
		t.Fatalf("request error %v", result.Err)
	}
}

func TestClientRequestProxy(t *testing.T) {
	c := feather.NewClient(feather.Options{ProxyURL: "http://127.0.0.1:7890"})

	pr, e := c.Get("https://www.google.com")
	if e != nil {
		t.Fatalf("request error %v", e)
	}

	result := pr.Wait()
	if result.Err != nil {
		t.Fatalf("request error %v", result.Err)
	}

}

type Result struct {
	*http.Response
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func TestClientUnmarshalBody(t *testing.T) {
	c := feather.NewClient(feather.Options{ProxyURL: feather.NoneProxy, BaseURI: "https://run.mocky.io/v3/"})

	pr, e := c.Get("58744bd4-a1ec-4555-9078-1be561b07043")
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

func TestInterceptor1(t *testing.T) {
	c := feather.NewClient(feather.Options{BaseURI: "https://run.mocky.io/v3/", ErrorUnsuccess: true})

	pr, e := c.Get("ed68fdb5-e9e9-4846-bb0f-f208e6820039")
	if e != nil {
		t.Fatalf("request error %v", e)
	}

	pr.Then(func(r Result) {
		t.Logf("response: %v", r)
	})

	pr.Catch(func(err error) {
		if ue, ok := err.(feather.UnsuccessError); ok {
			t.Logf("error: %v[%d]", err, ue.Response.StatusCode)
		}
	})
}

func TestInterceptor2(t *testing.T) {
	c := feather.NewClient(feather.Options{BaseURI: "https://run.mocky.io/v3/", ErrorUnsuccess: true})

	c.Interceptor(func(r *http.Request, next feather.Handler) *feather.Result {
		t.Logf("request [%s](%s)", r.Method, r.URL)
		return next(r)
	})

	pr, e := c.Get("58744bd4-a1ec-4555-9078-1be561b07043")
	if e != nil {
		t.Fatalf("request error %v", e)
	}
	var r Result
	pr.Then(&r)

	t.Logf("response: %v", r)

	pr.Catch(func(err error) {
		if ue, ok := err.(feather.UnsuccessError); ok {
			t.Logf("error: %v[%d]", err, ue.Response.StatusCode)
		} else {
			t.Logf("error: %v", err)
		}
	})
}
