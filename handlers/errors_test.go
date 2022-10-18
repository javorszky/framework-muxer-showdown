package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandlers(t *testing.T) {
	ch := make(chan error, 1)
	errHandler := ErrorHandler(zerolog.Nop(), ch)

	tests := []struct {
		name       string
		method     string
		path       string
		handler    fiber.Handler
		wantStatus int
		wantBody   []byte
	}{
		{
			name:       "app error",
			method:     http.MethodGet,
			path:       "/app-error",
			handler:    ReturnsApplicationError(),
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"message":"app error: application error: some error from someplace"}`),
		},
		{
			name:       "req error",
			method:     http.MethodGet,
			path:       "/req-error",
			handler:    ReturnsRequestError(),
			wantStatus: http.StatusBadRequest,
			wantBody:   []byte(`{"message":"bad request request error: hurr, bad request, yarr"}`),
		},
		{
			name:       "notfound error",
			method:     http.MethodGet,
			path:       "/notfound-error",
			handler:    ReturnsNotFoundError(),
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(`{"message":"not found: not found: not found the thing"}`),
		},
		{
			name:       "shutdown error",
			method:     http.MethodGet,
			path:       "/shutdown-error",
			handler:    ReturnsShutdownError(),
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   []byte(`{"message":"well this is bad: shutdown error: unrecoverable error"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				StrictRouting:     true,
				BodyLimit:         2 * 1024,
				ReadTimeout:       30 * time.Second,
				WriteTimeout:      2 * time.Minute,
				IdleTimeout:       2 * time.Minute,
				AppName:           "fiber-test",
				EnablePrintRoutes: true,
				ErrorHandler:      errHandler,
			})
			app.Add(tt.method, tt.path, tt.handler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			resp, err := app.Test(req)
			require.NoError(t, err)

			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, resp.StatusCode, "status codes don't match")
			assert.Equalf(t, tt.wantBody, b, "want: %s\ngot:  %s", tt.wantBody, b)
		})
	}
}
