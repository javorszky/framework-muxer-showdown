package handlers

import (
	"net/http"
	"testing"

	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestErrorHandlers(t *testing.T) {
	const path = "/error"

	testlogger := zerolog.Nop()
	testchan := make(chan error)
	mw := ErrorCatcher(testlogger, testchan)

	tests := []struct {
		name       string
		handler    fasthttp.RequestHandler
		wantStatus int
		wantBody   []byte
	}{
		{
			name:       "app error handler returns app error",
			handler:    ReturnsApplicationError(),
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"message":"app error: application error: some error from someplace"}`),
		},
		{
			name:       "Not found error handler returns not found error",
			handler:    ReturnsNotFoundError(),
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(`{"message":"not found: not found: not found the thing"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.New()
			r.GET(path, mw(tt.handler))

			ctx := new(fasthttp.RequestCtx)
			ctx.Request.SetRequestURI(path)

			r.Handler(ctx)

			assert.Equal(t, tt.wantStatus, ctx.Response.StatusCode())
			assert.Equalf(t, tt.wantBody, ctx.Response.Body(), "expected: %s\ngot: %s\n", tt.wantBody, ctx.Response.Body())

			t.Logf("and now trying without involving a router: just the handlers themself")

			newCtx := new(fasthttp.RequestCtx)
			mw(tt.handler)(newCtx)

			assert.Equal(t, tt.wantStatus, newCtx.Response.StatusCode())
			assert.Equalf(t, tt.wantBody, newCtx.Response.Body(), "expected: %s\ngot: %s\n", tt.wantBody, newCtx.Response.Body())

		})
	}
}
