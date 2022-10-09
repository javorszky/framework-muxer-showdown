package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReturnsAppError(t *testing.T) {
	tests := []struct {
		name       string
		handler    echo.HandlerFunc
		method     string
		path       string
		wantStatus int
		wantBody   []byte
	}{
		{
			name:       "application error returns 500",
			handler:    ReturnsAppError(),
			method:     http.MethodGet,
			path:       "/app-error",
			wantStatus: http.StatusInternalServerError,
			wantBody: []byte(`{"Message":"app error: application error: some error from someplace"}
`),
		},
		{
			name:       "request error returns 400",
			handler:    ReturnsRequestError(),
			method:     http.MethodGet,
			path:       "/req-error",
			wantStatus: http.StatusBadRequest,
			wantBody: []byte(`{"Message":"bad request request error: hurr, bad request, yarr"}
`),
		},
		{
			name:       "not found error returns 404",
			handler:    ReturnsNotFoundError(),
			method:     http.MethodGet,
			path:       "/notfound-error",
			wantStatus: http.StatusNotFound,
			wantBody: []byte(`{"Message":"not found: not found: not found the thing"}
`),
		},
		{
			name:       "shutdown error returns 500",
			handler:    ReturnsShutdownError(),
			method:     http.MethodGet,
			path:       "/shutdown-error",
			wantStatus: http.StatusServiceUnavailable,
			wantBody: []byte(`{"Message":"well this is bad: shutdown error: unrecoverable error"}
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			errChan := make(chan error, 1)
			l := zerolog.Nop()

			e := echo.New()
			e.HTTPErrorHandler = CustomErrorHandler(l, errChan)
			e.Add(tt.method, tt.path, tt.handler)

			e.ServeHTTP(rec, req)

			resultBody, err := io.ReadAll(rec.Result().Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, rec.Result().StatusCode)
			assert.Equalf(t, tt.wantBody, resultBody, "got:  %s\nwant: %s", resultBody, tt.wantBody)
		})
	}
}
