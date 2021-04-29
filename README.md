# Feather is a golang http client


## Usage

```go get -u github.com/enorith/feather```

```golang
import (
	"github.com/enorith/feather"
)

func FooRequest() {
    client := feather.NewClient(feather.Options{BaseURI: "https://run.mocky.io/v3/", HttpErrors: true})

	pr, e := client.Get("ed68fdb5-e9e9-4846-bb0f-f208e6820039")
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

```

## Interceptor

```golang

func requestLogger(logger *logger.Logger) feather.PipeFunc {
	return func(r *http.Request, next feather.Handler) *feather.Result {
        // print every http request
		logger.Infof("request [%s](%s)", r.Method, r.URL)
		return next(r)
	}
}

client := feather.NewClient()

client.Interceptor(requestLogger(logger.Default()))
```

## Json unmarshal

```golang
type Response struct {
    Code int `json:"code"`
    Message string `json:"message"`
}


client := feather.NewClient()

req, _ := client.Get("http://bar.com/foo.json")

var resp Response
req.Then(&resp)

// or

req.Then(func(resp *Response) {
    fmt.Println(resp.Code)
})
```

## File download

```golang

client := feather.NewClient()

file, _ := os.OpenFile("/tmp/temp.txt", os.O_CREATE|os.O_WRONLY, 0775)

req, _ := client.Get("http://bar.com/foo.json", feather.RequestOptions{
    Sink: file,
    OnProgress: func(now, total int64) {
        p := float64(now) / float64(total) * 100
        fmt.Printf("\rdownloading: [%s>%s] %.2f%%", strings.Repeat("=", int(p)), strings.Repeat(" ", 100-int(p)), p)
    },
})

req.Wait()
```