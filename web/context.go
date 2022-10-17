package web

import (
	"github.com/valyala/fasthttp"
)

const key string = "weirdErrorKey"

func AddError(ctx *fasthttp.RequestCtx, err error) {
	v, _ := ctx.UserValue(key).([]error)
	v = append(v, err)

	ctx.SetUserValue(key, v)
}

func GetErrors(ctx *fasthttp.RequestCtx) []error {
	if errs, ok := ctx.Value(key).([]error); ok {
		return errs
	}

	return nil
}
