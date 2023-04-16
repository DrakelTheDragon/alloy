package web

import (
	"context"
	"io"
	"net/http"

	"github.com/drakelthedragon/alloy/process"
	"github.com/drakelthedragon/alloy/web/request"
	"github.com/drakelthedragon/alloy/web/response"
)

type (
	Response interface {
		Code() int
		Body() any
		Header() http.Header
	}

	Result[T any]             interface{ Unwrap() (T, error) }
	RequestMapper[I any]      interface{ MapRequest(*http.Request, *I) error }
	ResponseMapper[O any]     interface{ MapResponse(Result[O]) Response }
	RequestMapperFunc[I any]  func(*http.Request, *I) error
	ResponseMapperFunc[O any] func(Result[O]) Response
	RequestProcessor          interface{ ProcessRequest(*http.Request) Response }

	requestProcessor[I, O any] struct {
		reqm RequestMapper[I]
		resm ResponseMapper[O]
		proc *process.Processor[I, O]
	}
)

func NewRequestProcessor[I, O any](reqm RequestMapper[I], resm ResponseMapper[O], fn func(context.Context, I) (O, error)) *requestProcessor[I, O] {
	return &requestProcessor[I, O]{
		reqm: reqm,
		resm: resm,
		proc: process.NewProcessor(fn),
	}
}

func NewHandler(rp RequestProcessor, logf func(error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := response.JSONify(w, rp.ProcessRequest(r)); err != nil {
			logf(err)
		}
	}
}

// WriteJSON writes a JSON response to the io.Writer or http.ResponseWriter.
func WriteJSON(w io.Writer, r Response) error { return response.JSONify(w, r) }
func ReadJSON(r *http.Request, dst any) error { return request.UnJSONify(r, dst) }

func (f RequestMapperFunc[I]) MapRequest(r *http.Request, input *I) error { return f(r, input) }
func (f ResponseMapperFunc[O]) MapResponse(result Result[O]) Response     { return f(result) }

func (rp *requestProcessor[I, O]) ProcessRequest(r *http.Request) Response {
	var input I

	if err := rp.reqm.MapRequest(r, &input); err != nil {
		return response.BadRequest("Could not parse request body")
	}

	result := rp.proc.Process(context.Background(), input)

	return rp.resm.MapResponse(result)
}
