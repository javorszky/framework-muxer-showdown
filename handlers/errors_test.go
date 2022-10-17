package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	localErrors "github.com/suborbital/framework-muxer-showdown/errors"
	"github.com/suborbital/framework-muxer-showdown/web"
)

func TestErrorHandlers(t *testing.T) {
	testlogger := zerolog.Nop()
	testchan := make(chan error)

	type args struct {
		handler http.Handler
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   []byte
		wantError  func(error) bool
	}{
		{
			name:       "app error handler returns app error",
			args:       args{handler: ErrorCatcher(testlogger, testchan)(ReturnsApplicationError(testlogger))},
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"Message":"app error: application error: some error from someplace"}`),
			wantError:  localErrors.IsApplicationError,
		},
		{
			name:       "Not found error handler returns not found error",
			args:       args{handler: ErrorCatcher(testlogger, testchan)(ReturnsNotFoundError())},
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(`{"Message":"not found: not found: not found the thing"}`),
			wantError:  localErrors.IsNotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/app-error", nil)
			w := httptest.NewRecorder()

			tt.args.handler.ServeHTTP(w, req)

			gotBody, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, w.Result().StatusCode)
			assert.Equalf(t, tt.wantBody, gotBody, "expected: %s\ngot: %s\n", tt.wantBody, gotBody)
			assert.True(t, tt.wantError(web.GetErrors(req.Context())[0]))
		})
	}
}
