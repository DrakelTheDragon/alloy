package result

type Result[T any] struct {
	Ok  T
	Err error
}

func Ok[T any](ok T) Result[T]              { return Result[T]{Ok: ok} }
func Err(err error) Result[any]             { return Result[any]{Err: err} }
func Wrap[T any](ok T, err error) Result[T] { return Result[T]{Ok: ok, Err: err} }
func Unwrap[T any](r Result[T]) (T, error)  { return r.Ok, r.Err }

func (r Result[T]) IsOk() bool         { return r.Err == nil }
func (r Result[T]) IsErr() bool        { return r.Err != nil }
func (r Result[T]) Unwrap() (T, error) { return r.Ok, r.Err }
