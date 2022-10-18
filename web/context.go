package web

import (
	"github.com/valyala/fasthttp"
)

const (
	errKey string = "weirdErrorKey"
	ridKey string = "___requestid"
)

func AddError(ctx *fasthttp.RequestCtx, err error) {
	v, _ := ctx.UserValue(errKey).([]error)
	v = append(v, err)

	ctx.SetUserValue(errKey, v)
}

func GetErrors(ctx *fasthttp.RequestCtx) []error {
	if errs, ok := ctx.Value(errKey).([]error); ok {
		return errs
	}

	return nil
}

func GetRequestID(ctx *fasthttp.RequestCtx) (string, bool) {
	v, ok := ctx.UserValue(ridKey).(string)
	if !ok {
		return "", false
	}

	if v == "" {
		return "", false
	}

	return v, true
}

func SetRequestID(ctx *fasthttp.RequestCtx, rid string) {
	ctx.SetUserValue(ridKey, rid)
}
