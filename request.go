package feather

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/enorith/container"
)

const (
	stateIdle = iota
	statePending
	stateSuccess
	stateFailed
)

type ErrorHandler func(error)

type PendingRequest struct {
	state      int
	result     *Result
	resultChan chan *Result
	handler    Handler
	request    *http.Request
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

		pr.resultChan <- pr.handler(pr.request)
	}()

	return pr
}

func (pr *PendingRequest) Then(cb interface{}) *PendingRequest {
	result := pr.Wait()
	if result.Err == nil {
		t := reflect.TypeOf(cb)
		if t.Kind() == reflect.Func {
			_, e := pr.container.Invoke(cb)
			if e != nil {
				return pr
			}
		}
		if t.Kind() == reflect.Ptr {
			unmarshalResponse(result, cb, reflect.Value{})
		}
	}

	return pr
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
	if pr.result.Err == nil {
		pr.state = stateSuccess
	} else {
		pr.state = stateFailed
	}

	pr.container.WithInjector(&responseInjector{result: pr.result})

	return pr.result
}

func (pr *PendingRequest) IsSuccess() bool {
	pr.Wait()
	return pr.state == stateSuccess
}

func newPendingRequest(req *http.Request, handler Handler) *PendingRequest {
	return &PendingRequest{
		handler:    handler,
		request:    req,
		resultChan: make(chan *Result, 1),
		container:  container.New(),
	}
}
