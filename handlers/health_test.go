package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	type args struct {
		method string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   []byte
	}{
		{
			name:       "health get",
			args:       args{method: http.MethodGet},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   []byte(``),
		},
		{
			name:       "health post",
			args:       args{method: http.MethodPost},
			wantStatus: http.StatusOK,
			wantBody:   []byte(`{"message":"everything working"}`),
		},
		{
			name:       "health options",
			args:       args{method: http.MethodOptions},
			wantStatus: http.StatusOK,
			wantBody:   []byte(`{"message":"everything working"}`),
		},
		{
			name:       "healt delete",
			args:       args{method: http.MethodDelete},
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   []byte(``),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.args.method, "/health", nil)
			w := httptest.NewRecorder()
			l := zerolog.Nop()

			g := gin.New()
			g.HandleMethodNotAllowed = true
			g.NoMethod(NoMethod())
			g.Handle(http.MethodPost, "/health", Health(l))
			g.Handle(http.MethodOptions, "/health", Health(l))

			g.ServeHTTP(w, req)

			body, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, w.Result().StatusCode)
			assert.Equalf(t, tt.wantBody, body, "expected: %s\ngot       %s\n", tt.wantBody, body)
		})
	}
}
