package web

import (
	"context"
)

type errKey int

const key errKey = 1

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
