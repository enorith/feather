package feather

import (
	"net/http"
)

func httpErrorInterceptor(r *http.Request, next Handler) *Result {
	result := next(r)

	if result.Err == nil && result.Response.StatusCode != 200 {
		result.Err = HttpError{result}
		return result
	}

	return result
}
