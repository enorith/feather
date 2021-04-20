package feather_test

import (
	"fmt"
	"net/http"
	"os"
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
	*feather.Result
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
	c := feather.NewClient(feather.Options{BaseURI: "https://run.mocky.io/v3/", HttpErrors: true})

	pr, e := c.Get("ed68fdb5-e9e9-4846-bb0f-f208e6820039")
	if e != nil {
		t.Fatalf("request error %v", e)
	}

	pr.Then(func(r Result) {
		t.Logf("response: %v", r)
	})

	pr.Catch(func(err error) {
		if ue, ok := err.(feather.HttpError); ok {
			t.Logf("error: %v[%d]", err, ue.Response.StatusCode)
		}
	})
}

func TestInterceptor2(t *testing.T) {
	c := feather.NewClient(feather.Options{BaseURI: "https://run.mocky.io/v3/", HttpErrors: true})

	c.Interceptor(requestLogger(t))

	pr, e := c.Get("58744bd4-a1ec-4555-9078-1be561b07043")
	if e != nil {
		t.Fatalf("request error %v", e)
	}
	pr.Then(func(r *Result) {
		t.Logf("response result: %v", r)
	})

	pr.Catch(func(err error) {
		if ue, ok := err.(feather.HttpError); ok {
			t.Logf("error: %v[%d]", err, ue.Response.StatusCode)
		} else {
			t.Logf("error: %v", err)
		}
	})
}

func TestDownloadFile(t *testing.T) {
	file, _ := os.OpenFile("test_download.exe", os.O_CREATE|os.O_WRONLY, 0775)
	c := feather.NewClient(feather.Options{ProxyURL: "http://127.0.0.1:7890"})
	c.Interceptor(requestLogger(t))

	resp, e := c.Get("https://github.com/Fndroid/clash_for_windows_pkg/releases/download/0.15.3/Clash.for.Windows.Setup.0.15.3.exe", feather.RequestOptions{
		Sink: file,
		OnProgress: func(now, total int64) {
			fmt.Printf("\rdownloading: %.2f%%", float64(now)/float64(total)*100)
		},
	})

	if e != nil {
		t.Fatal(e)
	}

	resp.Wait()
}

func requestLogger(t *testing.T) feather.PipeFunc {
	return func(r *http.Request, next feather.Handler) *feather.Result {
		t.Logf("request [%s](%s)", r.Method, r.URL)
		return next(r)
	}
}
