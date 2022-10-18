package web

import (
	"context"
)

type ctxKey int

const (
	errKey ctxKey = iota
	ridKey
)

func AddError(ctx context.Context, err error) context.Context {
	v, _ := ctx.Value(errKey).([]error)
	v = append(v, err)
	return context.WithValue(ctx, errKey, v)
}

func GetErrors(ctx context.Context) []error {
	if errs, ok := ctx.Value(errKey).([]error); ok {
		return errs
	}
	return nil
}

func ContextWithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, ridKey, rid)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(ridKey).(string)
	return rid, ok
}
