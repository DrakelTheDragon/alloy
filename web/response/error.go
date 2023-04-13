package response

import "net/http"

type Error struct {
	msg string
	src error
	res *Response
}

type Failures map[string][]string

func NewError(code int, msg string) *Error {
	return &Error{msg: msg, res: &Response{code: code}}
}

func NotFound() *Error {
	const msg string = "The requested resource could not be found"
	return NewError(http.StatusNotFound, msg)
}

func BadRequest(msg string) *Error {
	return NewError(http.StatusBadRequest, msg)
}

func UnprocessableEntity(failures Failures) *Error {
	const msg string = "The request could not be processed due to validation failures"
	return NewError(http.StatusUnprocessableEntity, msg).withResponseBody(
		struct {
			Message  string   `json:"message"`
			Failures Failures `json:"failures"`
		}{Message: msg, Failures: failures},
	)
}

func Unauthorized() *Error {
	const msg string = "The requested resource requires authentication"
	return NewError(http.StatusUnauthorized, msg)
}

func Forbidden() *Error {
	const msg string = "The requested resource requires authorization"
	return NewError(http.StatusForbidden, msg)
}

func InternalServerError(err error) *Error {
	const msg string = "The server encountered a problem and could not process your request"
	return NewError(http.StatusInternalServerError, msg).WithSource(err)
}

func (e *Error) Error() string { return e.msg }
func (e *Error) Unwrap() error { return e.src }

func (e *Error) WithMessage(msg string) *Error { e.msg = msg; return e }
func (e *Error) WithSource(err error) *Error   { e.src = err; return e }

func (e *Error) Response() *Response {
	if e.res.body == nil {
		e.res.body = struct {
			Message string `json:"message"`
		}{Message: e.msg}
	}
	return e.res
}

func (e *Error) JSONify(w http.ResponseWriter) error { return e.Response().JSONify(w) }

func (e *Error) withResponseBody(body any) *Error {
	e.res.body = body
	return e
}

func (f Failures) Add(field, msg string) { f[field] = append(f[field], msg) }
func (f Failures) Set(field, msg string) { f[field] = []string{msg} }
func (f Failures) Has(field string) bool { _, ok := f[field]; return ok }
func (f Failures) Get(field string) []string {
	if f.Has(field) {
		return f[field]
	}
	return nil
}
