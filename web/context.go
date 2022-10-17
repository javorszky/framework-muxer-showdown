package web

import (
	"context"

	"github.com/labstack/echo/v4"
)

const (
	errKey       string = "___errkey"
	requestIDkey string = "___requestIDKey"
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

func AddRequestID(c echo.Context, rid string) {
	c.Set(requestIDkey, rid)
}

func GetRequestID(c echo.Context) string {
	v := c.Get(requestIDkey)
}
