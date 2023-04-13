package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	code   int
	header http.Header
	body   any
}

func NewResponse(code int, data any) *Response {
	return &Response{code: code, body: data}
}

func OK(data any) *Response      { return NewResponse(http.StatusOK, data) }
func Created(data any) *Response { return NewResponse(http.StatusCreated, data) }

func (r Response) Code() int                           { return r.code }
func (r Response) Header() http.Header                 { return r.header }
func (r Response) Body() any                           { return r.body }
func (r Response) MarshalJSON() ([]byte, error)        { return json.Marshal(r.body) }
func (r Response) JSONify(w http.ResponseWriter) error { return JSONify(w, &r) }

func (r *Response) WithHeader(header http.Header) *Response { r.header = header; return r }
