package feather

import (
	"net/http"
	"net/url"
)

func MergeRequestOptions(ros ...RequestOptions) RequestOptions {
	ro := RequestOptions{
		Header: make(http.Header),
		Query:  make(url.Values),
	}

	for _, option := range ros {
		if option.Header != nil {
			for k, v := range option.Header {
				for _, hv := range v {
					ro.Header.Add(k, hv)
				}
			}
		}
		if option.Query != nil {
			for k, v := range option.Query {
				for _, qv := range v {
					ro.Query.Add(k, qv)
				}
			}
		}

		if option.FormParams != nil {
			if ro.FormParams == nil {
				ro.FormParams = make(url.Values)
			}
			for k, v := range option.FormParams {
				for _, qv := range v {
					ro.FormParams.Add(k, qv)
				}
			}
		}
		if option.Body != nil {
			ro.Body = option.Body
		}
		if option.Handler != nil {
			ro.Handler = option.Handler
		}
		if option.Method != "" {
			ro.Method = option.Method
		}
		if option.Upload != nil {
			ro.Upload = option.Upload
		}
		if option.Json != nil {
			ro.Json = option.Json
		}
		if option.Sink != nil {
			ro.Sink = option.Sink
		}
		if option.OnProgress != nil {
			ro.OnProgress = option.OnProgress
		}
	}

	return ro
}
