package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   []byte
	}{
		{
			name:       "healthpoint is 200",
			wantStatus: http.StatusOK,
			wantBody: []byte(`{"message":"everything working"}
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			err := Health()(c)

			body, err := io.ReadAll(rec.Result().Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, rec.Result().StatusCode)
			assert.Equalf(t, tt.wantBody, body, "got:  %s\nwant: %s", body, tt.wantBody)
		})
	}
}
