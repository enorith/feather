package client

import "net/http"

type Handler func(r *http.Request) (*http.Response, error)

type PipeFunc func(r *http.Request, next Handler) (*http.Response, error)

type Pipeline []PipeFunc

func (pl Pipeline) Push(pf PipeFunc) Pipeline {
	return append(pl, pf)
}

func (pl Pipeline) Resolve(r *http.Request, handler Handler) (*http.Response, error) {

	return func(r *http.Request) (*http.Response, error) {
		if pl != nil && len(pl) > 0 {
			next := pl.prepareNext(0, handler)
			return pl[0](r, next)
		}
		return handler(r)
	}(r)
}

func (pl Pipeline) prepareNext(now int, handler Handler) Handler {
	l := len(pl)
	var next Handler
	if now+1 >= l {
		next = handler
	} else {
		next = func(r *http.Request) (*http.Response, error) {
			return pl[now+1](r, pl.prepareNext(now+1, handler))
		}
	}

	return next
}