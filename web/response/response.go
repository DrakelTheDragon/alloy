package response

import (
	"encoding/json"
	"io"
	"net/http"
)

type (
	Response interface {
		Code() int
		Body() any
		Header() http.Header
	}

	response struct {
		code   int
		header http.Header
		body   any
	}

	Success struct{ *response }
	Failure struct{ *response }
)

var _, _, _ Response = (*response)(nil), (*Success)(nil), (*Failure)(nil)

func NewSuccess(code int, data any) Success { return Success{&response{code: code, body: data}} }
func NewFailure(code int, data any) Failure { return Failure{&response{code: code, body: data}} }

func OK(data any) Success      { return NewSuccess(http.StatusOK, data) }
func Created(data any) Success { return NewSuccess(http.StatusCreated, data) }

func NotFound() Failure {
	const msg string = "The requested resource could not be found"
	return NewFailure(http.StatusNotFound, NewError(msg).Payload())
}

func BadRequest(msg string) Failure {
	return NewFailure(http.StatusBadRequest, NewError(msg).Payload())
}

func Conflict(msg string) Failure {
	return NewFailure(http.StatusConflict, NewError(msg).Payload())
}

func UnprocessableEntity(validations Validations) Failure {
	const msg string = "The request could not be processed due to validation failures"
	return NewFailure(http.StatusUnprocessableEntity, NewError(msg).WithValidations(validations).Payload())
}

func Unauthorized() Failure {
	const msg string = "The requested resource requires authentication"
	return NewFailure(http.StatusUnauthorized, NewError(msg))
}

func Forbidden() Failure {
	const msg string = "The requested resource requires authorization"
	return NewFailure(http.StatusForbidden, NewError(msg))
}

func InternalServerError(err error) Failure {
	const msg string = "The server encountered a problem and could not process your request"
	return NewFailure(http.StatusInternalServerError, NewError(msg).WithSource(err).Payload())
}

func (r response) Code() int                              { return r.code }
func (r response) Header() http.Header                    { return r.header }
func (r response) Body() any                              { return r.body }
func (r response) MarshalJSON() ([]byte, error)           { return json.Marshal(r.body) }
func (r response) JSONify(w io.Writer) error              { return JSONify(w, &r) }
func (r response) WithHeader(header http.Header) response { r.header = header; return r }

func (f Failure) withBody(body any) Failure { f.body = body; return f }
