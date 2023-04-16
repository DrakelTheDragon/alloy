package request

import (
	"io"
)

const defaultMaxBytes int64 = 1_048_576

type request[T any] struct{ m mapper[T] }

type mapper[T any] func(io.Reader, *T) error

type loginInput struct {
	Username string
}

func jsonLoginMapper(r io.Reader, input *loginInput) error {
	var req struct {
		Username string `json:"username"`
	}

	if err := unJSONify(r, &req); err != nil {
		return err
	}

	input.Username = req.Username

	return nil
}

func New[T any](m mapper[T]) *request[T] { return &request[T]{m: m} }

func (r *request[T]) UnJSONify(rr io.Reader, t *T) error { return r.m(rr, t) }
