package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReturnsRequestError(t *testing.T) {
	testChan := make(chan error, 1)
	testLog := zerolog.Nop()

	tests := []struct {
		name       string
		handler    httprouter.Handle
		wantStatus int
		wantBody   []byte
	}{
		{
			name:       "returns application error",
			handler:    ErrorCatcher(testLog, testChan)(ReturnsApplicationError()),
			wantStatus: http.StatusInternalServerError,
			wantBody:   []byte(`{"message":"app error: application error: some error from someplace"}`),
		},
		{
			name:       "returns shutdown error",
			handler:    ErrorCatcher(testLog, testChan)(ReturnsShutdownError()),
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   []byte(`{"message":"well this is bad: shutdown error: unrecoverable error"}`),
		},
		{
			name:       "returns notfound error",
			handler:    ErrorCatcher(testLog, testChan)(ReturnsNotFoundError()),
			wantStatus: http.StatusNotFound,
			wantBody:   []byte(`{"message":"not found: not found: not found the thing"}`),
		},
		{
			name:       "returns request error",
			handler:    ErrorCatcher(testLog, testChan)(ReturnsRequestError()),
			wantStatus: http.StatusBadRequest,
			wantBody:   []byte(`{"message":"bad request request error: hurr, bad request, yarr"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			w := httptest.NewRecorder()

			r := httprouter.New()
			r.GET("/path", tt.handler)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Result().StatusCode)

			b, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			assert.Equalf(t, tt.wantBody, b, "want: %s\ngot:  %s", tt.wantBody, b)

		})
	}
}
