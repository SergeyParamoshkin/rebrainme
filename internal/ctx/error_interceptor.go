package fwctx

import "context"

type ErrorInterceptor struct {
	Err error
}

func WithErrorInterceptor(ctx context.Context) context.Context {
	return context.WithValue(ctx, KeyErrorInterceptor, &ErrorInterceptor{})
}

func RecordError(ctx context.Context, err error) context.Context {
	ei, _ := ctx.Value(KeyErrorInterceptor).(*ErrorInterceptor)
	if ei == nil {
		return nil
	}

	ei.Err = err

	return ctx
}

func ReadError(ctx context.Context) error {
	ei, _ := ctx.Value(KeyErrorInterceptor).(*ErrorInterceptor)
	if ei == nil {
		return nil
	}

	return ei.Err
}
