package client

import (
	"errors"
	"github.com/enorith/container"
	"net/http"
	"reflect"
	"time"
)

const (
	stateIdle = iota
	statePending
	stateSuccess
	stateFailed
)

type ErrorHandler func(error)

type Result struct {
	Response *http.Response
	Err      error
}

type PendingRequest struct {
	state      int
	result     *Result
	resultChan chan *Result
	handler    Handler
	request    *http.Request
	timeout    time.Duration
	container  *container.Container
}

func (pr *PendingRequest) do() *PendingRequest {
	if pr.state != stateIdle {
		return pr
	}
	pr.state = statePending

	go func() {
		defer func() {
			if x := recover(); x != nil {
				var err error
				if e, ok := x.(error); ok {
					err = e
				}
				if e, ok := x.(string); ok {
					err = errors.New(e)
				}
				pr.resultChan <- &Result{
					Response: nil,
					Err:      err,
				}
				return
			}
		}()

		resp, e := pr.handler(pr.request)
		pr.resultChan <- &Result{
			Response: resp,
			Err:      e,
		}
	}()

	return pr
}

func (pr *PendingRequest) Then(cb interface{}) error {
	result := pr.Wait()
	if result.Err == nil {
		t := reflect.TypeOf(cb)
		if t.Kind() == reflect.Func {
			_, e := pr.container.Invoke(cb)
			if e != nil {
				return e
			}
		}
		if t.Kind() == reflect.Ptr {
			unmarshalResponse(result.Response, cb, reflect.Value{})
		}
	}

	return nil
}

func (pr *PendingRequest) Catch(cb ErrorHandler) {
	result := pr.Wait()
	if result.Err != nil {
		cb(result.Err)
	}
}

func (pr *PendingRequest) Wait() *Result {
	if pr.result != nil {
		return pr.result
	}

	pr.result = <-pr.resultChan
	pr.container.WithInjector(&responseInjector{resp: pr.result.Response})

	return pr.result
}

func newPendingRequest(req *http.Request, handler Handler) *PendingRequest {
	return &PendingRequest{
		handler:    handler,
		request:    req,
		resultChan: make(chan *Result, 1),
		container: container.New(),
	}
}
