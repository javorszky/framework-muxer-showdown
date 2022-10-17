package web

import (
	"context"
)

type errKey int

const (
	key errKey = iota
	requestIDKey
)

func AddError(ctx context.Context, err error) context.Context {
	v, _ := ctx.Value(key).([]error)
	v = append(v, err)
	return context.WithValue(ctx, key, v)
}

func GetErrors(ctx context.Context) []error {
	if errs, ok := ctx.Value(key).([]error); ok {
		return errs
	}
	return nil
}

func ContextWithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestIDKey, rid)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(requestIDKey).(string)
	return rid, ok
}
