package pipelm

import "context"

type Handler interface {
	Call(ctx context.Context, values ...Values) (Values, error)
}

type HandlerFunc func(context.Context, ...Values) (Values, error)

func (f HandlerFunc) Call(ctx context.Context, values ...Values) (Values, error) {
	return f(ctx, values...)
}
