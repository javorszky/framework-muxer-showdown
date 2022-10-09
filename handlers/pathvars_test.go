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

func TestPathVars(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   []byte
	}{
		{
			name:       "correctly figures out path",
			path:       "/pathvars/testOne/metrics/testTwo",
			wantStatus: http.StatusOK,
			wantBody:   []byte(`pathvar1: testOne, pathvar2: testTwo`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			e := echo.New()
			e.Add(http.MethodGet, "/pathvars/:one/metrics/:two", PathVars())
			e.ServeHTTP(rec, req)

			body, err := io.ReadAll(rec.Result().Body)
			require.NoError(t, err)

			assert.Equalf(t, tt.wantStatus, rec.Result().StatusCode, "Status code")
			assert.Equalf(t, tt.wantBody, body, "got:  %s\nwant: %s", body, tt.wantBody)
		})
	}
}
