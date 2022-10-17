package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandlers(t *testing.T) {
	testlogger := zerolog.Nop()
	testchan := make(chan error)
	const path = "/error"

	type args struct {
		handler    http.Handler
		middleware func(http.Handler) http.Handler
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   []byte
	}{
		{
			name: "app error handler returns app error",
			args: args{
				handler:    ReturnsApplicationError(),
				middleware: ErrorCatcher(testlogger, testchan),
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"message":"app error: application error: some error from someplace"}`),
		},
		{
			name: "Not found error handler returns not found error",
			args: args{
				handler:    ReturnsNotFoundError(),
				middleware: ErrorCatcher(testlogger, testchan),
			},
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(`{"message":"not found: not found: not found the thing"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()

			r := httptreemux.NewContextMux()
			r.UseHandler(tt.args.middleware)

			r.GET(path, tt.args.handler.ServeHTTP)

			r.ServeHTTP(w, req)

			gotBody, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, w.Result().StatusCode)
			assert.Equalf(t, tt.wantBody, gotBody, "expected: %s\ngot: %s\n", tt.wantBody, gotBody)

			// Now the same without the router
			newW := httptest.NewRecorder()
			newReq := httptest.NewRequest(http.MethodGet, path, nil)

			tt.args.middleware(tt.args.handler).ServeHTTP(newW, newReq)
			gotBodyW, errW := io.ReadAll(newW.Result().Body)
			require.NoError(t, errW)

			assert.Equal(t, tt.wantStatus, newW.Result().StatusCode)
			assert.Equalf(t, tt.wantBody, gotBodyW, "expected: %s\ngot: %s\n", tt.wantBody, gotBodyW)
		})
	}
}
