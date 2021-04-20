package feather

import (
	"io"
	"net/http"
)

type progressWriter struct {
	onProgress ProgressHandler
	total      int64
	current    int64
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	l := len(p)
	pw.current += int64(l)

	if pw.onProgress != nil {
		pw.onProgress(pw.current, pw.total)
	}

	return l, nil
}

func sinkInterceptor(c *Client, url string, to io.Writer, pw *progressWriter) PipeFunc {
	return func(r *http.Request, next Handler) *Result {
		result := next(r)

		if result.Err == nil {
			defer result.Body.Close()
			io.Copy(to, io.TeeReader(result.Body, pw))
		}

		return result
	}
}

func httpErrorInterceptor(r *http.Request, next Handler) *Result {
	result := next(r)

	if result.Err == nil && result.Response.StatusCode != 200 {
		result.Err = HttpError{result}
		return result
	}

	return result
}
