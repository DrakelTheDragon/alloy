package process

import (
	"context"

	"github.com/drakelthedragon/alloy/result"
)

type Func[I, O any] func(context.Context, I) (O, error)

type Processor[I, O any] struct {
	process Func[I, O]
}

func NewProcessor[I, O any](fn Func[I, O]) *Processor[I, O] {
	return &Processor[I, O]{process: fn}
}

func (p *Processor[I, O]) Process(ctx context.Context, input I) result.Result[O] {
	return result.Wrap(p.process(ctx, input))
}

func (f Func[I, O]) Process(ctx context.Context, input I) result.Result[O] {
	return result.Wrap(f(ctx, input))
}
