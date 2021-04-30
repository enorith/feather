package feather_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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
	feather.DefaultClient.Interceptor(requestLogger(t))

	resp, e := feather.Download("https://www.baidu.com",
		"test_download.html", feather.RequestOptions{
			OnProgress: func(now, total int64) {
				p := float64(now) / float64(total) * 100
				fmt.Printf("\rdownloading: [%s>%s] %.2f%%", strings.Repeat("=", int(p)), strings.Repeat(" ", 100-int(p)), p)
			},
		})

	if e != nil {
		t.Fatal(e)
	}

	resp.Wait()
}

func TestUpload(t *testing.T) {
	file, _ := os.Open("D:\\workspace\\jpeg.jpg")

	resp, e := feather.NewClient(feather.Options{
		HttpErrors: true,
	}).Interceptor(requestLogger(t)).Post("http://127.0.0.1:8000/upload/image", feather.RequestOptions{
		Upload: feather.NewUploadFile(file, "file", "upload2.jpg"),
		Header: http.Header{
			"Accept": {"application/json"},
		},
	})

	if e != nil {
		t.Fatal(e)
	}

	resp.Then(func(res *feather.Result) {

		t.Log(res.ContentString())
	}).Catch(func(e error) {
		if err, ok := e.(feather.HttpError); ok {
			var em ErrorMessage
			err.Unmarshal(&em)
			t.Log(em)
		}

		t.Log(e.Error())
	})

}

func requestLogger(t *testing.T) feather.PipeFunc {
	return func(r *http.Request, next feather.Handler) *feather.Result {
		t.Logf("request [%s](%s)", r.Method, r.URL)
		return next(r)
	}
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
